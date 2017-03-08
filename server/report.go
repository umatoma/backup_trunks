package main

import (
	"bytes"
	"io"

	vegeta "github.com/tsenart/vegeta/lib"
)

// GetTextReporter make vegeta.Reporter from result binary
func GetTextReporter(bytesResult *[]byte) (vegeta.Reporter, error) {
	var reporter vegeta.Reporter
	var report vegeta.Report
	var m vegeta.Metrics
	reporter, report = vegeta.NewTextReporter(&m), &m

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

// GetPlotReporter make vegeta.Reporter from result binary
func GetPlotReporter(bytesResult *[]byte) (vegeta.Reporter, error) {
	var reporter vegeta.Reporter
	var report vegeta.Report
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
