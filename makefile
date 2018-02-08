build_main:pronunciation_bing translate_youdao_json
	go build src/translate_youdao/main.go

pronunciation_bing:
	go install pronunciation_bing

translate_youdao_json:
	go install translate_youdao/json
