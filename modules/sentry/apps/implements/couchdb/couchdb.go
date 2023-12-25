package couchdb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bellis-daemon/bellis/modules/sentry/apps/implements"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/option"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/status"
	"go.mongodb.org/mongo-driver/bson"
)

type CouchDB struct {
	implements.Template
	options couchDBOptions
	client  *http.Client
}

func (this *CouchDB) Fetch(ctx context.Context) (status.Status, error) {
	req, err := http.NewRequest("GET", this.options.URL, nil)
	if err != nil {
		return &couchDBStatus{}, err
	}

	if this.options.Username != "" || this.options.Password != "" {
		req.SetBasicAuth(this.options.Username, this.options.Password)
	}

	response, err := this.client.Do(req)
	if err != nil {
		return &couchDBStatus{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return &couchDBStatus{}, fmt.Errorf("failed to get stats from couchdb: HTTP responded %d", response.StatusCode)
	}

	stats := Stats{}
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&stats); err != nil {
		return &couchDBStatus{}, fmt.Errorf("failed to decode stats from couchdb: HTTP body %q", response.Body)
	}

	// for couchdb 2.0 API changes
	requestTime := metaData{
		Current: stats.Couchdb.RequestTime.Current,
		Sum:     stats.Couchdb.RequestTime.Sum,
		Mean:    stats.Couchdb.RequestTime.Mean,
		Stddev:  stats.Couchdb.RequestTime.Stddev,
		Min:     stats.Couchdb.RequestTime.Min,
		Max:     stats.Couchdb.RequestTime.Max,
	}

	httpdStatusCodesStatus200 := stats.HttpdStatusCodes.Status200
	httpdStatusCodesStatus201 := stats.HttpdStatusCodes.Status201
	httpdStatusCodesStatus202 := stats.HttpdStatusCodes.Status202
	httpdStatusCodesStatus301 := stats.HttpdStatusCodes.Status301
	httpdStatusCodesStatus304 := stats.HttpdStatusCodes.Status304
	httpdStatusCodesStatus400 := stats.HttpdStatusCodes.Status400
	httpdStatusCodesStatus401 := stats.HttpdStatusCodes.Status401
	httpdStatusCodesStatus403 := stats.HttpdStatusCodes.Status403
	httpdStatusCodesStatus404 := stats.HttpdStatusCodes.Status404
	httpdStatusCodesStatus405 := stats.HttpdStatusCodes.Status405
	httpdStatusCodesStatus409 := stats.HttpdStatusCodes.Status409
	httpdStatusCodesStatus412 := stats.HttpdStatusCodes.Status412
	httpdStatusCodesStatus500 := stats.HttpdStatusCodes.Status500
	// check if couchdb2.0 is used
	if stats.Couchdb.HttpdRequestMethods.Get.Value != nil {
		requestTime = stats.Couchdb.RequestTime.Value

		httpdStatusCodesStatus200 = stats.Couchdb.HttpdStatusCodes.Status200
		httpdStatusCodesStatus201 = stats.Couchdb.HttpdStatusCodes.Status201
		httpdStatusCodesStatus202 = stats.Couchdb.HttpdStatusCodes.Status202
		httpdStatusCodesStatus301 = stats.Couchdb.HttpdStatusCodes.Status301
		httpdStatusCodesStatus304 = stats.Couchdb.HttpdStatusCodes.Status304
		httpdStatusCodesStatus400 = stats.Couchdb.HttpdStatusCodes.Status400
		httpdStatusCodesStatus401 = stats.Couchdb.HttpdStatusCodes.Status401
		httpdStatusCodesStatus403 = stats.Couchdb.HttpdStatusCodes.Status403
		httpdStatusCodesStatus404 = stats.Couchdb.HttpdStatusCodes.Status404
		httpdStatusCodesStatus405 = stats.Couchdb.HttpdStatusCodes.Status405
		httpdStatusCodesStatus409 = stats.Couchdb.HttpdStatusCodes.Status409
		httpdStatusCodesStatus412 = stats.Couchdb.HttpdStatusCodes.Status412
		httpdStatusCodesStatus500 = stats.Couchdb.HttpdStatusCodes.Status500
	}

	return &couchDBStatus{
		AuthCacheMisses:   *stats.Couchdb.AuthCacheMisses.Current,
		AuthCacheHits:     *stats.Couchdb.AuthCacheHits.Current,
		DatabaseWrites:    *stats.Couchdb.DatabaseWrites.Current,
		DatabaseReads:     *stats.Couchdb.DatabaseReads.Current,
		OpenDatabases:     *stats.Couchdb.OpenDatabases.Current,
		OpenOsFiles:       *stats.Couchdb.OpenOsFiles.Current,
		RequestTime:       *requestTime.Current,
		Httpd2XXCounts:    *httpdStatusCodesStatus200.Current + *httpdStatusCodesStatus201.Current + *httpdStatusCodesStatus202.Current,
		Httpd3XXCounts:    *httpdStatusCodesStatus301.Current + *httpdStatusCodesStatus304.Current,
		Httpd4XXCounts:    *httpdStatusCodesStatus400.Current + *httpdStatusCodesStatus401.Current + *httpdStatusCodesStatus403.Current + *httpdStatusCodesStatus404.Current + *httpdStatusCodesStatus405.Current + *httpdStatusCodesStatus409.Current + *httpdStatusCodesStatus412.Current,
		Httpd500Counts:    *httpdStatusCodesStatus500.Current,
		HttpdRequests:     *stats.Httpd.Requests.Current,
		HttpdBulkRequests: *stats.Httpd.BulkRequests.Current,
	}, nil
}

