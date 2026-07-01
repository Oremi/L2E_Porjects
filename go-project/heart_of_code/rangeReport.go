package main

import "fmt"

type Report struct {
	EvenIndexSum, ScoreTotal, RuneCount int
}

func main() {
	fmt.Println(rangeReport([]int{2,3,4,5}, map[string]int{"a": 10, "b": 20}, "go"))
}


func rangeReport(nums []int, scores map[string]int, word string) Report {
	//report := new(Report)
	//return *report

	var report Report
	for idx, n := range nums {
		if idx == 0 || idx%2 == 0 {
			report.EvenIndexSum += n
		}
	}

	for _, score := range scores {report.ScoreTotal += score}
	for range word {report.RuneCount++}

	return report
}
