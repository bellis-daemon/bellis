package cryptoo

import (
	"fmt"
	"time"
)

func FormatDuration(duration time.Duration) string {
	years := duration / (365 * 24 * time.Hour)
	duration -= years * 365 * 24 * time.Hour

	months := duration / (30 * 24 * time.Hour)
	duration -= months * 30 * 24 * time.Hour

	days := duration / (24 * time.Hour)
	duration -= days * 24 * time.Hour

	hours := duration / time.Hour
	duration -= hours * time.Hour

	minutes := duration / time.Minute

	// 只保留两个有值且最大的数量级
	units := []struct {
		Value  time.Duration
		Suffix string
	}{
		{years, "y"},
		{months, "M"},
		{days, "d"},
		{hours, "h"},
		{minutes, "m"},
	}

	var result string
	count := 0
	for _, unit := range units {
		if unit.Value > 0 {
			result += fmt.Sprintf("%d%s", unit.Value, unit.Suffix)
			count++
		}
		if count == 2 {
			break
		}
	}

	return result
}
