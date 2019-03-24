package scraper

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Target interface {
	Scrape() ([]byte, error)
}

type Pool struct {
	Name                string    `json:"pool"`
	ProcessManager      string    `json:"process manager"`
	StartTime           Timestamp `json:"start time"`
	StartSince          int64     `json:"start since"`
	AcceptedConnections int64     `json:"accepted conn"`
	ListenQueue         int64     `json:"listen queue"`
	MaxListenQueue      int64     `json:"max listen queue"`
	ListenQueueLength   int64     `json:"listen queue len"`
	IdleProcesses       int64     `json:"idle processes"`
	ActiveProcesses     int64     `json:"active processes"`
	TotalProcesses      int64     `json:"total processes"`
	MaxActiveProcesses  int64     `json:"max active processes"`
	MaxChildrenReached  int64     `json:"max children reached"`
	SlowRequests        int64     `json:"slow requests"`
}

type Timestamp time.Time

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	ts := time.Time(*t).Unix()
	stamp := fmt.Sprint(ts)
	return []byte(stamp), nil
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	ts, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}
	*t = Timestamp(time.Unix(int64(ts), 0))
	return nil
}

func Scrape(target Target, registry *prometheus.Registry) error {
	var (
		startSinceGauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "start_since",
			Help: "The number of seconds since FPM has started.",
		})

		acceptedConnectionsGauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "accepted_connections",
			Help: "The number of requests accepted by the pool.",
		})

		listenQueueGauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "listen_queue",
			Help: "The number of requests in the queue of pending connections.",
		})

		maxListenQueueGauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "max_listen_queue",
			Help: "The maximum number of requests in the queue of pending connections since FPM has started.",
		})

		listenQueueLengthGauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "listen_queue_length",
			Help: "The size of the socket queue of pending connections.",
		})

		idleProcessesGauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "idle_processes",
			Help: "The number of idle processes.",
		})

		activeProcessesGauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "active_processes",
			Help: "The number of active processes.",
		})

		totalProcessesGauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "total_processes",
			Help: "The number of idle + active processes.",
		})

		maxActiveProcessesGauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "max_active_processes",
			Help: "The maximum number of active processes since FPM has started.",
		})
		maxChildrenReachedGauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "max_children_reached",
			Help: "The number of times, the process limit has been reached, when pm tries to start more children (works only for pm 'dynamic' and 'ondemand').",
		})

		slowRequestsGauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "slow_requests",
			Help: "The number of requests that exceeded your 'request_slowlog_timeout' value.",
		})
	)

	registry.MustRegister(
		startSinceGauge,
		acceptedConnectionsGauge,
		listenQueueGauge,
		maxListenQueueGauge,
		listenQueueLengthGauge,
		idleProcessesGauge,
		activeProcessesGauge,
		totalProcessesGauge,
		maxActiveProcessesGauge,
		maxChildrenReachedGauge,
		slowRequestsGauge,
	)

	data, err := target.Scrape()
	if err != nil {
		return err
	}

	pool := Pool{}

	if err := json.Unmarshal(data, &pool); err != nil {
		return errors.Wrap(err, "response unmarshal error")
	}

	startSinceGauge.Set(float64(pool.StartSince))
	acceptedConnectionsGauge.Set(float64(pool.AcceptedConnections))
	listenQueueGauge.Set(float64(pool.ListenQueue))
	maxListenQueueGauge.Set(float64(pool.MaxListenQueue))
	listenQueueLengthGauge.Set(float64(pool.ListenQueueLength))
	idleProcessesGauge.Set(float64(pool.IdleProcesses))
	activeProcessesGauge.Set(float64(pool.ActiveProcesses))
	totalProcessesGauge.Set(float64(pool.TotalProcesses))
	maxActiveProcessesGauge.Set(float64(pool.MaxActiveProcesses))
	maxChildrenReachedGauge.Set(float64(pool.MaxChildrenReached))
	slowRequestsGauge.Set(float64(pool.SlowRequests))

	return nil
}
