package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	//"bufio"
	"bufio"
)

//key words for search in html
const (
	us_key    = "美&#160"
	uk_key    = "美&#160"
	https_key = "https://"
	mp3_key   = ".mp3"
	file_type = "mp3"
)

//CreateRequestMsg	根据需要翻译的文本创建请求报文
func CreateRequestMsg(dict string) string {
	//修改空格为加号 ' ' to '+'
	dict_bytes := []byte(dict)
	for i, _ := range dict_bytes {
		if dict_bytes[i] == ' ' {
			dict_bytes[i] = '+'
		}
	}

	return "https://cn.bing.com/dict/search?q=" + string(dict_bytes) + "&qs=BD" + "&pq=" + string(dict_bytes) + "&sc=8-6" +
		"&cvid=7F3E2F5372A0492BB1B1600D57FCF269" + "&sp=1"
}

//PreaseText	根据html文本解析翻译文本，并格式化后返回
//基本纯字符串操作，对照注释和html文件更容易理解
func PreaseText(html string) string {
	//从html中提取出翻译文本	value-name: desc is description
	//"hello"的数据(含标签)，对照查看用
	//<meta name="desc" content="必应词典为您提供hello的释义，美[heˈləʊ]，英[hə'ləʊ]，int. 你好；喂；您好；哈喽； 网络释义： 哈罗；哈啰；大家好； " />

	//dict_desc_index 含html标签的翻译文本的索引, "<meta name=\"desc\"":检索使用的关键字
	dict_desc_index := strings.Index(html, "<meta name=\"description\"")

	//dict_desc_value 含html标签的翻译文本的数据, "/>": 结束的关键字
	dict_desc_value := html[dict_desc_index : dict_desc_index+strings.Index(html[dict_desc_index:], "/>")+len("/>")]

	//str 是最终提取出的翻译文本，含全角符号，不含html标签(hello演示):
	//美[heˈləʊ]，英[hə'ləʊ]，int. 你好；喂；您好；哈喽； 网络释义： 哈罗；哈啰；大家好；
	str := dict_desc_value[strings.Index(dict_desc_value, "释义，")+len("释义，") : len(dict_desc_value)-5]

	//格式化的翻译文本，作返回值用
	var format_tran string
	//临时存储网络释义(网络释义放在最后显示，但是先提取出来方便操作)
	var str_web string

	//查找音标是否存在，若存在则提前提取并写入format_tran
	if strings.Index(str, "]" /*]　作为音标是否存在的关键字*/) != -1 {
		last_index := strings.LastIndex(str, "]" /*最后一个 ]，即音标的结尾*/) + 1 /*左闭右开[)，所以加1*/
		//将音标写入 format_tran，这里不需要"\n"，后面的释义会前置写入"\n"
		format_tran = str[:last_index]
		//更新str
		str = str[last_index+len("，")/*去除分号，是全角符号，所以使用len(...)*/: ]
	}

	//查找网络释义是否存在，若存在则提前提取并保存到　str_web
	if strings.Index(str, "网络释义") != -1 {
		n := strings.Index(str, "网络释义")
		str_web = str[n:]
		str = str[:n] //更新str
	}

	//pos_num　代表释义类型的数量(名词，动词，形容词等，一种就需要另起一行，除开网络释义)
	var pos_num int = strings.Count(str, ". " /*'.'代表一种释义，如: n. v. adj. 等*/)

	//如果除网络释义外，只有一个释义或没有，则直接与网络释义一同写入format_tran并返回
	if pos_num == 0 || pos_num == 1 {
		format_tran = str + "\n" + str_web
		return format_tran
	}

	//如果除网络释义外，还有两个及以上释义，则用此循环处理
	for {
		//处理方法不是寻找第一个释义的开始，而是通过"."寻找第二个释义的index，然后向前寻找"；"(全角分号)，
		//然后将前一段写入format_tran并更新str，然后重复操作直到最后一个释义

		//第一个"."
		index_first := strings.Index(str, ". ") + 1 /*跳过已经找到的第一个"."*/
		//第二个"."，重要，用于判断是否到了最有一个释义，如果是，则会置为 -1
		index_second := strings.Index(str[index_first:], ". ")
		//如果不是最有一个释义，则为 index_second　赋正确的值
		if index_second != -1 {
			index_second = index_first + index_second /*index_second是在index_second的基础上所引到的*/
		}
		//如果是最后一个
		if index_second == -1 {
			format_tran = format_tran + "\n" + str[:len(str)-len("； ")] + "\n" + str_web
			return format_tran
		} else {
			//第一个释义的尾索引
			end_first := strings.LastIndex(str[:index_second], "；" /*全角分号*/)
			//写入format_tran
			format_tran = format_tran + "\n" + str[:end_first]
			//更新str
			str = str[end_first+len("； "):]
		}
	}
}

