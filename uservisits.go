package main

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"sort"
	"strconv"
)

type UserVisits struct {
	Visits []UserVisit `json:"visits"`
}

type UserVisit struct {
	Mark       uint   `json:"mark"`
	Visited_at int    `json:"visited_at"`
	Place      string `json:"place"`
}

type UserVisitsFilter struct {
	fromDate   *int
	toDate     *int
	country    *string
	toDistance *uint
}

func userVisitsRequestHandler(ctx *fasthttp.RequestCtx, entityId uint, query *fasthttp.Args) {
	if user := usersMap.Get(entityId); user == nil {
		ctx.NotFound()
		return
	}

	var filters = UserVisitsFilter{}
	if query.Len() > 0 {
		if fromDate := query.Has("fromDate"); fromDate {
			if fromDateInt, err := strconv.Atoi(string(query.Peek("fromDate"))); err != nil {
				ctx.Error("{}", 400)
			} else {
				filters.fromDate = &fromDateInt
			}
		}
		if toDate := query.Has("toDate"); toDate {
			if toDateInt, err := strconv.Atoi(string(query.Peek("toDate"))); err != nil {
				ctx.Error("{}", 400)
			} else {
				filters.toDate = &toDateInt
			}
		}
		if country := query.Has("country"); country {
			countryName := string(query.Peek("country"))
			filters.country = &countryName
		}
		if toDistance := query.Has("toDistance"); toDistance {
			// get location id by Country
			if distanceInt, err := strconv.Atoi(string(query.Peek("toDistance"))); err != nil {
				ctx.Error("{}", 400)
			} else {
				distanceInt := uint(distanceInt)
				filters.toDistance = &distanceInt
			}
		}
	}

	response, _ := json.Marshal(UserVisits{getUserVisits(entityId, filters)})
	ctx.Success("application/json", response)
}

func getUserVisits(userId uint, filters UserVisitsFilter) []UserVisit {
	var userVisits = make(map[int]UserVisit, 0)
	for _, visit := range visitsByUserMap[userId] {
		if filters.fromDate != nil && visit.Visited_at < *filters.fromDate {
			continue
		}
		if filters.toDate != nil && visit.Visited_at > *filters.toDate {
			continue
		}
		location := locationsMap.Get(visit.Location)
		if filters.country != nil && location.Country != *filters.country {
			continue
		}
		if filters.toDistance != nil && location.Distance >= *filters.toDistance {
			continue
		}
		userVisits[visit.Visited_at] = UserVisit{visit.Mark, visit.Visited_at, location.Place}
	}

	visitedAtList := make([]int, 0, len(userVisits))
	for visitedAt := range userVisits {
		visitedAtList = append(visitedAtList, visitedAt)
	}
	sort.Ints(visitedAtList)

	sortedUserVisits := make([]UserVisit, 0, len(userVisits))
	for _, visitedAt := range visitedAtList {
		sortedUserVisits = append(sortedUserVisits, userVisits[visitedAt])
	}

	return sortedUserVisits
}
