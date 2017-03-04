package main

import (
	"log"
	"time"
	"encoding/json"
	"encoding/base64"
	"net/http"
	"bytes"

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

		target := &vegeta.Target{
			Method: tgt.Method,
			URL: tgt.URL,
			Body: []byte(tgt.Body),
			Header: headers,
		}

		vegetaTargets[i] = *target
	}

	targeter := vegeta.NewStaticTargeter(vegetaTargets...)
	attacker := vegeta.NewAttacker()

	var resBuffer bytes.Buffer
	enc := vegeta.NewEncoder(&resBuffer)
	for res := range attacker.Attack(targeter, rate, duration) {
		log.Println("response:", res)
		enc.Encode(res)
	}

	return base64.StdEncoding.EncodeToString(resBuffer.Bytes()), nil
}

func TaskAdd(args ...int64) (int64, error) {
	sum := int64(0)
	for _, arg := range args {
		sum += arg
	}
	return sum, nil
}
