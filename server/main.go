package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const HOST string{"15.165.203.82"}
const PORT string{"9090"}
const VERSION string{"v1"}

const PARAMMETRIC string{"metricID"}
const PARAMRANK string{"rankID"}
const PARAMDURATION string{"durationID"}

func getCountAPIQuery(w http.ResponseWriter, r *http.Request) {
	c, err := ovs_prom_client.NewOVSPClilent(HOST, PORT, VERSION)
	pathParams := mux.Vars(r)
	metricID := ""

	if val, ok := pathParams[PARAMMETRIC]; ok {
		metricID, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "need a number"}`))
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	/*
		1. Make Query String according to the metricID
		2. Call OVSClient API : countQuery(metric string) ([]TSMetricObj, error)
		3. Marsha JSON
	*/
	w.Write([]byte(`{"message": "get called"}`))
}

func getTopkAPIQuery(w http.ResponseWriter, r *http.Request) {
	c, err := ovs_prom_client.NewOVSPClilent(HOST, PORT, VERSION)
	pathParams := mux.Vars(r)
	metricID := ""
	durationID := ""
	rankID := ""

	if val, ok := pathParams[PARAMMETRIC]; ok {
		metricID, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "need a metric"}`))
			return
		}
	}

	if val, ok := pathParams[PARAMDURATION]; ok {
		durationID, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "need a duration"}`))
			return
		}
	}

	if val, ok := pathParams[PARAMRANK]; ok {
		rankID, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "need a rank"}`))
			return
		}
	}

	/*
		1. Make Query String
		2. Call OVSClient API : ntopQueryWithRate(rankSize string, metric string, duration string) ([]TSMetricObj, error)
		3. Marsha JSON
	*/
	w.Write([]byte(`{"message": "get called"}`))
}

func getGroupbyAPIQueryRange(w http.ResponseWriter, r *http.Request) {
	c, err := ovs_prom_client.NewOVSPClilent(HOST, PORT, VERSION)
	pathParams := mux.Vars(r)
	metricID := ""
	durationID := ""

	if val, ok := pathParams[PARAMMETRIC]; ok {
		metricID, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "need a metric"}`))
			return
		}
	}

	if val, ok := pathParams[PARAMDURATION]; ok {
		durationID, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "need a duration"}`))
			return
		}
	}

	/*
		1. Make Query String
		2. Call OVSClient API : avgbyQueryWithRate(metric string, duration string) ([]TSMetricObj, error)
		3. Marsha JSON
	*/
	w.Write([]byte(`{"message": "get called"}`))
}

func post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "post called"}`))
}

func put(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"message": "put called"}`))
}

func delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "delete called"}`))
}

func params(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	userID := -1
	var err error
	if val, ok := pathParams["userID"]; ok {
		userID, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "need a number"}`))
			return
		}
	}

	commentID := -1
	if val, ok := pathParams["commentID"]; ok {
		commentID, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "need a number"}`))
			return
		}
	}

	query := r.URL.Query()
	location := query.Get("location")

	w.Write([]byte(fmt.Sprintf(`{"userID": %d, "commentID": %d, "location": "%s" }`, userID, commentID, location)))
}

func main() {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/count/metric/{metricID}", getCountAPIQuery).Methods(http.MethodGet)
	api.HandleFunc("/topk/metric/{metricID}/duration/{durationID}/rank/{rankID}", getTopkAPIQuery).Methods(http.MethodGet)
	api.HandleFunc("/groupby/metric/{metricID}/duration/{durationID}", getGroupbyAPIQueryRange).Methods(http.MethodGet)

	// Sample
	api.HandleFunc("", post).Methods(http.MethodPost)
	api.HandleFunc("", put).Methods(http.MethodPut)
	api.HandleFunc("", delete).Methods(http.MethodDelete)
	api.HandleFunc("/user/{userID}/comment/{commentID}", params).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8081", r))
}
