package common

import (
	"fmt"
	"github.com/robfig/config"
	"io/ioutil"
	"os"
)

//写入内容到文件
func SaveFile(strFileName string, strData string) (ok bool) {
	f, err := os.Create(strFileName)
	if err != nil {
		fmt.Println("create file faild error:", err)
		return false
	}

	_, err_w := f.Write([]byte(strData))
	if err_w != nil {
		fmt.Println("Server start faild error:", err_w)
		return false
	}
	return true
}

//追加内容到文件
func AppendFile(strFileName string, strData string) (ok bool) {
	strData = ReadFile(strFileName) + "\n" + strData
	fmt.Println("数据", strData)
	return SaveFile(strFileName, strData)
}

//从文件中读取内容
func ReadFile(filePth string) string {
	bytes, err := ioutil.ReadFile(filePth)
	if err != nil {
		fmt.Println("读取文件失败: ", err)
		return ""
	}

	return string(bytes)
}

//截取固定位置以前的字符串
func SubstrBefore(s string, l int) string {
	if len(s) <= l {
		return s
	}
	ret, rs := "", []rune(s)

	for i, r := range rs {
		if i >= l {
			break
		}

		ret += string(r)
	}
	return ret
}

//截取固定位置以后的字符串
func SubstrAfter(s string, l int) string {
	if len(s) <= l {
		return s
	}
	ret, rs := "", []rune(s)

	for i, r := range rs {
		if i <= l {
			continue
		}

		ret += string(r)
	}
	return ret
}

//错误处理
func ErrorHandle(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

//获取配置文件信息
func GetConfig(section string, option string) string {
	c, _ := config.ReadDefault("conf/config.cfg")
	value, err := c.String(section, option)
	if err != nil {
		fmt.Println("读取配置文件出错", err)
		return ""
	}
	return value
}

//设置配置文件信息
func SetConfig(section string, option string, value string) {
	c, _ := config.ReadDefault("conf/config.cfg")
	c.AddOption(section, option, value)
	c.WriteFile("config.cfg", 0644, "A header for this file")
}
