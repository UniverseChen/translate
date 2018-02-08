/*
需要事先安装sox:	sudo apt-get install sox
通过解析从 https://dict.bing.com 返回的 html page, 获取音频文件(.mp3)的下载链接，然后进行下载-保存-播放
保存目录默认为 $HOME/Music/Pronunciation_bing
文件名默认为 dictname_US.mp3


https://cn.bing.com/dict/search?q=sir&qs=BD&pq=sir&sc=8-6&cvid=7F3E2F5372A0492BB1B1600D57FCF269&sp=1

https://cn.bing.com/dict/
search?q=
indent
&qs=BD
&pq=
indent
&sc=8-6
&cvid=7F3E2F5372A0492BB1B1600D57FCF269
&sp=1

*/

package pronunciation_bing
