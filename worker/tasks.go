package main

import (
	"log"
	"fmt"
	"encoding/json"
	"net/http"

	vegeta "github.com/tsenart/vegeta/lib"
)

// AttackerOption is attack options.
type AttackerOption struct {}

// AttackTarget is a HTTP request blueprint.
type AttackTarget struct {
	Method string `json:"method"`
	URL string `json:"url"`
	Body string `json:"body"`
	Headers map[string][]string `json:"headers"`
}

// TaskAttack performs load test and returns the measurement results.
func TaskAttack(tgts string, rate, duration uint64) (string, error) {
	attackTargets := make([]AttackTarget, 0)
	if err := json.Unmarshal([]byte(tgts), &attackTargets); err != nil {
		log.Println(err)
	}
	fmt.Println(attackTargets)

	vegetaTargets := make([]vegeta.Target, len(attackTargets))
	for _, tgt := range attackTargets {
		headers := http.Header{}
		for k, v := range tgt.Headers {
			headers[k] = v
		}

		target := &vegeta.Target{
			Method: tgt.Method,
			URL: tgt.URL,
			Body: []byte(tgt.Body),
			Header: headers,
		}

		vegetaTargets = append(vegetaTargets, *target)
	}

	return "", nil
}

func TaskAdd(args ...int64) (int64, error) {
	sum := int64(0)
	for _, arg := range args {
		sum += arg
	}
	return sum, nil
}
