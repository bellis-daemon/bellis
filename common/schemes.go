package common

import (
	"github.com/bellis-daemon/bellis/common/cryptoo"
	"github.com/minoic/glgf"
	"os/exec"
)

var (
	hostname  string
	BuildTime string
	GoVersion string
)

func Hostname() string {
	if hostname == "" {
		b, err := exec.Command("sh", "hostname.sh").Output()
		if err != nil {
			glgf.Error(err)
		}
		hostname = string(b)
		if hostname == "" {
			hostname = "RND" + cryptoo.RandString(5)
		}
		if len(hostname) > 8 {
			hostname = hostname[:8]
		}
	}
	return hostname
}

// Redis stream
const (
	EntityOfflineAlert = "EntityOfflineAlert"
	EntityOnlineAlert  = "EntityOnlineAlert"
	EntityClaim        = "EntityClaim"
	EntityUpdate       = "EntityUpdate"
	EntityDelete       = "EntityDelete"
	CaptchaToEmail     = "CaptchaToEmail"
)

type PriorityLevel int32

const (
	Predicting = 4000 + iota
	Warning
	Critical
)
