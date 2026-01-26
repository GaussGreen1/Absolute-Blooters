package server

import (
	"blooters/internal/handler"
	"blooters/internal/middleware"
	"bytes"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/klauspost/compress/snappy"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/prometheus/prompb"
)

type Server struct {
	mux http.Handler
}

func NewServer() *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/ping", handler.PingHandler)
	mux.HandleFunc("GET /api/games", handler.GamesHandler)
	mux.Handle("/metrics", promhttp.Handler())

	// Create remote write client for Grafana Cloud
	url := os.Getenv("GRAFANA_REMOTE_WRITE_URL")
	if url == "" {
		url = "https://prometheus-prod-65-prod-eu-west-2.grafana.net/api/prom/push"
	}
	username := os.Getenv("GRAFANA_USERNAME")
	password := os.Getenv("GRAFANA_PASSWORD")

	if username != "" && password != "" {
		go pushMetrics(url, username, password)
	}

	//Logging middleware:
	handler := middleware.LoggingMiddleware(mux)

	return &Server{
		mux: handler,
	}
}

func (s *Server) Start(addr string) error {
	log.Printf("Server starting on %s", addr)
	return http.ListenAndServe(addr, s.mux)
}

func pushMetrics(url, username, password string) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		mfs, err := prometheus.DefaultGatherer.Gather()
		if err != nil {
			log.Println("Error gathering metrics:", err)
			continue
		}
		ts := toTimeSeries(mfs)
		if len(ts) == 0 {
			continue
		}
		req := &prompb.WriteRequest{
			Timeseries: ts,
		}
		data, err := proto.Marshal(req)
		if err != nil {
			log.Println("Error marshaling:", err)
			continue
		}
		compressed := snappy.Encode(nil, data)
		httpClient := &http.Client{Timeout: 30 * time.Second}
		reqHttp, err := http.NewRequest("POST", url, bytes.NewReader(compressed))
		if err != nil {
			log.Println("Error creating request:", err)
			continue
		}
		reqHttp.SetBasicAuth(username, password)
		reqHttp.Header.Set("Content-Type", "application/x-protobuf")
		reqHttp.Header.Set("Content-Encoding", "snappy")
		reqHttp.Header.Set("X-Prometheus-Remote-Write-Version", "0.1.0")
		resp, err := httpClient.Do(reqHttp)
		if err != nil {
			log.Println("Error sending:", err)
			continue
		}
		resp.Body.Close()
		if resp.StatusCode != 200 {
			log.Println("Error status:", resp.StatusCode)
		}
	}
}

func toTimeSeries(mfs []*dto.MetricFamily) []prompb.TimeSeries {
	var ts []prompb.TimeSeries
	for _, mf := range mfs {
		for _, m := range mf.Metric {
			labels := []prompb.Label{{
				Name:  "__name__",
				Value: mf.GetName(),
			}}
			for _, l := range m.Label {
				labels = append(labels, prompb.Label{
					Name:  l.GetName(),
					Value: l.GetValue(),
				})
			}
			var timestamp int64
			if m.TimestampMs != nil {
				timestamp = *m.TimestampMs
			} else {
				timestamp = time.Now().UnixNano() / 1e6
			}
			switch mf.GetType() {
			case dto.MetricType_GAUGE:
				value := m.GetGauge().GetValue()
				samples := []prompb.Sample{{
					Value:     value,
					Timestamp: timestamp,
				}}
				ts = append(ts, prompb.TimeSeries{
					Labels:  labels,
					Samples: samples,
				})
			case dto.MetricType_COUNTER:
				value := m.GetCounter().GetValue()
				samples := []prompb.Sample{{
					Value:     value,
					Timestamp: timestamp,
				}}
				ts = append(ts, prompb.TimeSeries{
					Labels:  labels,
					Samples: samples,
				})
			case dto.MetricType_HISTOGRAM:
				h := m.GetHistogram()
				// Send count
				labelsCount := make([]prompb.Label, len(labels))
				copy(labelsCount, labels)
				labelsCount[0] = prompb.Label{Name: "__name__", Value: mf.GetName() + "_count"}
				ts = append(ts, prompb.TimeSeries{
					Labels:  labelsCount,
					Samples: []prompb.Sample{{Value: float64(h.GetSampleCount()), Timestamp: timestamp}},
				})
				// Send sum
				labelsSum := make([]prompb.Label, len(labels))
				copy(labelsSum, labels)
				labelsSum[0] = prompb.Label{Name: "__name__", Value: mf.GetName() + "_sum"}
				ts = append(ts, prompb.TimeSeries{
					Labels:  labelsSum,
					Samples: []prompb.Sample{{Value: h.GetSampleSum(), Timestamp: timestamp}},
				})
			}
		}
	}
	return ts
}
