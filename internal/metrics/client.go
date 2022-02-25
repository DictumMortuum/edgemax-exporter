package metrics

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type Client struct {
	httpClient http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

func (c *Client) Metrics() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		stats, err := c.getStatistics()
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		c.setMetrics(stats)

		promhttp.Handler().ServeHTTP(writer, request)
	}
}

func (c *Client) setMetrics(data map[string]Stats) {
	for iface, stats := range data {
		var isEnabled int = 0
		if stats.Status == true {
			isEnabled = 1
		}

		Status.WithLabelValues(iface).Set(float64(isEnabled))
		Pings.WithLabelValues(iface).Set(float64(stats.Pings))
		Fails.WithLabelValues(iface).Set(float64(stats.Fails))
		Drops.WithLabelValues(iface).Set(float64(stats.Drops))
	}
}

func (c *Client) getStatistics() (map[string]Stats, error) {
	stats, err := GetStatistics()
	if err != nil {
		return nil, err
	}

	rs := map[string]Stats{}

	for _, stat := range stats {
		rs[stat.Interface] = stat
	}

	return rs, nil
}
