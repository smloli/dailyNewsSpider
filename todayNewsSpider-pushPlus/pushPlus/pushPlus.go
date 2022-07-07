package pushPlus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Loli struct {
	// 密钥
	Token string
	// 返回状态码
	StartCode struct {
		Code int
		Msg  string
		Data string
	}
}

func post(url string, params *[]byte, headers *map[string]string) *[]byte {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, bytes.NewReader(*params))
	req.Header.Set("Content-Type", "application/json")
	if headers != nil {
		for k, v := range *headers {
			req.Header.Set(k, v)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return &body
}

func (loli *Loli) Send(title string, content string, topic string, template string) (startCode int, msg string) {
	url := "http://www.pushplus.plus/send"
	params := map[string]string{
		// 密钥
		"token": loli.Token,
		// 标题
		"title": title,
		// 内容
		"Content": content,
		// 群组名称 为空则转发自己，否则转发群组
		"topic": topic,
		// 消息类型 markdown
		"template": template,
	}
	body, _ := json.Marshal(&params)
	resp := post(url, &body, nil)
	if resp == nil {
		return
	}
	json.Unmarshal(*resp, &loli.StartCode)
	return loli.StartCode.Code, loli.StartCode.Msg
}
