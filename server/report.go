package main

import (
	"bytes"
	"io"

	vegeta "github.com/tsenart/vegeta/lib"
)

// GetPlotReporter generates vegeta.Reporter from result binary
func GetPlotReporter(bytesResult *[]byte) (vegeta.Reporter, error) {
	var reporter vegeta.Reporter
	var report vegeta.Report
	// text reporter
	// var metrics vegeta.Metrics
	// reporter, report = vegeta.NewTextReporter(&metrics), &metrics

	// plot reporter
	var rs vegeta.Results
	reporter, report = vegeta.NewPlotReporter("Vegeta Plot", &rs), &rs

	decoder := vegeta.NewDecoder(bytes.NewReader(*bytesResult))
	for {
		var r vegeta.Result
		if err := decoder.Decode(&r); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		report.Add(&r)
	}

	if c, ok := report.(vegeta.Closer); ok {
		c.Close()
	}

	return reporter, nil
}
