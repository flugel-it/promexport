package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"strings"
	"time"
)

func main() {
	fmt.Println("promexporter")

	showprefix := flag.Bool("showprefix", false,
		"Display prefix and time series label")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := api.NewClient(api.Config{Address: "http://localhost:9090"})
	if err != nil {
		fmt.Errorf("Ooops connecting to the API", err)
	}

	query_api := v1.NewAPI(client)
	// Getting up metric first, to get all the instances being monitored.
	// Getting all the metrics for all the instances is timming out.
	match := []string{"up"}
	date_from := time.Now().AddDate(0, -1, 0)
	date_to := time.Now()

	series, err := query_api.Series(ctx, match, date_from, date_to)
	if err != nil {
		fmt.Errorf("Ooops getting series", err)
	}

	for _, serie := range series {
		instance_match := []string{fmt.Sprintf("{instance=\"%s\"}",
			serie["instance"])}
		instance_series, err := query_api.Series(ctx, instance_match,
			date_from, date_to)
		if err != nil {
			fmt.Errorf("Ooops getting series", err)
		}
		for _, instance_serie := range instance_series {
			if *showprefix {
				name_fields := strings.Split(fmt.Sprintf("%s",
					instance_serie["__name__"]), "_")
				fmt.Println(name_fields[0], instance_serie["__name__"])
			} else {
				fmt.Println(instance_serie)
			}
		}
	}
}
