package wxPusher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Loli struct {
	// 密钥
	AppToken string `json:"appToken"`
	// 内容
	Content string `json:"content"`
	// 消息摘要
	Summary string `json:"summary"`
	// 内容类型
	ContentType int `json:"contentType"`
	// 主题id 为空不转发
	TopicIds []int `json:"topicIds"`
	// 个人id 为空不转发
	Uids []string `json:"uids"`
}

func post(url string, requestBody *[]byte, headers *map[string]string) *[]byte {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, bytes.NewReader(*requestBody))
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

func (loli *Loli) Send(appToken string, content string, summary string, contentType int, topicIds []int, uIds []string) *[]byte {
	url := "https://wxpusher.zjiecode.com/api/send/message"
	loli = &Loli{appToken, content, summary, contentType, topicIds, uIds}
	requestBody, _ := json.Marshal(&loli)
	resp := post(url, &requestBody, nil)
	if resp == nil {
		return nil
	}
	return resp
}
