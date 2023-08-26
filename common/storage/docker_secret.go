package storage

import (
	"github.com/minoic/glgf"
	"os"
)

func Secret(key string) string {
	file, err := os.ReadFile("/run/secrets/" + key)
	if err != nil {
		glgf.Warn(err)
		return ""
	}
	return string(file)
}
