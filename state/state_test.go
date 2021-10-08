package state

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/camptocamp/terraboard/config"
	"google.golang.org/api/option"
	raw "google.golang.org/api/storage/v1"
)

func TestConfigure(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.Header().Set("X-RateLimit-Limit", "30")
		w.Header().Set("TFP-API-Version", "34.21.9")
		w.WriteHeader(204)
	}))
	defer ts.Close()

	gotURL := make(chan *url.URL, 1)
	hClient, close := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(ioutil.Discard, r.Body)
		gotURL <- r.URL
		if strings.Contains(r.URL.String(), "/rewriteTo/") {
			res := &raw.RewriteResponse{Done: true}
			bytes, err := res.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}
			w.Write(bytes)
		} else {
			fmt.Fprintf(w, "{}")
		}
	})
	defer close()
	option.WithHTTPClient(hClient)

	config := config.Config{
		AWS: []config.AWSConfig{
			{
				AccessKey:       "AKIAIOSFODNN7EXAMPLE",
				SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
				Region:          "us-east-1",
				Endpoint:        "http://localhost:8000",
				DynamoDBTable:   "test-locks",
				S3: []config.S3BucketConfig{
					{
						Bucket:        "test",
						FileExtension: []string{".tfstate"},
					},
				},
			},
		},
		TFE: []config.TFEConfig{
			{
				Address: ts.URL,
				Token:   "abcd1234",
			},
		},
		GCP: []config.GCPConfig{
			{
				HTTPClient: hClient,
				GCSBuckets: []string{
					"test-bucket",
				},
			},
		},
		Gitlab: []config.GitlabConfig{
			{
				Address: "http://localhost:8081",
				Token:   "test-token",
			},
		},
	}

	providers, err := Configure(&config)
	if err != nil {
		t.Error(err)
	} else if len(providers) != 4 {
		t.Errorf("Expected 4 providers, got %d", len(providers))
	}
}

func newTestServer(handler func(w http.ResponseWriter, r *http.Request)) (*http.Client, func()) {
	ts := httptest.NewTLSServer(http.HandlerFunc(handler))
	tlsConf := &tls.Config{InsecureSkipVerify: true}
	tr := &http.Transport{
		TLSClientConfig: tlsConf,
		DialTLS: func(netw, addr string) (net.Conn, error) {
			return tls.Dial("tcp", ts.Listener.Addr().String(), tlsConf)
		},
	}
	return &http.Client{Transport: tr}, func() {
		tr.CloseIdleConnections()
		ts.Close()
	}
}
