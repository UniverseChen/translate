package json

import (
	"encoding/json"
	"errors"
	"strings"
)

// 错误返回码的含义
const (
	ERR_101 = "缺少必填的参数，出现这个情况还可能是et的值和实际加密方式不对应"
	ERR_102 = "不支持的语言类型"
	ERR_103 = "翻译文本过长"
	ERR_104 = "不支持的API类型"
	ERR_105 = "不支持的签名类型"
	ERR_106 = "不支持的响应类型"
	ERR_107 = "不支持的传输加密类型"
	ERR_108 = "	appKey无效，注册账号， 登录后台创建应用和实例并完成绑定， 可获得应用ID和密钥等信息，其中应用ID就是appKey（ 注意不是应用密钥）"
	ERR_109 = "batchLog格式不正确"
	ERR_110 = "无相关服务的有效实例"
	ERR_111 = "开发者账号无效，可能是账号为欠费状态"
	ERR_201 = "解密失败，可能为DES,BASE64,URLDecode的错误"
	ERR_202 = "签名检验失败"
	ERR_203 = "访问IP地址不在可访问IP列表"
	ERR_301 = "	辞典查询失败"
	ERR_302 = "翻译查询失败"
	ERR_303 = "服务端的其它异常"
	ERR_401 = "账户已经欠费停用"
)

// JSON结构 web项单一元素类型
type jsonWeb struct {
	Value []string `value`
	Key   string   `key`
}

type jsonBasic struct {
	Phonetic    string `phonetic`
	Us_phonetic string `us-phonetic`
	Uk_phonetic string `uk-phonetic`
}

// JSON结构 内部结构，json解析用
type YouDaoApiJson struct {
	ErrorCode   int       `errorCode`
	Query       string    `query`
	Translation []string  `translation`
	L           string    `l`
	Web         []jsonWeb `web`
	Basic       jsonBasic `basic`
}

// json中的 us-phonetic 和 uk-phonetic 两项由于 '-' 无法被json解析，所以更改为 '_'
func changeUsUk(str *string) {

	// 要处理的子串
	var sep_us, sep_uk string = "us-phonetic", "uk-phonetic"

	*str = strings.Replace(*str, sep_uk, "uk_phonetic", 1)
	*str = strings.Replace(*str, sep_us, "us_phonetic", 1)

}

// func ParseJson		json解析
func ParseJson(response string) (string, error) {

	// 修改 '-' 为 '_'
	changeUsUk(&response)

	var s YouDaoApiJson                  // json实例，用于解析用
	json.Unmarshal([]byte(response), &s) // 解析json

	// 处理数据
	if s.ErrorCode == 0 { // 返回的错误码，== 0 代表正常
		var str string // 作返回值

		// 处理 translation
		// Translation 是一个 字符串数组，所以先拼接为字符串, 格式为: 元素1; 元素2; 元素3;
		var strTranslation string
		for _, v := range s.Translation {
			strTranslation += v + "; "
		}
		// 拼接基础数据并赋值给str，格式为:
		// 原文: 翻译
		str = "\t" + s.Query + ": " + strTranslation + "\n"

		// 处理 web
		// web 格式为，[](key(string), value([]string))
		for _, v := range s.Web {
			var strValue string
			for _, vv := range v.Value {
				strValue += vv + "; "
			}
			// 将 处理过的web数据拼接到str
			str += "\t" + v.Key + ": " + strValue + "\n"
		}

		// 处理 basic.phonetic
		str += "\t" + "phonetic: " + "\t" + s.Basic.Phonetic + "\n" +
			"\t" + "us-phonetic: " + "\t" + s.Basic.Us_phonetic + "\n" +
			"\t" + "uk-phonetic: " + "\t" + s.Basic.Uk_phonetic + "\n"

		return str, nil
	} else { // 不正常的错误码，做错误处理
		switch s.ErrorCode {
		case 101:
			return "", errors.New(ERR_101)
		case 102:
			return "", errors.New(ERR_102)
		case 103:
			return "", errors.New(ERR_103)
		case 104:
			return "", errors.New(ERR_104)
		case 105:
			return "", errors.New(ERR_105)
		case 106:
			return "", errors.New(ERR_106)
		case 107:
			return "", errors.New(ERR_107)
		case 108:
			return "", errors.New(ERR_108)
		case 109:
			return "", errors.New(ERR_109)
		case 110:
			return "", errors.New(ERR_110)
		case 111:
			return "", errors.New(ERR_111)
		case 201:
			return "", errors.New(ERR_201)
		case 202:
			return "", errors.New(ERR_202)
		case 203:
			return "", errors.New(ERR_203)
		case 301:
			return "", errors.New(ERR_301)
		case 302:
			return "", errors.New(ERR_302)
		case 303:
			return "", errors.New(ERR_303)
		case 401:
			return "", errors.New(ERR_401)
		default:
			return "", errors.New("未知错误，返回码无效")
		}
	}
}
