package xk6reporter

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"go.k6.io/k6/js/modules"
)

//go:embed templates/report.tmpl
var templateString string

func init() {
	modules.Register("k6/x/reporter", new(Reporter))
}

type Reporter struct {
}

type ResultData struct {
	Title             string
	Metrics           map[string]interface{}
	ThresholdFailures int
	ThresholdTotal    int
	CheckFailures     int
	CheckPasses       int
	RootGroup         ReportRootGroup `json:"root_group"`
}

type ReportRootGroup struct {
	Groups []Group
	Checks []Check
}

type Group struct {
	Name   string
	Checks []Check
}

type Check struct {
	Name   string
	Passes int
	Fails  int
}

func (r *Reporter) GenerateReport(data []byte, reportTitle string) []byte {

	tmpl, err := template.New("").Funcs(sprig.FuncMap()).Parse(templateString)
	if err != nil {
		fmt.Println("Template file error", err)
		os.Exit(1)
	}

	resultData := ResultData{}
	json.Unmarshal(data, &resultData)

	resultData.Title = reportTitle

	thresholdFailures := 0
	thresholdTotal := 0
	for _, metric := range resultData.Metrics {
		metricMap := metric.(map[string]interface{})
		if metricMap["thresholds"] != nil {
			thresholds := metricMap["thresholds"].(map[string]interface{})
			thresholdTotal++
			for _, thres := range thresholds {
				thresMap := thres.(map[string]interface{})
				if thresMap["ok"] == false {
					thresholdFailures++
				}
			}
		}
	}

	resultData.ThresholdFailures = thresholdFailures
	resultData.ThresholdTotal = thresholdTotal

	checkFailures := 0
	checkPasses := 0

	for _, group := range resultData.RootGroup.Groups {
		for _, check := range group.Checks {
			checkFailures += check.Fails
			checkPasses += check.Passes
		}
	}

	for _, check := range resultData.RootGroup.Checks {
		checkFailures += check.Fails
		checkPasses += check.Passes
	}

	resultData.CheckFailures = checkFailures
	resultData.CheckPasses = checkPasses

	var b bytes.Buffer
	_ = tmpl.Execute(&b, resultData)
	return b.Bytes()
}
