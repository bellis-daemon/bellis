package palworld

import (
	"context"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/implements"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/option"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/status"
	"github.com/gorcon/rcon"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

type Palworld struct {
	implements.Template
	options palworldOptions
}

type palworldOptions struct {
	Address  string `json:"Address" bson:"Address"`
	Password string `json:"Password" bson:"Password"`
}

func (this *Palworld) Fetch(ctx context.Context) (status.Status, error) {
	conn, err := rcon.Dial(this.options.Address, this.options.Password)
	if err != nil {
		return &palworldStatus{}, err
	}
	defer conn.Close()
	resp, err := conn.Execute("Info")
	if err != nil {
		return &palworldStatus{}, err
	}
	ret := &palworldStatus{}
	var v1Started, v1Ended, nameStarted bool
	for _, char := range resp {
		if char == '[' {
			v1Started = true
			continue
		}
		if v1Started && char != ']' {
			ret.ServerVersion += string(char)
		}
		if char == ']' {
			v1Ended = true
			v1Started = false
			continue
		}
		if v1Ended && char != ' ' && !nameStarted {
			nameStarted = true
		}
		if nameStarted {
			ret.ServerName += string(char)
		}
	}

	if ret.ServerName == "" {
		ret.ServerName = resp
	}
	ret.ServerName = strings.TrimRight(ret.ServerName, "\n")

	if ret.ServerVersion == "" {
		ret.ServerVersion = "Unknown"
	}

	resp2, err := conn.Execute("ShowPlayers")
	if err != nil {
		return &palworldStatus{}, err
	}
	lines := strings.Split(resp2, "\n")
	ret.OnlinePlayer = len(lines) - 2
	return ret, nil
}

type palworldStatus struct {
	OnlinePlayer  int    `json:"OnlinePlayer"`
	ServerVersion string `json:"ServerVersion"`
	ServerName    string `json:"ServerName"`
}

func (p palworldStatus) PullTrigger(triggerName string) *status.TriggerInfo {
	panic("implement me")
}

func init() {
	implements.Register("palworld", func(options bson.M) implements.Implement {
		return &Palworld{options: option.ToOption[palworldOptions](options)}
	})
}
