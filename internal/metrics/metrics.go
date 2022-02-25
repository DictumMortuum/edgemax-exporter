package metrics

import (
	"bytes"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Stats struct {
	Interface string `json:"interface"`
	Status    bool   `json:"status"`
	Pings     int    `json:"pings"`
	Fails     int    `json:"fails"`
	Drops     int    `json:"drops"`
}

var (
	Status = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "status",
			Namespace: "router",
		},
		[]string{"interface"},
	)

	Pings = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "pings",
			Namespace: "router",
		},
		[]string{"interface"},
	)

	Fails = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "fails",
			Namespace: "router",
		},
		[]string{"interface"},
	)

	Drops = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "drops",
			Namespace: "router",
		},
		[]string{"interface"},
	)

	re_pppoe = regexp.MustCompile(`(pppoe.)\s+status:\s+(\S+)\s+pings:\s+(\d+)\s+fails:\s+(\d+)\s+run fails: \S+\s+route drops:\s+(\d+)`)
)

func Init() {
	initMetric("router_status", Status)
	initMetric("router_pings", Pings)
	initMetric("router_fails", Fails)
	initMetric("router_drops", Drops)
}

func initMetric(name string, metric *prometheus.GaugeVec) {
	prometheus.MustRegister(metric)
	log.Printf("New Prometheus metric registered: %s", name)
}

func hal(cmd string, args ...string) (string, error) {
	var out bytes.Buffer
	c := exec.Command(cmd, args...)
	c.Stdout = &out

	err := c.Run()
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

func atoi(s string) int {
	i, _ := strconv.Atoi(strings.TrimSpace(s))
	return i
}

func GetStatistics() ([]Stats, error) {
	stats := []Stats{}

	raw, err := hal("/usr/sbin/ubnt-hal", "wlbGetWdStatus")
	if err != nil {
		return nil, err
	}

	refs := re_pppoe.FindAllStringSubmatch(raw, -1)
	for _, ref := range refs {
		if len(ref) == 6 {
			stats = append(stats, Stats{
				Interface: ref[1],
				Status:    ref[2] == "OK",
				Pings:     atoi(ref[3]),
				Fails:     atoi(ref[4]),
				Drops:     atoi(ref[5]),
			})
		}
	}

	return stats, nil
}
