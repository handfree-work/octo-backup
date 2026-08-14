package http_

import (
	"context"
	"handfree-work/web-restic/internal/base/error_"
	"net/http"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
)

type Request struct {
	Url     string
	Method  string
	Body    interface{}
	Headers map[string]string
	// Query   *map[string]string
	Proxy string
}
type Response struct {
	Body   string
	Status int
}

func DoRequest(in *Request) (res *Response, err error) {

	client := g.Client()
	if in.Proxy != "" {
		client.SetProxy(in.Proxy)
	}
	if in.Headers != nil {
		client.SetHeaderMap(in.Headers)
	}

	if in.Method == "" {
		in.Method = "POST"
	}

	url := in.Url
	var r *gclient.Response
	if in.Method == "POST" {
		jsonStr, ok := in.Body.(string)
		if !ok {
			jsonStr, err = gjson.EncodeString(in.Body)
			if err != nil {
				return nil, err
			}
		}
		r, err = client.Post(context.Background(), url, jsonStr)
	} else if in.Method == "GET" {
		r, err = client.Get(context.Background(), url)
	}

	if err != nil {
		return nil, err
	}

	if r.StatusCode != http.StatusOK {
		return nil, error_.NewTextError("请求失败")
	}
	defer r.Close()
	body := r.ReadAllString()
	return &Response{
		Body:   body,
		Status: r.StatusCode,
	}, nil

}
