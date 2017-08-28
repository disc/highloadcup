package main

import (
	"bytes"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"strconv"
)

type LocationAvg struct {
	Avg float64 `json:"avg"`
}

type LocationAvgFilter struct {
	fromDate *int
	toDate   *int
	fromAge  *int
	toAge    *int
	gender   *[]byte
}

func locationAvgRequestHandler(ctx *fasthttp.RequestCtx, locationId uint, query *fasthttp.Args) {
	if location := locationsMap.Get(locationId); location == nil {
		ctx.NotFound()
		return
	}

	var filters = LocationAvgFilter{}
	if fromDate := query.Has("fromDate"); fromDate {
		if fromDateInt, err := strconv.Atoi(string(query.Peek("fromDate"))); err != nil {
			ctx.Error("{}", 400)
			return
		} else {
			filters.fromDate = &fromDateInt
		}
	}
	if toDate := query.Has("toDate"); toDate {
		if toDateInt, err := strconv.Atoi(string(query.Peek("toDate"))); err != nil {
			ctx.Error("{}", 400)
			return
		} else {
			filters.toDate = &toDateInt
		}
	}
	if fromAge := query.Has("fromAge"); fromAge {
		if fromAgeInt, err := strconv.Atoi(string(query.Peek("fromAge"))); err != nil {
			ctx.Error("{}", 400)
			return
		} else {
			filters.fromAge = &fromAgeInt
		}
	}
	if toAge := query.Has("toAge"); toAge {
		if toAgeInt, err := strconv.Atoi(string(query.Peek("toAge"))); err != nil {
			ctx.Error("{}", 400)
			return
		} else {
			filters.toAge = &toAgeInt
		}
	}
	if gender := query.Has("gender"); gender {
		if genderStr := query.Peek("gender"); len(genderStr) > 0 && bytes.ContainsAny(genderStr, "mf") {
			filters.gender = &genderStr
		} else {
			ctx.Error("{}", 400)
			return
		}
	}

	response, _ := json.Marshal(LocationAvg{Round(getLocationAvg(locationId, filters), .5, 5)})
	ctx.Success("application/json", response)
}

func getLocationAvg(locationId uint, filters LocationAvgFilter) float64 {
	marks := make([]uint, 0)
	var marksSum uint
	for _, visit := range visitsByLocationMap[locationId] {
		if filters.fromDate != nil && visit.Visited_at < *filters.fromDate {
			continue
		}
		if filters.toDate != nil && visit.Visited_at > *filters.toDate {
			continue
		}
		user := usersMap.Get(visit.User)
		if filters.fromAge != nil || filters.toAge != nil {
			if filters.fromAge != nil && user.Birth_date > getTimestampByAge(filters.fromAge, now) {
				continue
			}
			if filters.toAge != nil && user.Birth_date <= getTimestampByAge(filters.toAge, now) {
				continue
			}
		}
		if filters.gender != nil && string(*filters.gender) != user.Gender {
			continue
		}
		marksSum += visit.Mark
		marks = append(marks, visit.Mark)
	}

	if len(marks) > 0 {
		return float64(marksSum) / float64(len(marks))
	}

	return 0
}
