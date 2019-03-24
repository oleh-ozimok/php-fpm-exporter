package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oleh-ozimok/php-fpm-exporter/pkg/scraper"
	"github.com/oleh-ozimok/php-fpm-exporter/pkg/scraper/target"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

type options struct {
	address string
}

func main() {
	options := &options{}
	command := &cobra.Command{
		Use: "php-fpm-exporter",
		RunE: func(cmd *cobra.Command, args []string) error {
			return options.Run()
		},
	}

	command.Flags().StringVar(&options.address, "address", ":9253", "The address to listen on for HTTP requests")

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

func (o *options) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", metricsHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<html>
    		<head><title>PHP-FPM Exporter</title></head>
    		<body>
    		<h1>PHP-FPM Exporter</h1>
    		</body>
    		</html>`,
		))
	})

	server := http.Server{
		Addr:    o.address,
		Handler: mux,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	var g errgroup.Group

	g.Go(server.ListenAndServe)
	g.Go(func() error {
		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		return server.Shutdown(ctx)
	})

	return g.Wait()
}

func metricsHandler(responseWriter http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	request = request.WithContext(ctx)

	targetURL := request.URL.Query().Get("target")
	if targetURL == "" {
		http.Error(responseWriter, "target parameter is missing", http.StatusBadRequest)
		return
	}

	scrapeTarget, err := target.Create(targetURL)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	registry := prometheus.NewRegistry()

	collectDurationGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "collect_duration_seconds",
		Help: "Returns how long the collect took to complete in seconds",
	})

	registry.MustRegister(collectDurationGauge)

	start := time.Now()

	if err := scraper.Scrape(scrapeTarget, registry); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	collectDurationGauge.Set(time.Since(start).Seconds())

	promhttp.
		HandlerFor(registry, promhttp.HandlerOpts{}).
		ServeHTTP(responseWriter, request)
}
