package http

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/implements"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/status"
	"net/http"
	"net/http/httptrace"
	"strings"
	"time"
)

type HTTP struct {
	options httpOptions
}

func (this *HTTP) Fetch(ctx context.Context) (status.Status, error) {
	ret := &httpStatus{}
	req, err := http.NewRequest(this.options.Method, this.options.URL, nil)
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
			switch state.Version {
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
			if len(state.PeerCertificates) > 0 {
				ret.TLSStartTime = state.PeerCertificates[0].NotBefore.Format(time.RFC3339Nano)
				ret.TLSExpireTime = state.PeerCertificates[0].NotAfter.Format(time.RFC3339Nano)
				if len(state.PeerCertificates[0].Issuer.Organization) > 0 {
					ret.TLSIssuer = state.PeerCertificates[0].Issuer.Organization[0]
				}
			}
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return ret, err
	}
	if resp.StatusCode >= 400 {
		return ret, errors.New(resp.Status)
	}
	ret.StatusCode = resp.StatusCode
	ret.ContentLength = resp.ContentLength
	return ret, nil
}

func (this *HTTP) Init(setOptions func(options any) error) error {
	return setOptions(&this.options)
}

type httpStatus struct {
	IP            string
	StatusCode    int
	ContentLength int64
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
	implements.Add("http", func() implements.Implement {
		return &HTTP{}
	})
}
