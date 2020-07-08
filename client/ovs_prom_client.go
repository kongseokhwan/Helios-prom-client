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
	"strings"
	"time"

	"github.com/prometheus/client_golang/api"
	ovs_prom_ctx "github.com/kongseokhwan/Helios-prom-client/client"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/log"
)

type TSMetricObj struct {
	Label      string
	Vals       []string
	TimeSeries []string
}

type OVSClient struct {
	Host    string
	Port    int
	Version string
}

func NewOVSPClilent(host string, port string, version string) (*OVSClient, error) {
	// TODO: opts validation check

	c := OVSClient{
		Host:    opts.Host,
		Port:    opts.Port,
		Version: opts.Version,
	}

	log.Debug("NewOVSPClilent() initialized successfully")
	return &c, nil
}

func parseCountMetric(res string) map[string][]string {
	metricMap := make(map[string][]string)
	res = strings.Split(res, "=>")[1]
	repTimestamp := strings.NewReplacer(
		"[", "",
		"]", "",
		"@", "")

	res = repTimestamp.Replace(res)

	metricMap["count"] = []string{""}
	metricMap["count"] = append(metricMap["count"], res)

	return metricMap
}

func parseTopkMetric(res string) map[string][]string {
	var keyStr string

	metricMap := make(map[string][]string)
	repTimestamp := strings.NewReplacer(
		"[", "",
		"]", "",
		"@", "",
		"{", "",
		"}", "",
	)

	testTmp1 := strings.Split(res, "\n")

	for _, line := range testTmp1 {
		metricStr := strings.Split(line, "=>")

		keyStr = repTimestamp.Replace(metricStr[0])
		valStr := repTimestamp.Replace(metricStr[1])

		metricMap[keyStr] = []string{""}
		metricMap[keyStr] = append(metricMap[keyStr], valStr)
	}

	for key, val := range metricMap {
		fmt.Printf("key: %s, val: %s\n", key, val)
	}

	return metricMap
}

func parseGroupByMetric(res string) map[string][]string {
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

func countAPIQuery(host string, port int, query string) ([]TSMetricObj, error) {
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

	result, warnings, err := v1api.Query(ctx, query, time.Now())
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		return nil, err
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}

	fmt.Printf("Result Strnings: %v\n", result)

	resMetric := parseCountMetric(result.String())

	for key, val := range resMetric {
		var metricObj TSMetricObj
		metricObj.Label = key
		for i, v := range val {
			if i > 0 {
				metricList := strings.Fields(v)
				metricObj.Vals = append(metricObj.Vals, metricList[0])
				metricObj.TimeSeries = append(metricObj.TimeSeries, metricList[1])
				fmt.Printf("Val: %s, Time : %s \n", metricList[0], metricList[1])
			}
		}
		queryResult = append(queryResult, metricObj)
	}
	return queryResult, err
}

func topkAPIQuery(host string, port int, query string) ([]TSMetricObj, error) {
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

	result, warnings, err := v1api.Query(ctx, query, time.Now())
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		return nil, err
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}

	fmt.Printf("Result Strnings: %v\n", result)

	resMetric := parseTopkMetric(result.String())

	for key, val := range resMetric {
		var metricObj TSMetricObj
		metricObj.Label = key
		for i, v := range val {
			if i > 0 {
				metricList := strings.Fields(v)
				metricObj.Vals = append(metricObj.Vals, metricList[0])
				metricObj.TimeSeries = append(metricObj.TimeSeries, metricList[1])
				fmt.Printf("Val: %s, Time : %s \n", metricList[0], metricList[1])
			}
		}
		queryResult = append(queryResult, metricObj)
	}

	return queryResult, nil
}

func groupbyAPIQueryRange(host string, port int, query string) ([]TSMetricObj, error) {
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

	resMetric := parseGroupByMetric(result.String())

	for key, val := range resMetric {
		var metricObj TSMetricObj
		metricObj.Label = key
		for i, v := range val {
			if i > 0 {
				metricList := strings.Fields(v)
				metricObj.Vals = append(metricObj.Vals, metricList[0])
				metricObj.TimeSeries = append(metricObj.TimeSeries, metricList[1])
				fmt.Printf("Val: %s, Time : %s \n", metricList[0], metricList[1])
			}
		}
		queryResult = append(queryResult, metricObj)
	}
	return queryResult, nil
}

func (c *OVSClient) ntopQueryWithRate(rankSize string, metric string, duration string) ([]TSMetricObj, error) {
	// Make Query String
	ovs_prom_ctx.
	query := fmt.Sprintf(ovs_prom_ctx.ntopQueryWithRate, rankSize, metric, duration)

	// Call ovsAPIQueryRange() & return result
	return topkAPIQuery(c.Host, c.Port, query)
}

func (c *OVSClient) countQuery(metric string) ([]TSMetricObj, error) {
	// Make Query String
	query := fmt.Sprintf(ovs_prom_ctx.countQuery, metric)

	// Call ovsAPIQueryRange() & return result
	return countAPIQuery(c.Host, c.Port, query)
}

func (c *OVSClient) avgbyQueryWithRate(metric string, duration string) ([]TSMetricObj, error) {
	// Make Query String
	query := fmt.Sprintf(ovs_prom_ctx.avgbyQueryWithRate, metric, duration)

	// Call ovsAPIQueryRange() & return result
	return groupbyAPIQueryRange(c.Host, c.Port, query)
}
