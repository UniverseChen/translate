/*
语言代码
		语言	代码
		中文	zh-CHS
		日文	ja
		英文	EN
		韩文	ko
		法文	fr
		俄文	ru
		葡萄牙文	pt
		西班牙文	es


请求报文格式(dict: good)
		http://
		openapi.youdao.com/api?
		q=good&
		from=EN&
		to=zh_CHS&
		appKey=ff889495-4b45-46d9-8f48-946554334f2a&
		salt=2&
		sign=1995882C5064805BC30A39829B779D7B


返回结果字段信息
		返回的结果是json格式，包含字段与FROM和TO的值有关，具体说明如下：

		字段名	类型	含义	备注
		errorCode	text	错误返回码	一定存在
		query	text	源语言	查询正确时，一定存在
		translation	text	翻译结果	查询正确时一定存在
		basic	text	词义	基本词典,查词时才有
		web	text	词义	网络释义，该结果不一定存在
		l	text	源语言和目标语言	一定存在
		dict	text	词典deeplink	查询语种为支持语言时，存在
		webdict	text	webdeeplink	查询语种为支持语言时，存在


示例: hello 的翻译结果
		{"web":[{"value":["你好","您好","hello"],"key":"Hello"},{"value":["凯蒂猫","昵称","匿称"],"key":"Hello Kitty"},{"value":["哈乐哈乐","乐扣乐扣"],"key":"Hello Bebe"}],"query":"hello","translation":["你好"],"errorCode":"0","dict":{"url":"yddict://m.youdao.com/dict?le=eng&q=hello"},"webdict":{"url":"http://m.youdao.com/dict?le=eng&q=hello"},"basic":{"us-phonetic":"həˈlo","phonetic":"həˈləʊ","uk-phonetic":"həˈləʊ","explains":["n. 表示问候， 惊奇或唤起注意时的用语","int. 喂；哈罗","n. (Hello)人名；(法)埃洛"]},"l":"EN2zh-CHS"}

		分解:
			{
			"web":
				[
				{"value":["你好","您好","hello"],"key":"Hello"},
				{"value":["凯蒂猫","昵称","匿称"],"key":"HelloKitty"},
				{"value":["哈乐哈乐","乐扣乐扣"],"key":"Hello Bebe"}
				],
				"query":
					"hello",
				"translation":
					[
					"你好"
					],
				"errorCode":
					"0",
				"dict":
					{"url":"yddict://m.youdao.com/dict?le=eng&q=hello"},
				"webdict":
					{"url":"http://m.youdao.com/dict?le=eng&q=hello"},
				"basic":
					{"us-phonetic":"həˈlo",
					"phonetic":"həˈləʊ",
					"uk-phonetic":"həˈləʊ",
					"explains":["n. 表示问候， 惊奇或唤起注意时的用语",
					"int. 喂；哈罗",
					"n. (Hello)人名；(法)埃洛"]},
					"l":"EN2zh-CHS"}
		}


errorCode　信息
		101	缺少必填的参数，出现这个情况还可能是et的值和实际加密方式不对应
		102	不支持的语言类型
		103	翻译文本过长
		104	不支持的API类型
		105	不支持的签名类型
		106	不支持的响应类型
		107	不支持的传输加密类型
		108	appKey无效，注册账号， 登录后台创建应用和实例并完成绑定， 可获得应用ID和密钥等信息，其中应用ID就是appKe) (注意不是应用密匙)
		109	batchLog格式不正确
		110	无相关服务的有效实例
		111	开发者账号无效，可能是账号为欠费状态
		201	解密失败，可能为DES,BASE64,URLDecode的错误
		202	签名检验失败
		203	访问IP地址不在可访问IP列表
		301	辞典查询失败
		302	翻译查询失败
		303	服务端的其它异常
		401	账户已经欠费停
*/

package translate_youdao
