package storage

import "os"

func Secret(key string) string {
	file, err := os.ReadFile("/run/secrets/" + key)
	if err != nil {
		panic(err)
	}
	return string(file)
}
