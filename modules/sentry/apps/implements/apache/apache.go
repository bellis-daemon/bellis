package apache

import (
	"bufio"
	"context"
	"fmt"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/implements"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/option"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/status"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"net/url"
	"strings"
)

type Apache struct {
	options apacheOptions
	client  *http.Client
}

func (this *Apache) Fetch(ctx context.Context) (status.Status, error) {
	addr, err := url.Parse(this.options.Url)
	if err != nil {
		return &apacheStatus{}, fmt.Errorf("unable to parse address %q: %w", this.options.Url, err)
	}
	req, err := http.NewRequestWithContext(ctx, "GET", addr.String(), nil)
	if err != nil {
		return &apacheStatus{}, fmt.Errorf("error on new request to %q: %w", addr.String(), err)
	}

	if len(this.options.Username) != 0 && len(this.options.Password) != 0 {
		req.SetBasicAuth(this.options.Username, this.options.Password)
	}
	resp, err := this.client.Do(req)
	if err != nil {
		return &apacheStatus{}, fmt.Errorf("error on request to %q: %w", addr.String(), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &apacheStatus{}, fmt.Errorf("%s returned HTTP status %s", addr.String(), resp.Status)
	}
	sc := bufio.NewScanner(resp.Body)
	ret := &apacheStatus{}
	for sc.Scan() {
		line := sc.Text()
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			key, part := strings.ReplaceAll(parts[0], " ", ""), strings.TrimSpace(parts[1])
			switch key {
			case "BusyWorkers":
				ret.BusyWorkers = cast.ToFloat32(part)
			case "IdleWorkers":
				ret.IdleWorkers = cast.ToFloat32(part)
			case "CPULoad":
				ret.CPULoad = cast.ToFloat32(part)
			case "ReqPerSec":
				ret.ReqPerSec = cast.ToFloat32(part)
			case "BytesPerSec":
				ret.BytesPerSec = cast.ToFloat32(part)
			case "Uptime":
				ret.Uptime = cast.ToFloat32(part)
			case "ConnsTotal":
				ret.ConnsTotal = cast.ToFloat32(part)
			case "TotalAccesses":
				ret.TotalAccesses = cast.ToFloat32(part)
			case "TotalkBytes":
				ret.TotalkBytes = cast.ToFloat32(part)
			default:
				continue
			}
		}
	}
	return ret, nil
}

type apacheOptions struct {
	// readable version of the mod_status page including the auto query string.
	// Default is "http://localhost/server-status?auto".
	Url string `json:"Url"`
	// Credentials for basic HTTP authentication.
	Username string `json:"Username"`
	Password string `json:"Password"`
}

type apacheStatus struct {
	BusyWorkers   float32 `json:"BusyWorkers"`
	IdleWorkers   float32 `json:"IdleWorkers"`
	CPULoad       float32 `json:"CPULoad"`
	ReqPerSec     float32 `json:"ReqPerSec"`
	BytesPerSec   float32 `json:"BytesPerSec"`
	Uptime        float32 `json:"Uptime"`
	ConnsTotal    float32 `json:"ConnsTotal"`
	TotalAccesses float32 `json:"TotalAccesses"`
	TotalkBytes   float32 `json:"TotalkBytes"`
}

func (this *apacheStatus) PullTrigger(triggerName string) *status.TriggerInfo {
	switch triggerName {

	}
	return nil
}

func init() {
	implements.Register("apache", func(options bson.M) implements.Implement {
		return &Apache{
			options: option.ToOption[apacheOptions](options),
			client:  http.DefaultClient,
		}
	})
}
