package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus"
)

type Operation struct {
	ID        string   `json:"id,omitempty"`
	Name string   `json:"name,omitempty"`
	Scope  string   `json:"scope,omitempty"`
}

var operations []Operation

type metricsHandler struct {
	http.Handler
	duration *prometheus.HistogramVec
	requests *prometheus.CounterVec
}

func GetOperations(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(operations)
}

func AddOperation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var operation Operation
	_ = json.NewDecoder(r.Body).Decode(&operation)
	operation.ID=params["id"]
	operation.Name=params["name"]
	operation.Scope=params["scope"]
	operations = append(operations,operation)
	json.NewEncoder(w).Encode(operations)
}

func NewMetricsHandler(handler http.Handler) http.Handler {
	m := &metricsHandler{
		duration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "monitoring",
				Subsystem: "rest",
				Name:      "http_durations_histogram_seconds",
				Help:      "Request time duration.",
			},
			[]string{"code", "method"},
		),
		requests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "monitoring",
				Subsystem: "rest",
				Name:      "http_requests_total",
				Help:      "Total number of requests received.",
			},
			[]string{"code", "method"},
		),
	}

	prometheus.DefaultRegisterer.Register(m)

	return promhttp.InstrumentHandlerDuration(
		m.duration,
		promhttp.InstrumentHandlerCounter(
			m.requests,
			handler,
		),
	)
}

func main() {
	operations = append(operations, Operation{ID: "1", Name: "Upload", Scope: "PLMN-PLMN/MRBTS-1"})
	operations = append(operations, Operation{ID: "2", Name: "Upload", Scope: "PLMN-PLMN/MRBTS-2"})
	operations = append(operations, Operation{ID: "3", Name: "Upload", Scope: "PLMN-PLMN/MRBTS-3"})

	router := mux.NewRouter();
	router.HandleFunc("/operations", GetOperations).Methods("GET")
	router.HandleFunc("/operations/{id}", AddOperation).Methods("POST")

	router.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{},
	))

	log.Fatal(http.ListenAndServe("127.0.0.1:8002", NewMetricsHandler(router)))
}
