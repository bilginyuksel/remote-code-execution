package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/yudai/pp"
)

const (
	Endpoint    = "http://3.67.10.139:8888/v2/codexec"
	ContentType = "application/json"

	TotalIterationCount       = 4
	CycleStartingRequestCount = 500
	MaxIncreaseCoefficient    = 4
)

type RequestBody struct {
	Lang    string   `json:"lang"`
	Content string   `json:"content"`
	Args    []string `json:"args"`
}

var requestBodies = []RequestBody{
	{
		Lang:    "python3",
		Content: "import sys\ninp = sys.argv[1]\nfor i in range(int(inp)):\n\tprint('*' * (i +1))",
		Args:    []string{"20"},
	},
	{
		Lang:    "",
		Content: ContentType,
		Args:    []string{},
	},
	{
		Lang:    "",
		Content: ContentType,
		Args:    []string{},
	},
	{
		Lang:    "",
		Content: ContentType,
		Args:    []string{},
	},
}

func main() {
	var iterationReports [][]Report

	for i := 0; i < TotalIterationCount; i++ {
		var reports []Report

		randomIncreaseCoefficient := rand.Intn(MaxIncreaseCoefficient) + 1
		totalRequestCountForCycle := CycleStartingRequestCount * randomIncreaseCoefficient

		gw := sync.WaitGroup{}
		gw.Add(totalRequestCountForCycle)

		log.Printf("Iteration %d starts with %d concurrent requests\n", i, totalRequestCountForCycle)
		for j := 0; j < totalRequestCountForCycle; j++ {
			go func(idx int) {
				singleReport := hit(idx % len(requestBodies))
				reports = append(reports, singleReport)
			}(j)
		}

		gw.Wait()
		iterationReports = append(iterationReports, reports)

		log.Printf("Iteration %d finished..\nWaiting for 3 seconds..\n", i)
		time.Sleep(3 * time.Second)
	}

	summaries := calculateReportsSummaries(iterationReports)
	pp.Println(summaries)
}

type ReportsSummary struct {
	MaxExecutionTimeMs int64
	AvgExecutionTimeMs int64
	MinExecutionTimeMs int64

	MaxServerExecutionTimeMs int64
	AvgServerExecutionTimeMs int64
	MinServerExecutionTimeMs int64

	TotalCount   int64
	GarbageCount int64
	FailCount    int64
	SuccessCount int64

	SuccessRate float64
	FailRate    float64
	GarbageRate float64
}

func calculateReportsSummaries(iterationReports [][]Report) []ReportsSummary {
	var summaries []ReportsSummary

	for _, reports := range iterationReports {
		summary := ReportsSummary{}
		for idx := range reports {
			report := reports[idx]
			if report.Failed {
				summary.FailCount++
			} else if report.Garbage {
				summary.GarbageCount++
			} else {
				summary.SuccessCount++
			}

			summary.MaxExecutionTimeMs = max(summary.MaxExecutionTimeMs, report.ExecutionTime.Milliseconds())
			summary.MinExecutionTimeMs = max(summary.MinExecutionTimeMs, report.ExecutionTime.Milliseconds())
			summary.AvgExecutionTimeMs += report.ExecutionTime.Milliseconds()

			summary.MaxServerExecutionTimeMs = max(summary.MaxServerExecutionTimeMs, report.ServerExecutionTime.Milliseconds())
			summary.MinServerExecutionTimeMs = max(summary.MinServerExecutionTimeMs, report.ServerExecutionTime.Milliseconds())
			summary.AvgServerExecutionTimeMs += report.ServerExecutionTime.Milliseconds()
		}

		summary.TotalCount = summary.FailCount + summary.SuccessCount + summary.GarbageCount
		summary.SuccessRate = float64(summary.SuccessCount) / float64(summary.TotalCount) * 100
		summary.FailRate = float64(summary.FailCount) / float64(summary.TotalCount) * 100
		summary.GarbageRate = float64(summary.GarbageCount) / float64(summary.TotalCount) * 100
		summary.AvgExecutionTimeMs = summary.AvgExecutionTimeMs / summary.TotalCount
		summary.AvgServerExecutionTimeMs = summary.AvgServerExecutionTimeMs / summary.TotalCount
	}

	return summaries
}

func min(t1, t2 int64) int64 {
	if t1 < t2 {
		return t1
	}
	return t2
}

func max(t1, t2 int64) int64 {
	if t1 > t2 {
		return t1
	}
	return t2
}

type ServerResponse struct {
	Output          string `json:"output"`
	Message         string `json:"message"`
	ExecutionTimeMs int64  `json:"execution_time_ms"`
}

type Report struct {
	Garbage             bool
	Cause               string
	Failed              bool
	ExecutionTime       time.Duration
	ServerExecutionTime time.Duration
}

func hit(idx int) Report {
	var report Report
	randomRequestBody := requestBodies[idx]
	reqBytes, _ := json.Marshal(randomRequestBody)

	startTime := time.Now()
	res, err := http.DefaultClient.Post(Endpoint, ContentType, bytes.NewBuffer(reqBytes))
	report.ExecutionTime = time.Since(startTime)
	if err != nil {
		report.Garbage = true
		return report
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		report.Failed = true
		report.Cause = err.Error()
		return report
	}

	var response ServerResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		report.Cause = err.Error()
		report.Failed = true
		return report
	}
	report.ServerExecutionTime, _ = time.ParseDuration(fmt.Sprintf("%dms", response.ExecutionTimeMs))

	isServerResponseOK := res.StatusCode >= 200 && res.StatusCode < 400
	if !isServerResponseOK {
		report.Cause = response.Message
		report.Failed = true
	}

	return report
}
