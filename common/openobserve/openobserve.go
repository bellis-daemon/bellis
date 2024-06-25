package openobserve

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bellis-daemon/bellis/common"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/minoic/glgf"
)

type openObserveWriter struct {
	level string
	c     chan []byte
}

func (this *openObserveWriter) Write(p []byte) (n int, err error) {
	this.c <- p
	return len(p), nil
}

type openObserve struct {
	org      string
	username string
	password string
	stream   string
	client   *http.Client
	writers  []openObserveWriter
}

func (this *openObserve) send(logs []map[string]any) error {
	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(logs)
	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.openobserve.ai/api/%s/%s/_json", this.org, this.stream), &buf)
	if err != nil {
		return err
	}
	req.SetBasicAuth(this.username, this.password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")
	_, err = this.client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func (this *openObserve) run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := make(chan map[string]any, 16)
	for i := range this.writers {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case msg := <-this.writers[i].c:
					c <- map[string]any{
						"level":    this.writers[i].level,
						"hostname": common.Hostname(),
						"log":      string(msg),
					}
				}
			}
		}()
	}
	var list []map[string]any
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			if len(list) != 0 {
				err := this.send(list)
				if err != nil {
					fmt.Println(err)
				} else {
					list = []map[string]any{}
				}
			}
		case msg := <-c:
			list = append(list, msg)
			if len(list) >= 10 {
				err := this.send(list)
				if err != nil {
					fmt.Println(err)
				} else {
					list = []map[string]any{}
				}
			}
		}
	}
}

func RegisterGlgf() {
	instance := &openObserve{
		org:      storage.Config().OpenObserveOrg,
		username: storage.Config().OpenObserveUsername,
		password: storage.Config().OpenObservePassword,
		stream:   common.AppName,
		client:   http.DefaultClient,
		writers: []openObserveWriter{
			{
				glgf.DEBG.String(),
				make(chan []byte),
			},
			{
				glgf.INFO.String(),
				make(chan []byte),
			},
			{
				glgf.WARN.String(),
				make(chan []byte),
			},
			{
				glgf.ERR.String(),
				make(chan []byte),
			},
			{
				glgf.OK.String(),
				make(chan []byte),
			},
		},
	}
	glgf.Get().
		SetMode(glgf.BOTH).
		AddLevelWriter(glgf.DEBG, &instance.writers[0]).
		AddLevelWriter(glgf.INFO, &instance.writers[1]).
		AddLevelWriter(glgf.WARN, &instance.writers[2]).
		AddLevelWriter(glgf.ERR, &instance.writers[3]).
		AddLevelWriter(glgf.OK, &instance.writers[4])
	go instance.run()
}
