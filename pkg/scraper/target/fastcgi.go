package target

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/oleh-ozimok/php-fpm-exporter/pkg/scraper"
	"github.com/pkg/errors"
	"github.com/tomasen/fcgi_client"
)

type fastCGITarget struct {
	url *url.URL
}

func NewFastCGITarget(url *url.URL) scraper.Target {
	return &fastCGITarget{
		url: url,
	}
}

func (t *fastCGITarget) Scrape() ([]byte, error) {
	client, err := fcgiclient.DialTimeout(t.url.Scheme, t.url.Host, 3*time.Second)
	if err != nil {
		return nil, errors.Wrap(err, "fastcgi dial failed")
	}

	defer client.Close()

	envs := map[string]string{
		"SCRIPT_FILENAME": t.url.Path,
		"SCRIPT_NAME":     t.url.Path,
		"SERVER_SOFTWARE": "go / php-fpm-exporter",
		"QUERY_STRING":    "json&full",
	}

	response, err := client.Get(envs)
	if err != nil {
		return nil, errors.Wrap(err, "fastcgi get failed")
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != 0 {
		return nil, errors.Errorf("unexpected status: %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read fastcgi body")
	}

	return body, nil
}
