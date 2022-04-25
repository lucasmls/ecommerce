package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	receivedMessagesCounterMetric := promauto.NewCounter(prometheus.CounterOpts{
		Name: "products_received_messages_total",
		Help: "The total number of received messages",
	})

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	for {
		receivedMessagesCounterMetric.Inc()
		time.Sleep(time.Second * 1)
	}
}
