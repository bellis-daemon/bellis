package http

import (
	"context"
	"errors"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/status"
	"net/http"
	"time"
)

type HTTP struct {
	options httpOptions
}

func (this *HTTP) Fetch(ctx context.Context) (status.Status, error) {
	method := "GET"
	if this.options.Method != "" {
		method = this.options.Method
	}
	req, err := http.NewRequest(method, this.options.URL, nil)
	if err != nil {
		return &httpStatus{}, err
	}
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return &httpStatus{}, err
	}
	if resp.StatusCode >= 400 {
		return &httpStatus{}, errors.New(resp.Status)
	}
	return &httpStatus{
		IP:     resp.Request.Host,
		Status: resp.Status,
	}, nil
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
