package implements

import (
	"bytes"
	"errors"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/status"
	"golang.org/x/net/context"
	"net"
)

type DNS struct {
	options dnsOptions
}

func (this *DNS) Fetch(ctx context.Context) (status.Status, error) {
	var buf bytes.Buffer
	switch this.options.Method {
	case "IP":
		ips, err := net.LookupIP(this.options.DomainName)
		if err != nil {
			return &dnsStatus{}, err
		} else {
			for _, ip := range ips {
				buf.WriteString(ip.String() + " ")
			}
			return &dnsStatus{
				Results: buf.String(),
			}, nil
		}
	case "NS":
		nss, err := net.LookupNS(this.options.DomainName)
		if err != nil {
			return &dnsStatus{}, err
		} else {
			for _, ns := range nss {
				buf.WriteString(ns.Host + " ")
			}
			return &dnsStatus{
				Results: buf.String(),
			}, nil
		}
	case "MX":
		mxs, err := net.LookupMX(this.options.DomainName)
		if err != nil {
			return &dnsStatus{}, err
		} else {
			for _, mx := range mxs {
				buf.WriteString(mx.Host + " ")
			}
			return &dnsStatus{
				Results: buf.String(),
			}, nil
		}
	case "TXT":
		texts, err := net.LookupTXT(this.options.DomainName)
		if err != nil {
			return &dnsStatus{}, err
		} else {
			for _, text := range texts {
				buf.WriteString(text + " ")
			}
			return &dnsStatus{
				Results: buf.String(),
			}, nil
		}
	case "CNAME":
		cname, err := net.LookupCNAME(this.options.DomainName)
		if err != nil {
			return &dnsStatus{}, err
		} else {
			return &dnsStatus{
				Results: cname,
			}, nil
		}
	case "ADDR":
		names, err := net.LookupAddr(this.options.DomainName)
		if err != nil {
			return &dnsStatus{}, err
		} else {
			for _, name := range names {
				buf.WriteString(name + " ")
			}
			return &dnsStatus{
				Results: buf.String(),
			}, nil
		}
	case "SRV":
		cname, srvs, err := net.LookupSRV(this.options.SRVService, this.options.SRVProtocol, this.options.DomainName)
		if err != nil {
			return &dnsStatus{}, err
		} else {
			buf.WriteString(cname + " ")
			for _, srv := range srvs {
				buf.WriteString(srv.Target + " ")
			}
			return &dnsStatus{
				Results: buf.String(),
			}, nil
		}
	default:
		return &dnsStatus{}, errors.New("错误的解析模式：" + this.options.Method)
	}
}

func (this *DNS) Init(setOptions func(options any) error) error {
	return setOptions(&this.options)
}

type dnsStatus struct {
	Results string `json:"results"`
}

func (this *dnsStatus) PullTrigger(triggerName string) *status.TriggerInfo {
	switch triggerName {

	}
	return nil
}

type dnsOptions struct {
	DomainName  string `json:"domain_name"`
	Method      string `json:"method"`
	SRVService  string `json:"srv_service"`
	SRVProtocol string `json:"srv_protocol"`
}
