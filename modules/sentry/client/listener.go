package client

import (
	"bytes"
	"encoding/json"
	"github.com/bellis-daemon/bellis/common/storage"
	"net/http"
	"time"

	"github.com/bellis-daemon/bellis/common"
	"github.com/minoic/glgf"
)

type SentrySingletonEvent struct {
	EventType    string         `json:"EventType"`
	EventContent map[string]any `json:"EventContent"`
}

type SentrySingletonEvents []SentrySingletonEvent

func ServeHttpEventListener(token string) {
	ticker := time.NewTicker(5 * time.Second)
	cl := http.DefaultClient

	url := storage.Config().WebEndpoint + "/api/sentry-singleton/refresh"
	for range ticker.C {
		go func() {
			var reqEvents SentrySingletonEvents
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(reqEvents)
			if err != nil {
				glgf.Error(err)
				return
			}
			req, err := http.NewRequest("GET", url, &buf)
			if err != nil {
				glgf.Error(err)
				return
			}
			req.Header.Add("Request-Token", token)
			resp, err := cl.Do(req)
			if err != nil {
				if resp.StatusCode == http.StatusUnauthorized {
					panic("Invalid token")
				} else {
					glgf.Error(err)
					return
				}
			}
			var respEvents SentrySingletonEvents
			err = json.NewDecoder(resp.Body).Decode(&respEvents)
			if err != nil {
				glgf.Error(err)
				return
			}
			for _, event := range respEvents {
				switch event.EventType {
				case common.EntityClaim:
				case common.EntityUpdate:
				case common.EntityDelete:
					// todo: implement event handlers
				}
			}
		}()
	}
}
