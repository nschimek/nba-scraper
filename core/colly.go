package core

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Danny-Dasilva/CycleTLS/cycletls"
	"github.com/gocolly/colly/v2"
)

const (
	allowedDomain = "www.basketball-reference.com"
	domainGlob    = "*" + allowedDomain + "*"
	userAgent     = "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:87.0) Gecko/20100101 Firefox/87.0"
	ja3           = "771,4865-4867-4866-49195-49199-52393-52392-49196-49200-49162-49161-49171-49172-51-57-47-53-10,0-23-65281-10-11-35-16-5-51-43-13-45-28-21,29-23-24-25-256-257,0"
)

type transport struct{}

var LimitRule = &colly.LimitRule{
	DomainGlob:  domainGlob,
	Parallelism: 1,
	RandomDelay: 5 * time.Second,
}

func (t *transport) RoundTrip(r *http.Request) (*http.Response, error) {
	client := cycletls.Init()
	response, err := client.Do(r.URL.String(), cycletls.Options{
		Body:      "",
		Ja3:       ja3,
		UserAgent: userAgent,
	}, "GET")
	if err != nil {
		return nil, err
	}

	return &http.Response{
		Request:       r,
		Header:        http.Header{"Content-Type": {"html"}}, // REQUIRED OR COLLY IGNORES EVERYTHING!
		StatusCode:    response.Status,
		Status:        fmt.Sprintf("%d %s", response.Status, http.StatusText(response.Status)),
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          ioutil.NopCloser(bytes.NewBufferString(response.Body)),
		ContentLength: int64(len(response.Body)),
		Close:         true,
	}, nil
}

func createColly() *colly.Collector {
	c := colly.NewCollector(colly.AllowedDomains(allowedDomain))
	c.WithTransport(&transport{})
	c.Limit(LimitRule)
	return c
}

func CloneColly(colly *colly.Collector) *colly.Collector {
	c := colly.Clone()
	c.OnRequest(onRequestVisit)
	c.OnError(onError)
	return c
}

func onRequestVisit(r *colly.Request) {
	Log.Infof("Visiting: %s", r.URL.String())
}

func onError(r *colly.Response, err error) {
	Log.Info(string(r.Body))
	Log.Fatalf("Scraping resulted in error: %s", err)
}
