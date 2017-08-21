package main

import (
	"os"
	"testing"
	//"math/rand"
)

func TestWorkingDirectory(t *testing.T) {
	wd, _ := os.Getwd()
	t.Log(wd)

	parseFile("data/users_1.json")
	parseFile("data/locations_1.json")
	parseFile("data/visits_1.json")
}

//func BenchmarkGetLocationAvg(b *testing.B)  {
//	getLocationAvg(uint(rand.Intn(1000)), LocationAvgFilter{})
//}
//
//func BenchmarkAddNewVisit(b *testing.B) {
//	visit := Visit{9111,1,2, 7890, 5}
//	updateVisitsMaps(visit, nil)
//}
//
//func BenchmarkUpdateVisit(b *testing.B) {
//	visit := visitsMap[uint(rand.Intn(1000))]
//	updatedVisit := *visit
//
//	updatedVisit.Location = visit.Location + 1
//	updatedVisit.User = visit.User + 1
//
//	updateVisitsMaps(updatedVisit, visit)
//}

func BenchmarkGetUserVisits(b *testing.B) {
	for n := 0; n < b.N; n++ {
		//userId := uint(rand.Intn(1000))
		getUserVisits(1, UserVisitsFilter{})
	}

}
