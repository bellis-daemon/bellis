package docker

import (
	"context"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/implements"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/option"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/status"
	"github.com/moby/moby/api/types"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"net/http"
	"time"
)

type Docker struct {
	options dockerOptions
}

func (this *Docker) Fetch(ctx context.Context) (status.Status, error) {
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

type dockerStatus struct {
	Info       types.Info
	Containers []types.Container
	Images     []types.ImageSummary
	Networks   []types.NetworkResource
	Plugins    types.PluginsListResponse
}

func (this *dockerStatus) PullTrigger(triggerName string) *status.TriggerInfo {
	switch triggerName {

	}
	return nil
}

type dockerOptions struct {
	URL string
}

func init() {
	implements.Register("docker", func(options bson.M) implements.Implement {
		return &Docker{options: option.ToOption[dockerOptions](options)}
	})
}