type couchDBOptions struct {
	URL      string
	Username string
	Password string
}

type couchDBStatus struct {
	AuthCacheMisses   float64
	AuthCacheHits     float64
	DatabaseWrites    float64
	DatabaseReads     float64
	OpenDatabases     float64
	OpenOsFiles       float64
	RequestTime       float64
	Httpd2XXCounts    float64
	Httpd3XXCounts    float64
	Httpd4XXCounts    float64
	Httpd500Counts    float64
	HttpdRequests     float64
	HttpdBulkRequests float64
}

func (this *couchDBStatus) PullTrigger(triggerName string) *status.TriggerInfo {
	switch triggerName {

	}
	return nil
}

func init() {
	implements.Register("dns", func(options bson.M) implements.Implement {
		return &CouchDB{
			options: option.ToOption[couchDBOptions](options),
			client:  http.DefaultClient,
		}
	})
}

type (
	metaData struct {
		Current *float64 `json:"current"`
		Sum     *float64 `json:"sum"`
		Mean    *float64 `json:"mean"`
		Stddev  *float64 `json:"stddev"`
		Min     *float64 `json:"min"`
		Max     *float64 `json:"max"`
		Value   *float64 `json:"value"`
	}

	oldValue struct {
		Value metaData `json:"value"`
		metaData
	}

	couchdb struct {
		AuthCacheHits       metaData            `json:"auth_cache_hits"`
		AuthCacheMisses     metaData            `json:"auth_cache_misses"`
		DatabaseWrites      metaData            `json:"database_writes"`
		DatabaseReads       metaData            `json:"database_reads"`
		OpenDatabases       metaData            `json:"open_databases"`
		OpenOsFiles         metaData            `json:"open_os_files"`
		RequestTime         oldValue            `json:"request_time"`
		HttpdRequestMethods httpdRequestMethods `json:"httpd_request_methods"`
		HttpdStatusCodes    httpdStatusCodes    `json:"httpd_status_codes"`
	}

	httpdRequestMethods struct {
		Put    metaData `json:"PUT"`
		Get    metaData `json:"GET"`
		Copy   metaData `json:"COPY"`
		Delete metaData `json:"DELETE"`
		Post   metaData `json:"POST"`
		Head   metaData `json:"HEAD"`
	}

	httpdStatusCodes struct {
		Status200 metaData `json:"200"`
		Status201 metaData `json:"201"`
		Status202 metaData `json:"202"`
		Status301 metaData `json:"301"`
		Status304 metaData `json:"304"`
		Status400 metaData `json:"400"`
		Status401 metaData `json:"401"`
		Status403 metaData `json:"403"`
		Status404 metaData `json:"404"`
		Status405 metaData `json:"405"`
		Status409 metaData `json:"409"`
		Status412 metaData `json:"412"`
		Status500 metaData `json:"500"`
	}

	httpd struct {
		BulkRequests             metaData `json:"bulk_requests"`
		Requests                 metaData `json:"requests"`
		TemporaryViewReads       metaData `json:"temporary_view_reads"`
		ViewReads                metaData `json:"view_reads"`
		ClientsRequestingChanges metaData `json:"clients_requesting_changes"`
	}

	Stats struct {
		Couchdb             couchdb             `json:"couchdb"`
		HttpdRequestMethods httpdRequestMethods `json:"httpd_request_methods"`
		HttpdStatusCodes    httpdStatusCodes    `json:"httpd_status_codes"`
		Httpd               httpd               `json:"httpd"`
	}
)
