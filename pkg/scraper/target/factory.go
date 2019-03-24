package target

import (
	"net/url"

	"github.com/oleh-ozimok/php-fpm-exporter/pkg/scraper"
	"github.com/pkg/errors"
)

func Create(rawURL string) (scraper.Target, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, errors.Wrap(err, "url parsing error")
	}

	switch u.Scheme {
	case "tcp":
		return NewFastCGITarget(u), nil
	case "http", "https":
		return NewHTTPTarget(u), nil
	}

	return nil, errors.New("target scheme not supported: " + u.Scheme)
}
