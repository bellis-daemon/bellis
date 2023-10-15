package models

import "time"

type Timezone string

const DefaultTimezone = "Asia/Shanghai"

func (this Timezone) Location() *time.Location {
	s := string(this)
	if this == "" {
		s = DefaultTimezone
	}
	loc, _ := time.LoadLocation(s)
	return loc
}
