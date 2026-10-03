package http_

import (
	"bufio"
	"bytes"
	"handfree-work/octo-backup/internal/base/error_"
	"handfree-work/octo-backup/internal/base/log_"
	"io"
	"net/http"
	"net/url"

	"github.com/gogf/gf/v2/encoding/gjson"
)

func StreamRequest(in *Request, stream func(data *SseMessage) error) error {
	jsonStr, err := gjson.EncodeString(in.Body)
	if err != nil {
		return err
	}
	// 创建 HTTP 请求
	req, err := http.NewRequest(in.Method, in.Url, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf8")
	// in.headers 循环设置 HTTP 请求头
	for k, v := range in.Headers {
		req.Header.Set(k, v)
	}

	req.Header.Set("Accept", "text/event-stream")

	Transport := &http.Transport{}
	if in.Proxy != "" {
		proxyUrl, _ := url.Parse(in.Proxy)
		Transport.Proxy = http.ProxyURL(proxyUrl)
	}

	client := &http.Client{
		Transport: Transport,
	}

	log_.Info("请求url", in.Url)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log_.Error("请求失败", err)
			return err
		}
		bodyStr := string(body)
		log_.Error("请求失败", bodyStr)
		message := "接口请求失败"
		return error_.NewTextError(message)
	}
	// 处理 HTTP 响应
	reader := bufio.NewReader(resp.Body)

	for {
		var sse SseMessage

		// Read fields until end of message
		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				sse.Event = "end"
				break
			}
			if err != nil {
				log_.Error("读取stream失败", err)
				return err
			}
			if line == "\n" {
				break
			}
			sse.Merge(parseSSE(line))
		}

		// Check for end of stream
		if sse.IsEmpty() {
			return nil
		}
		if sse.Data == "[DONE]" {
			sse.Event = "end"
			sse.Data = ""
		}
		err := stream(&sse)
		if err != nil {
			log_.Error("sse on Message error", err)
			return nil
		}
		if sse.Event == "end" {
			return nil
		}
	}

}

func parseSSE(line string) SseMessage {
	sse := SseMessage{}
	switch {
	case line == "\n":
		return sse
	case line[:3] == "id:":
		sse.ID = line[4 : len(line)-1]
	case line[:5] == "data:":
		sse.Data = line[6 : len(line)-1]
	case line[:6] == "event:":
		sse.Event = line[7 : len(line)-1]
	}
	return sse
}

type SseMessage struct {
	Lines []string
	Event string // 事件名称
	Data  string // 数据内容
	ID    string // 消息 ID
}

func (sse *SseMessage) Merge(other SseMessage) {
	if other.Event != "" {
		sse.Event = other.Event
	}
	if other.Data != "" {
		sse.Data = sse.Data + other.Data
	}
	if other.ID != "" {
		sse.ID = other.ID
	}
}

func (sse *SseMessage) IsEmpty() bool {
	if sse.Event == "" && sse.Data == "" && sse.ID == "" {
		return true
	}
	return false
}
