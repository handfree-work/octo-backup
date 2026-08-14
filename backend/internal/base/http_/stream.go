package http_

import (
	"bufio"
	"bytes"
	"handfree-work/web-restic/base/log_"
	"io"
	"net/http"
	"net/url"

	"github.com/gogf/gf/v2/encoding/gjson"
)

func RequestStream(in *Request) (streamReader *StreamReader, err error) {
	jsonStr, ok := in.Body.(string)
	if !ok {
		jsonStr, err = gjson.EncodeString(in.Body)
	}

	if err != nil {
		return nil, err
	}
	// 创建 HTTP 请求
	req, err := http.NewRequest(in.Method, in.Url, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		return nil, err
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
		return nil, err
	}

	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log_.Error("请求失败", err)
			return nil, err
		}
		bodyStr := string(body)
		log_.Error("请求失败", bodyStr)
		return &StreamReader{
			IsFinished: true,
			Response:   resp,
			RequestError: &RequestError{
				Body:   bodyStr,
				Status: resp.StatusCode,
			},
		}, nil
	}
	// 处理 HTTP 响应
	reader := bufio.NewReader(resp.Body)

	sr := &StreamReader{
		IsFinished: false,
		Reader:     reader,
		Response:   resp,
	}
	return sr, nil
}

type RequestError struct {
	Body   string
	Status int
}
type StreamReader struct {
	IsFinished   bool
	Reader       *bufio.Reader
	Response     *http.Response
	RequestError *RequestError
}

func (c *StreamReader) Recv() (data *SseMessage, err error) {
	if c.IsFinished {
		err = io.EOF
		return
	}

	data, err = c.processLines()
	return
}

func (c *StreamReader) processLines() (data *SseMessage, err error) {
	for {
		var sse SseMessage

		// Read fields until end of message
		for {
			line, err := c.Reader.ReadString('\n')
			if err == io.EOF {
				c.IsFinished = true
				return nil, io.EOF
			}
			if err != nil {
				log_.Error("读取stream失败", err)
				return nil, err
			}
			sse.Lines = append(sse.Lines, line)
			if line == "\n" {
				break
			}
			sse.Merge(parseSSE(line))
		}

		// Check for end of stream
		if sse.IsEmpty() {
			return nil, nil
		}
		if sse.Data == "[DONE]" {
			c.IsFinished = true
		}
		return &sse, nil
	}
}

func (stream *StreamReader) Close() {
	stream.Response.Body.Close()
}

type OnStream func(data *SseMessage) error

func RequestStreamOn(in *Request, onStream OnStream) (reqErr *RequestError, streamErr error) {
	streamReader, err := RequestStream(in)
	if err != nil {
		log_.Error("请求失败", err)
		return &RequestError{
			Status: 400,
			Body:   err.Error(),
		}, nil
	}
	if streamReader.RequestError != nil {
		return streamReader.RequestError, nil
	}
	defer streamReader.Close()
	for {
		if streamReader.IsFinished {
			return nil, nil
		}
		data, streamErr := streamReader.Recv()
		if streamErr == io.EOF || streamErr == io.ErrUnexpectedEOF {
			return nil, nil
		}
		if streamErr != nil {
			return nil, streamErr
		}
		streamErr = onStream(data)
		if streamErr != nil {
			return nil, streamErr
		}
	}
}
