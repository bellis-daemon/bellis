package geo

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

//go:embed ip2region.xdb
var xdbFile []byte

type GeoData struct {
	IP       string `json:"query"`
	Status   string `json:"status"`
	Country  string `json:"country"`
	Region   string `json:"regionName"`
	City     string `json:"city"`
	Timezone string `json:"timezone"`
	ISP      string `json:"isp"`
}

func (this *GeoData) String() string {
	if this.Status != "success" {
		return ""
	}
	var buf bytes.Buffer
	var count = 0
	if this.IP != "" && this.IP != "0" {
		if count != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(this.IP)
		count++
	}
	if this.Country != "" && this.Country != "0" {
		if count != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(this.Country)
		count++
	}
	if this.Region != "" && this.Region != "0" {
		if count != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(this.Region)
		count++
	}
	if this.City != "" && this.City != "0" {
		if count != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(this.City)
		count++
	}
	if this.ISP != "" && this.ISP != "0" {
		if count != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(this.ISP)
		count++
	}
	return buf.String()
}

func isIPV4(ipAddress string) bool {
	if ipAddress == "" {
		return false
	}
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return false
	}
	ipv4 := ip.To4()
	if ipv4 == nil {
		return false
	}
	return true
}

func FromAPI(address ...string) (*GeoData, error) {
	url := "http://ip-api.com/json/?fields=25369"
	if len(address) != 0 {
		if !isIPV4(address[0]) {
			return nil, errors.New("invalid ipv4 address:" + address[0])
		}
		url = fmt.Sprintf("http://ip-api.com/json/%s?fields=25369", address[0])
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	var g GeoData
	err = json.NewDecoder(resp.Body).Decode(&g)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func FromLocal(address string) (*GeoData, error) {
	if !isIPV4(address) {
		return nil, errors.New("invalid ipv4 address:" + address)
	}
	searcher, err := xdb.NewWithBuffer(xdbFile)
	if err != nil {
		return nil, err
	}
	result, err := searcher.SearchByStr(address)
	if err != nil {
		return nil, err
	}
	splits := strings.Split(result, "|")
	if len(splits) < 5 {
		return nil, fmt.Errorf("failed to search region data, result: %s", result)
	}
	return &GeoData{
		IP:       address,
		Status:   "success",
		Country:  splits[0],
		Region:   splits[2],
		City:     splits[3],
		Timezone: "",
		ISP:      splits[4],
	}, nil
}

var self *GeoData

func Self() (*GeoData, error) {
	if self == nil {
		var err error
		self, err = FromAPI()
		if err != nil {
			self = nil
			return nil, err
		}
	}
	return self, nil
}
