package status

import "github.com/bellis-daemon/bellis/common"

type Status interface {
	PullTrigger(triggerName string) *TriggerInfo
}

type TriggerInfo struct {
	Name     string
	Message  string
	Priority common.PriorityLevel
}
