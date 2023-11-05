package http

import (
	"context"
	"errors"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/implements"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/status"
	"net/http"
	"net/http/httptrace"
	"strings"
	"time"
)

type HTTP struct {
	options httpOptions
}

func (this *HTTP) Fetch(ctx context.Context) (status.Status, error) {
	req, err := http.NewRequest(this.options.Method, this.options.URL, nil)
	if err != nil {
		return &httpStatus{}, err
	}
	var addr string
	trace := &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			addr = connInfo.Conn.RemoteAddr().String()
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return &httpStatus{}, err
	}
	if resp.StatusCode >= 400 {
		return &httpStatus{}, errors.New(resp.Status)
	}
	ret := &httpStatus{
		Status: resp.Status,
	}
	splits := strings.Split(addr, ":")
	if len(splits) > 0 {
		ret.IP = splits[0]
	}
	return ret, nil
}

func (this *HTTP) Init(setOptions func(options any) error) error {
	return setOptions(&this.options)
}

type httpStatus struct {
	IP     string `json:"ip"`
	Status string `json:"status"`
}

func (h httpStatus) PullTrigger(triggerName string) *status.TriggerInfo {
	switch triggerName {

	}
	return nil
}

type httpOptions struct {
	URL    string `json:"url"`
	Method string `json:"method"`
}

func init() {
	implements.Add("http", func() implements.Implement {
		return &HTTP{}
	})
}
