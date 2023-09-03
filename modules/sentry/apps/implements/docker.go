package implements

import (
	"context"
	"github.com/moby/moby/api/types"
	"io"
	"net/http"
	"time"
)

type Docker struct {
	options dockerOptions
}

func (this *Docker) Fetch(ctx context.Context) (any, error) {
	req, err := http.NewRequest("GET", this.options.URL, nil)
	if err != nil {
		return &dockerStatus{}, err
	}
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return &dockerStatus{}, err
	}
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return &dockerStatus{}, err
	}
	// todo: 修改docker状态格式
	return &dockerStatus{}, err
}

func (this *Docker) Init(setOptions func(options any) error) error {
	return setOptions(&this.options)
}

type dockerStatus struct {
	Info       types.Info
	Containers []types.Container
	Images     []types.ImageSummary
	Networks   []types.NetworkResource
	Plugins    types.PluginsListResponse
}

type dockerOptions struct {
	URL string
}
