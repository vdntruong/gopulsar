package app

import (
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strconv"
)

var clusterAddr = "pulsar://localhost:6650"

func Run() {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: clusterAddr,
	})
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	prometheusPort := 2112
	log.Printf("Starting Prometheus metrics at http://localhost:%v/metrics\n", prometheusPort)

	http.Handle("/metrics", promhttp.Handler())
	err = http.ListenAndServe(":"+strconv.Itoa(prometheusPort), nil)
	if err != nil {
		log.Fatal(err)
	}
}
