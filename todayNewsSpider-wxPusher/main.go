package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"loli/wxPusher"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
)

type Config struct {
	// 密钥
	AppToken string `json:"appToken"`
	// 主题id 为空不转发
	TopicIds []int `json:"topicIds"`
	// 个人id 为空不转发
	Uids []string `json:"uids"`
	// 状态码
	Data []map[string]interface{}
}

type Image struct {
	ImageContent struct {
		Image struct {
			Url string
		}
	}
}

type BingWallpaper struct {
	MediaContents []Image
}

// 获取Bing当日壁纸
func getBingWallpaper() string {
	var bing BingWallpaper
	url := "https://cn.bing.com/hp/api/model"
	resp := get(url, nil)
	json.Unmarshal(*resp, &bing)
	return bing.MediaContents[0].ImageContent.Image.Url
}

func get(url string, headers *map[string]string) *[]byte {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")
	if headers != nil {
		for k, v := range *headers {
			req.Header.Set(k, v)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("网络连接超时!", err)
		return nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return &body
}

func parse() {
	url := "https://www.163.com/dy/media/T1603594732083.html"
	resp := get(url, nil)
	if resp == nil {
		return
	}
	html, _ := htmlquery.Parse(bytes.NewReader(*resp))
	root := htmlquery.Find(html, `//li[@class="js-item item"][1]/a/@href`)
	todayUrl := htmlquery.InnerText(root[0])
	headers := map[string]string{
		"referer": url,
	}
	resp = get(todayUrl, &headers)
	if resp == nil {
		return
	}
	html, _ = htmlquery.Parse(bytes.NewReader(*resp))
	root = htmlquery.Find(html, `//div[@class="post_body"]/p[2]/text()[position()>1]`)
	var con string
	for i, v := range root {
		if i == 0 {
			con += "![](" + getBingWallpaper() + ")\n\n**<center>" + htmlquery.InnerText(v) + "</center>**\n\n"
			continue
		}
		if i == len(root)-1 {
			con += "\n> " + htmlquery.InnerText(v)
			break
		}
		con += strings.Replace(htmlquery.InnerText(v), "、", ". ", 1) + "\n"
	}
	push(con, "今日早报")
}

func push(content string, summary string) {
	var loli wxPusher.Loli
	var con Config
	readConfig(&con)
	if con.AppToken == "" || (con.TopicIds[0] == 0 && con.Uids[0] == "") {
		fmt.Println("config.json配置错误!")
		return
	}
	if con.TopicIds[0] == 0 {
		con.TopicIds = nil
	}
	if con.Uids[0] == "" {
		con.Uids = nil
	}
	res := loli.Send(con.AppToken, content, summary, 3, con.TopicIds, con.Uids)
	json.Unmarshal(*res, &con)
	fmt.Println(con.Data[0]["code"], con.Data[0]["status"])
}

func readConfig(con *Config) {
	path := filepath.Dir(os.Args[0]) + filepath.FromSlash("/") + "config.json"
	if _, err := os.Stat(path); err != nil {
		f, err := os.Create(path)
		if err != nil {
			fmt.Println("config.json创建失败！", err)
			time.Sleep(3 * time.Second)
			os.Exit(0)
		}
		defer f.Close()
		con.TopicIds = make([]int, 1)
		con.Uids = make([]string, 1)
		c, _ := json.MarshalIndent(con, "", "    ")
		f.Write(c)
		fmt.Println("请在config.json中填入必要的配置！")
		time.Sleep(3 * time.Second)
		os.Exit(0)
	}
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	fi, _ := f.Stat()
	data := make([]byte, fi.Size())
	f.Read(data)
	json.Unmarshal(data, &con)
}

func main() {
	parse()
}
