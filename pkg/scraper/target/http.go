package target

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/oleh-ozimok/php-fpm-exporter/pkg/scraper"
	"github.com/pkg/errors"
)

var httpClient = http.Client{
	Timeout: time.Duration(3 * time.Second),
}

type httpTarget struct {
	url *url.URL
}

func NewHTTPTarget(url *url.URL) scraper.Target {
	return &httpTarget{
		url: url,
	}
}

func (t *httpTarget) Scrape() ([]byte, error) {
	response, err := httpClient.Get(t.url.String() + "?json&full")
	if err != nil {
		return nil, errors.Wrap(err, "HTTP request failed")
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.Errorf("unexpected HTTP status: %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read http body")
	}

	return body, nil
}
