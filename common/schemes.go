package common

import (
	"github.com/bellis-daemon/bellis/common/cryptoo"
	"github.com/minoic/glgf"
	"os/exec"
)

var Measurements = map[int]string{
	0:  "null",
	1:  "bt",
	2:  "ping",
	3:  "http",
	4:  "minecraft",
	5:  "v2ray",
	6:  "dns",
	7:  "vps",
	8:  "nps",
	9:  "frp",
	10: "html",
	11: "cloudflare",
	12: "pterodactyl",
	13: "qiniu",
	14: "upyun",
	15: "docker",
	16: "source",
	17: "jmx",
	18: "redis",
	19: "elasticsearch",
	20: "oracle",
	21: "postgresql",
	22: "mssql",
	23: "mariadb",
	24: "mongodb",
	25: "ntp",
	26: "smb",
	27: "snmp",
	28: "weblogic",
	29: "webdav",
	30: "ftp",
	31: "ssh",
	32: "vnc",
	33: "rdp",
	34: "rtmp",
	35: "rtsp",
	36: "websocket",
	37: "rss",
	38: "jellyfin",
	39: "plex",
	40: "photoprism",
	41: "synology",
	42: "qnap",
	43: "unraid",
	44: "homeassistant",
	45: "discuz",
	46: "gitlab",
}

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
