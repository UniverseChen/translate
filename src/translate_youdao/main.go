package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"pronunciation_bing"
	"strconv"
	tydjson "translate_youdao/json"
)

// FLAG
var tranLanguage *int = flag.Int("l", 0, "0 is EN->zh_CHS(default)\n\t1 is zh_CHS->EN")

// 有道API地址
const httpUrl = "http://openapi.youdao.com/api"

// 私人密匙，需自行到 http://ai.youdao.com/　进行申请，填写正确后即可运行
// appKey		用于请求报文
// sercetKey	用于生成随机md5验证码
const (
	appKey    = "004a62508d92b47f"
	secretKey = "dQ4D8it8KPOulzE4RfywtDDjQ6rl2oaV"
)

// 语言
const (
	simpleChinese = "zh_CHS"
	english       = "EN"
)

// 默认的翻译文本
var translateText string

// GetMd5	md5验证码生成
func GetMd5(str string) (string, error) {

	// 空字符串处理
	if str == "" {
		return "", errors.New("str is nil")
	}

	h := md5.New()
	h.Write([]byte(str))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr), nil
}

// GetRequest	拼接请求报文
func GetRequest(httpUrl, q, from, to, appKey, salt, sign string) string {
	/*
		请求报文格式:
		http://openapi.youdao.com/api?					// 有道api网址
		q=good&											// q:	代表要翻译的文本
		from=EN&										// from:源语言
		to=zh_CHS&										// to:	目标语言
		appKey=ff889495-4b45-46d9-8f48-946554334f2a&	// appKey:自行申请
		salt=2&											// 随机数
		sign=1995882C5064805BC30A39829B779D7B			// 生成的md5验证码
	*/
	return httpUrl + "?" + "q=" + q + "&" + "from=" + from + "&" + "to=" + to +
		"&" + "appKey=" + appKey + "&" + "salt=" + salt + "&" + "sign=" + sign
}

// GetResponse	发送http请求，获取响应报文
func GetResponse(requestion string) (string, error) {
	resp, err := http.Get(requestion)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	return string(body), nil
}

// init	初始化
func init() {
	flag.Parse()
}

func main() {

	// 翻译语言处理
	var languageForm, languageTo string
	switch *tranLanguage {
	case 0:
		languageForm, languageTo = english, simpleChinese
	case 1:
		languageForm, languageTo = simpleChinese, english
	}

	// 主循环
	for {

		// 从键盘获取输入
		fmt.Scanf("%s", &translateText)

		// 生成 0-9 随机数，并转换为字符串
		salt := strconv.Itoa(rand.Intn(9))
		// 生成md5验证码，appkey + 翻译文本 + 随机数 + secretkey
		m := appKey + translateText + salt + secretKey
		md5Str, err := GetMd5(m)
		if err != nil {
			fmt.Printf("GetMd5 Failed: %s", err.Error())
			break
		}

		// 请求报文
		requestion := GetRequest(httpUrl, translateText, languageForm,
			languageTo, appKey, salt, md5Str)

		// 获取响应报文，包含 HTTP 请求
		tran, err := GetResponse(requestion)
		if err != nil {
			fmt.Printf("GetResponse Failed: %s", err.Error())
			break
		}

		// 解析JSON数据并输出
		tran, err = tydjson.ParseJson(tran)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(tran)
		}

		go pronunciation_bing.RunPlay(translateText)
	}
}
