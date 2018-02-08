package pronunciation_bing

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

//key words for search in html
const (
	us_key    = "美&#160"
	uk_key    = "美&#160"
	https_key = "https://"
	mp3_key   = ".mp3"
	file_type = "mp3"
)

// CreateRequestMsg	request message	请求报文
func CreateRequestMsg(dict string) string {
	return "https://cn.bing.com/dict/search?q=" + dict + "&qs=BD" + "&pq=" + dict + "&sc=8-6" +
		"&cvid=7F3E2F5372A0492BB1B1600D57FCF269" + "&sp=1"
}

// GetUS	search the USaudio's url in html	搜索美式发音的音频url
func GetUS(html string) (url_us_mp3 string) {
	us_index := strings.Index(html, us_key)
	if us_index == -1 {
		return url_us_mp3
	} else {
		https_index := strings.Index(html[us_index:], https_key)
		if https_index != -1 {
			mp3_index := strings.Index(html[us_index:], mp3_key)
			if mp3_index != -1 {
				url_us_mp3 = html[us_index+https_index : us_index+mp3_index+len(mp3_key)]
			}
		}
	}
	return url_us_mp3
}

//SaveFile save the stream file	保存流文件
func SaveFile(name string, src io.Reader) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, src)

	return nil
}

// PlayCmd	play mp3 (sox)	通过命令行播放，使用了linux的命令行软件sox
func PlayCmd(name string) {
	playCmd := exec.Command("play", name)
	playCmd.Run()
}

// FileExist	is the file exist	检测文件是否存在
func FileExist(name string) bool {
	_, err := os.Stat(name)
	return err == nil || os.IsExist(err)
}

var path_home, path_video_bin string

func init() {
	path_home = os.Getenv("HOME")
	path_video_bin = path_home + "/" + "Music/" + "Pronunciation_bing/"
	c := exec.Command("mkdir", "-p", path_video_bin)
	c.Run()
}

// RunPlay	run and play 运行并播放
// 没有认真写，错误处理基本敷衍
func RunPlay(dict /*, pith*/ string) {
	//dict of translate

	//file name
	f_path := path_video_bin + dict + "_US" + "." + file_type

	if FileExist(f_path) {
		PlayCmd(f_path)
		return
	}

	// get pronunciation's url from []byte-type
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

	//US
	resp, err := http.Get(GetUS(string(robots)))
	if err != nil {
		fmt.Printf("Not Find \"%s\"'s pronunciation\n", dict)
		return
	}

	pix, _ := ioutil.ReadAll(resp.Body)
	SaveFile(f_path, bytes.NewReader(pix))
	PlayCmd(f_path)

	resp.Body.Close()
}
