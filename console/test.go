package console

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
)

func GetBase64ByPic1(path string) {
	srcByte, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	res := base64.StdEncoding.EncodeToString(srcByte)
	fmt.Println(res)
}

func GetBase64ByPic2(path string) {

	//读取文件
	ff, err := os.Open(path)
	if err != nil {
		return
	}
	defer ff.Close()
	//获取文件大小
	stat, _ := ff.Stat()

	//读取文件
	bs := make([]byte, stat.Size())
	_, err = ff.Read(bs)
	if err != nil {
		return
	}
	//编码
	getBase64 := base64.StdEncoding.EncodeToString(bs)

	fmt.Println(getBase64)
}

func main() {
	path := "/Users/mr.kevin/Downloads/123123.png"
	fmt.Println(path)
	//console.GetBase64ByPic1(path)
	//console.GetBase64ByPic2(path)
}
