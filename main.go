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
package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

type StatsObj struct {
	Label string
	Vals  []string
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

func parseInterfaceRxBytes(res string) []StatsObj {
	var statsObjList []StatsObj

	parseResult := strings.Split(res, "{")[1:]
	for _, vector := range parseResult {
		var tmpStatsObj StatsObj
		tmpResult := strings.Split(vector, "=>")

		tmpLabel := tmpResult[0]
		tmpVal := strings.Split(tmpResult[1], "\n")

		tmpStatsObj.Label = tmpLabel
		tmpStatsObj.Vals = tmpVal

		statsObjList = append(statsObjList, tmpStatsObj)
	}
	return statsObjList
}

func exampleAPIQueryRange() {
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
	r := v1.Range{
		Start: time.Now().Add(-time.Hour),
		End:   time.Now(),
		Step:  time.Minute,
	}
	result, warnings, err := v1api.QueryRange(ctx, "rate(ovs_interface_receive_bytes_total[5m])", r)
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	resSlice := parseInterfaceRxBytes(result.String())

	for _, statsObj := range resSlice {
		fmt.Printf("Labels: %v\n", statsObj.Label)
		fmt.Printf("Vals: %v\n", statsObj.Vals)
	}
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

func main() {
	exampleAPIQueryRange()
}
