// Copyright 2019 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package v1_test provides examples making requests to Prometheus using the
// Golang client.
package ovs_prom_client

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
	"github.com/kongseokhwan/Helios-prom-client/client"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

type TSMetricObj struct {
	Label      string
	Vals       []string
	TimeSeries []string
}

type OVSClient struct {
	Host 		string
	Port		int
	Version		string
}

func NewOVSPClilent(opts Options) (*OVSClient, error) {
	// TODO: opts validation check

	c := OVSClient{
		Host: opts.Host,
		Port: opts.Port,
		Version: opts.Version,
	}

	log.Debug("NewOVSPClilent() initialized successfully")
	return &c, nil
}

func parseMetric(res string) map[string][]string {
	var keyStr string

	metricMap := make(map[string][]string)
	repTimestamp := strings.NewReplacer(
		"[", "",
		"]", "",
		"@", "")

	res = strings.Replace(res, "=>", "", -1)
	testTmp1 := strings.Split(res, "\n")

	for _, line := range testTmp1 {
		if strings.Contains(line, "{") {
			keyStr = line
			metricMap[keyStr] = []string{""}
		} else {
			line = repTimestamp.Replace(line)
			metricMap[keyStr] = append(metricMap[keyStr], line)
		}
	}

	for key, val := range metricMap {
		fmt.Printf("key: %s, val: %s\n", key, val)
	}

	return metricMap
}

func ovsAPIQueryRange(host string, port int, query string) ([]TSMetricObj, error) {
	var queryResult []TSMetricObj

	client, err := api.NewClient(api.Config{
		Address: fmt.Sprint("http://%s:%s", host, port),
	})

	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return nil, err
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r := v1.Range{
		Start: time.Now().Add(-time.Hour),
		End:   time.Now(),
		Step:  time.Minute,
	}

	result, warnings, err := v1api.QueryRange(ctx, query, r)
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		return nil, err
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}

	resMetric := parseMetric(result.String())

	for key, val := range resMetric {
		var metricObj TSMetricObj
		metricObj.Label = key
		for i, v := range val {
			if i > 0 {
				metricList := strings.Fields(v)
				metricObj.Vals = append(metricObj.Vals, metricList[0])
				metricObj.TimeSeries = append(metricObj.TimeSeries, metricList[1])
			}
		}
		queryResult = append(queryResult, metricObj)
	}
	return queryResult, nil
}

func exampleAPIQuery() {
	client, err := api.NewClient(api.Config{
		Address: "http://13.209.193.98:9090",
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := v1api.Query(ctx, "up", time.Now())
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	fmt.Printf("Result:\n%v\n", result)
}

func exampleAPISeries() {
	client, err := api.NewClient(api.Config{
		Address: "http://13.209.193.98:9090",
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	lbls, warnings, err := v1api.Series(ctx, []string{
		"{__name__=\"ovs_flow_flow_bytes_total\", job=\"prometheus\"}",
	}, time.Now().Add(-time.Hour), time.Now())
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	fmt.Println("Result:")
	for _, lbl := range lbls {
		fmt.Println(lbl)
	}
}

func (c *OVSClient) ntopQueryWithRate(rankSize string, metric string, duration string) ([]TSMetricObj, error) {
	// Make Query String
	query := fmt.Sprintf(ovs_prom_ctx.ntopQueryWithRate, rankSize, metric, duration)

	// Call ovsAPIQueryRange() & return result
	return ovsAPIQueryRange(c.Host, c.Port, query)
}

func (c *OVSClient) countQuery(metric string) ([]TSMetricObj, error) {
	// Make Query String
	query := fmt.Sprintf(ovs_prom_ctx.countQuery, metric)

	// Call ovsAPIQueryRange() & return result
	return ovsAPIQueryRange(c.Host, c.Port, query)
}

func (c *OVSClient) avgbyQueryWithRate(metric string, duration string) ([]TSMetricObj, error) {
	// Make Query String
	query := fmt.Sprintf(ovs_prom_ctx.avgbyQueryWithRate, metric, duration)

	// Call ovsAPIQueryRange() & return result
	return ovsAPIQueryRange(c.Host, c.Port, query)
}