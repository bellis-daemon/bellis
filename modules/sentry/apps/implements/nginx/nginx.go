package nginx

import (
	"bufio"
	"context"
	"fmt"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/implements"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/option"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/status"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Nginx struct {
	options nginxOptions
	client  *http.Client
}

func (this *Nginx) Fetch(ctx context.Context) (status.Status, error) {
	addr, err := url.Parse(this.options.Url)
	if err != nil {
		return &nginxStatus{}, fmt.Errorf("error parsing url: %s", this.options.Url)
	}

	resp, err := this.client.Get(addr.String())
	if err != nil {
		return &nginxStatus{}, fmt.Errorf("error making HTTP request to %q: %w", addr.String(), err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return &nginxStatus{}, fmt.Errorf("%s returned HTTP status %s", addr.String(), resp.Status)
	}
	r := bufio.NewReader(resp.Body)

	// Active connections
	_, err = r.ReadString(':')
	if err != nil {
		return &nginxStatus{}, err
	}
	line, err := r.ReadString('\n')
	if err != nil {
		return &nginxStatus{}, err
	}
	active, err := strconv.ParseUint(strings.TrimSpace(line), 10, 64)
	if err != nil {
		return &nginxStatus{}, err
	}

	// Server accepts handled requests
	_, err = r.ReadString('\n')
	if err != nil {
		return &nginxStatus{}, err
	}
	line, err = r.ReadString('\n')
	if err != nil {
		return &nginxStatus{}, err
	}
	data := strings.Fields(line)
	accepts, err := strconv.ParseUint(data[0], 10, 64)
	if err != nil {
		return &nginxStatus{}, err
	}

	handled, err := strconv.ParseUint(data[1], 10, 64)
	if err != nil {
		return &nginxStatus{}, err
	}
	requests, err := strconv.ParseUint(data[2], 10, 64)
	if err != nil {
		return &nginxStatus{}, err
	}

	// Reading/Writing/Waiting
	line, err = r.ReadString('\n')
	if err != nil {
		return &nginxStatus{}, err
	}
	data = strings.Fields(line)
	reading, err := strconv.ParseUint(data[1], 10, 64)
	if err != nil {
		return &nginxStatus{}, err
	}
	writing, err := strconv.ParseUint(data[3], 10, 64)
	if err != nil {
		return &nginxStatus{}, err
	}
	waiting, err := strconv.ParseUint(data[5], 10, 64)
	if err != nil {
		return &nginxStatus{}, err
	}
	return &nginxStatus{
		Active:   active,
		Accepts:  accepts,
		Handled:  handled,
		Requests: requests,
		Reading:  reading,
		Writing:  writing,
		Waiting:  waiting,
	}, nil
}

type nginxOptions struct {
	Url string `json:"Url"`
}

type nginxStatus struct {
	Active   uint64 // Active connections
	Accepts  uint64 // Server accepts handled requests
	Handled  uint64
	Requests uint64
	Reading  uint64
	Writing  uint64
	Waiting  uint64
}

func (this *nginxStatus) PullTrigger(triggerName string) *status.TriggerInfo {
	switch triggerName {

	}
	return nil
}

func init() {
	implements.Register("nginx", func(options bson.M) implements.Implement {
		return &Nginx{options: option.ToOption[nginxOptions](options), client: http.DefaultClient}
	})
}