//PreasePronunciationUS	查找并返回美式发音的mp3文件的url
func PreasePronunciationUS(html string) string {
	//us_key = "美&#160",是寻找美式发音url的关键字
	//下面为从html(hello)截取的一段，包含了音频文件url，且可以通过 us_key 索引到
	//<div class="hd_prUS">美&#160;</div><div class="hd_tf"><a class="bigaud" onmouseover="this.className='bigaud_f';javascript:BilingualDict.Click(this,'https://dictionary.blob.core.chinacloudapi.cn/media/audio/tom/9c/79/9C79B32346052BE464270F9083081FB7.mp3','akicon.png',false,'dictionaryvoiceid')"
	us_index := strings.Index(html, us_key)
	if us_index == -1 { //如果没有发音，则返回空字符串
		return ""
	} else {
		https_index := strings.Index(html[us_index:], https_key) //索引到url前缀(https://)
		if https_index != -1 {                                   //url前缀存在(即url存在)
			mp3_index := strings.Index(html[us_index:], mp3_key) //索引到url后缀(.mp3)
			if mp3_index != -1 {                                 //url后缀存在(可以完整索引到url)
				return html[us_index+https_index : us_index+mp3_index+len(mp3_key)] //返回完整url
			}
		}
	}
	return "" //若url前缀或url后缀不存在，则返回空字符串
}

//SaveFile	保存流文件(下载的mp3)
func SaveFile(name string, src io.Reader) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, src)

	return nil
}

//PlayCmd	与系统相关, 通过命令行播放，使用了linux的命令行软件sox
func PlayCmd(name string) {
	playCmd := exec.Command("play", name)
	playCmd.Run()
}

//FileExist	检测文件是否存在
func FileExist(name string) bool {
	_, err := os.Stat(name)
	return err == nil || os.IsExist(err)
}

//下载文件保存的路径，默认为 path_home = $HOME; path_video_bin = $HOME/Music/Pronunciation/bing/
var path_home, path_video_bin string

//RunPlay	运行并播放
//没有认真写，错误处理基本敷衍
func RunPlay(dict /*, pith*/ string) {
	//dict of translate

	//file name		path_video_bin/xxxx(name)_US.mp3
	f_path := path_video_bin + dict + "_US" + "." + file_type

	//通过Get请求获取响应报文，然后将Body读入 robots(type is []byte)
	res, err := http.Get(CreateRequestMsg(dict))
	if err != nil {
		fmt.Printf("Get request failed: %s\n", err.Error())
		return
	}
	robots, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Read html is failed: %s\n", err.Error())
		return
	}

	//打印翻译文本
	fmt.Println(PreaseText(string(robots)))
	fmt.Println()

	//音频处理
	//如果已经存在此mp3文件，则直接播放，然后返回
	if FileExist(f_path) {
		PlayCmd(f_path)
		return
	}

	//解析并获取音频文件的url，然后下载mp3文件
	resp, err := http.Get(PreasePronunciationUS(string(robots)))
	if err != nil {
		fmt.Printf("Not Find %s's pronunciation\n\n", dict)
		return
	}
	defer resp.Body.Close()
	//保存下载的数据
	pix, _ := ioutil.ReadAll(resp.Body)
	SaveFile(f_path, bytes.NewReader(pix))
	//播放
	PlayCmd(f_path)

}

func init() {
	path_home = os.Getenv("HOME") //从环境变量获取 $HOME　并赋值
	path_video_bin = path_home + "/" + "Music/" + "Pronunciation/bing/"
	c := exec.Command("mkdir", "-p", path_video_bin) //调用shell命令创建文件夹
	c.Run()
}

func main() {

	for {
		//读入一行数据而不是遇到空格结束
		reader := bufio.NewReader(os.Stdin)
		str_bytes, _, _ := reader.ReadLine()

		RunPlay(string(str_bytes))
	}

}
