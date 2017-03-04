package main

import (
	"fmt"
	"time"
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
func TaskAttack(tgts string, rate, du uint64) (string, error) {
	duration := time.Duration(du) * time.Second
	attackTargets := make([]AttackTarget, 0)
	if err := json.Unmarshal([]byte(tgts), &attackTargets); err != nil {
		return "", err
	}

	vegetaTargets := make([]vegeta.Target, len(attackTargets))
	for i, tgt := range attackTargets {
		headers := http.Header{}
		for k, v := range tgt.Headers {
			headers[k] = v
		}

		fmt.Println("target:", tgt.Method, tgt.URL)
		target := &vegeta.Target{
			Method: tgt.Method,
			URL: tgt.URL,
			Body: []byte(tgt.Body),
			Header: headers,
		}

		vegetaTargets[i] = *target
	}

	fmt.Println(vegetaTargets)
	targeter := vegeta.NewStaticTargeter(vegetaTargets...)
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration) {
		fmt.Println("response:", res)
		metrics.Add(res)
	}
	metrics.Close()

	jsonMetrics, err := json.Marshal(metrics)
	if err != nil {
		return "", err
	}

	return string(jsonMetrics), nil
}

func TaskAdd(args ...int64) (int64, error) {
	sum := int64(0)
	for _, arg := range args {
		sum += arg
	}
	return sum, nil
}
