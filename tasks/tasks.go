package tasks

import (
	"log"
	"time"
	"encoding/json"
	"encoding/base64"
	"net/http"
	"bytes"

	"github.com/RichardKnop/machinery/v1/signatures"
	vegeta "github.com/tsenart/vegeta/lib"
)

// AttackerOption is attack options.
type AttackerOption struct {}

// AttackTarget is a HTTP request blueprint.
type AttackTarget struct {
	Method string `validate:"required"`
	URL string `validate:"required"`
	Body string
	Headers map[string][]string
}

// GetVegetaTarget generates *vegeta.Target
func (t *AttackTarget) GetVegetaTarget() (*vegeta.Target) {
	headers := http.Header{}
	for k, v := range t.Headers {
		headers[k] = v
	}

	return &vegeta.Target{
		Method: t.Method,
		URL: t.URL,
		Body: []byte(t.Body),
		Header: headers,
	}
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
		vegetaTargets[i] = *tgt.GetVegetaTarget()
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

// AttackTaskSignature generates signatures.TaskSignature
func AttackTaskSignature(targets *[]AttackTarget, rate, duration uint64) (*signatures.TaskSignature, error) {
	jsonTargets, err := json.Marshal(&targets)
	if err != nil {
		return nil, err
	}

	taskSignature := &signatures.TaskSignature{
		Name: "attack",
		Args: []signatures.TaskArg{
			signatures.TaskArg{
				Type: "string",
				Value: string(jsonTargets),
			},
			signatures.TaskArg{
				Type: "uint64",
				Value: rate,
			},
			signatures.TaskArg{
				Type: "uint64",
				Value: duration,
			},
		},
	}
	return taskSignature, nil
}

func TaskAdd(args ...int64) (int64, error) {
	sum := int64(0)
	for _, arg := range args {
		sum += arg
	}
	return sum, nil
}
