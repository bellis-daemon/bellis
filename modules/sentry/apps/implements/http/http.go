package http

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"net/http/httptrace"
	"strings"
	"time"

	"github.com/bellis-daemon/bellis/modules/sentry/apps/implements"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/option"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/status"
	"go.mongodb.org/mongo-driver/bson"
)

var httpClient = http.Client{
	Timeout: 6 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:    100,
		IdleConnTimeout: 30 * time.Second,
	},
}

type HTTP struct {
	implements.Template
	options  httpOptions
	tlsState *tls.ConnectionState
}

func (this *HTTP) Fetch(ctx context.Context) (status.Status, error) {
	ret := &httpStatus{}
	req, err := http.NewRequestWithContext(ctx, this.options.Method, this.options.URL, nil)
	if err != nil {
		return ret, err
	}
	trace := &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			splits := strings.Split(connInfo.Conn.RemoteAddr().String(), ":")
			if len(splits) > 0 {
				ret.IP = splits[0]
			}
		},
		TLSHandshakeDone: func(state tls.ConnectionState, err error) {
			if this.tlsState == nil && strings.Contains(this.options.URL, "https://") {
				this.tlsState = &state
			}
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	resp, err := httpClient.Do(req)
	if err != nil {
		return ret, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return ret, errors.New(resp.Status)
	}
	ret.StatusCode = resp.StatusCode
	ret.ContentType = resp.Header.Get("Content-Type")
	ret.ContentLength = resp.ContentLength
	if ret.ContentLength == -1 {
		length, _ := io.CopyN(io.Discard, resp.Body, 100000)
		ret.ContentLength = length
	}
	if this.tlsState != nil {
		switch this.tlsState.Version {
		case tls.VersionTLS10:
			ret.TLSVersion = "TLS10"
		case tls.VersionTLS11:
			ret.TLSVersion = "TLS11"
		case tls.VersionTLS12:
			ret.TLSVersion = "TLS12"
		case tls.VersionTLS13:
			ret.TLSVersion = "TLS13"
		default:
			ret.TLSVersion = "None"
		}
		if len(this.tlsState.PeerCertificates) > 0 {
			ret.TLSStartTime = this.tlsState.PeerCertificates[0].NotBefore.Format(time.RFC3339Nano)
			ret.TLSExpireTime = this.tlsState.PeerCertificates[0].NotAfter.Format(time.RFC3339Nano)
			if len(this.tlsState.PeerCertificates[0].Issuer.Organization) > 0 {
				ret.TLSIssuer = this.tlsState.PeerCertificates[0].Issuer.Organization[0]
			}
		}
	}
	return ret, nil
}

type httpStatus struct {
	IP            string
	StatusCode    int
	ContentLength int64
	ContentType   string
	TLSVersion    string
	TLSStartTime  string
	TLSExpireTime string
	TLSIssuer     string
}

func (h httpStatus) PullTrigger(triggerName string) *status.TriggerInfo {
	switch triggerName {

	}
	return nil
}

type httpOptions struct {
	URL    string `json:"url"`
	Method string `json:"method"`
}

func init() {
	implements.Register("http", func(options bson.M) implements.Implement {
		o := option.ToOption[httpOptions](options)
		if o.Method == "" {
			o.Method = "GET"
		}
		if !strings.Contains(o.URL, "http://") && !strings.Contains(o.URL, "https://") {
			o.URL = "http://" + o.URL
		}
		return &HTTP{options: o}
	})
}
