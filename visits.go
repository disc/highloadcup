package main

import (
	"encoding/json"
	"errors"
	"github.com/valyala/fasthttp"
)

type Visit struct {
	Id         uint `json:"id"`
	Location   uint `json:"location"`
	User       uint `json:"user"`
	Visited_at int  `json:"visited_at"`
	Mark       uint `json:"mark"`
}

func getVisitRequestHandler(ctx *fasthttp.RequestCtx, entityId uint) {
	if visit, ok := visitsMap[entityId]; ok {
		response, _ := json.Marshal(visit)
		ctx.Success("application/json", response)
		return
	}
	ctx.NotFound()
}

func createVisitRequestHandler(ctx *fasthttp.RequestCtx) {
	if visit, err := createVisit(ctx.PostBody()); err == nil {
		ctx.SetConnectionClose()
		ctx.Success("application/json", []byte("{}"))

		updateVisitsMaps(*visit, nil)
		return
	}
	ctx.Error("{}", 400)
}

func updateVisitRequestHandler(ctx *fasthttp.RequestCtx, entityId uint) {
	if visit, ok := visitsMap[entityId]; ok {
		if updatedVisit, err := updateVisit(ctx.PostBody(), *visit); err == nil {
			ctx.SetConnectionClose()
			ctx.Success("application/json", []byte("{}"))
			updateVisitsMaps(*updatedVisit, visit)
			return
		}
		ctx.Error("{}", 400)
		return
	}
	ctx.NotFound()
}

func createVisit(postData []byte) (*Visit, error) {
	visit := Visit{}
	if err := json.Unmarshal(postData, &visit); err != nil {
		return nil, err
	}

	if visit.Id == 0 || visit.Location == 0 || visit.User == 0 || visit.Mark > 5 {
		return nil, errors.New("Validation error")
	}
	if _, ok := visitsMap[visit.Id]; ok {
		return nil, errors.New("Visit already exists")
	}

	return &visit, nil
}

func updateVisit(postData []byte, visit Visit) (*Visit, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(postData, &data); err != nil {
		return nil, err
	}

	updatedVisit := visit

	if location, ok := data["location"]; ok {
		if location != nil {
			updatedVisit.Location = uint(location.(float64))
		} else {
			return nil, errors.New("Field validation error")
		}
	}
	if user, ok := data["user"]; ok {
		if user != nil {
			updatedVisit.User = uint(user.(float64))
		} else {
			return nil, errors.New("Field validation error")
		}
	}
	if visitedAt, ok := data["visited_at"]; ok {
		if visitedAt != nil {
			updatedVisit.Visited_at = int(visitedAt.(float64))
		} else {
			return nil, errors.New("Field validation error")
		}
	}
	if mark, ok := data["mark"]; ok {
		if mark != nil && uint(mark.(float64)) <= 5 {
			updatedVisit.Mark = uint(mark.(float64))
		} else {
			return nil, errors.New("Field validation error")
		}
	}

	return &updatedVisit, nil
}

func updateVisitsMaps(visit Visit, prevRef *Visit) {
	var (
		prevVisit       *Visit
		ok              bool
		userChanged     = false
		locationChanged = false
	)
	if prevRef == nil {
		visitsMap[visit.Id] = &visit
	} else {
		if prevVisit, ok = visitsMap[visit.Id]; ok {
			if prevVisit.User != visit.User {
				userChanged = true
				for key, visit := range visitsByUserMap[prevVisit.User] {
					if prevVisit == visit {
						array := visitsByUserMap[prevVisit.User]
						visitsByUserMap[prevVisit.User] = append(array[:key], array[key+1:]...)
					}
				}
			}
			if prevVisit.Location != visit.Location {
				locationChanged = true
				for key, visit := range visitsByLocationMap[prevVisit.Location] {
					if visit == prevVisit {
						array := visitsByLocationMap[prevVisit.Location]
						visitsByLocationMap[prevVisit.Location] = append(array[:key], array[key+1:]...)
					}
				}
			}
		}
		*prevRef = visit
	}

	if prevRef == nil || userChanged {
		//if _, ok := usersVisitsByVisitedAtMap[visit.User]; !ok {
		//	usersVisitsByVisitedAtMap[visit.User] = make(map[int]*Visit)
		//}
		//
		//usersVisitsByVisitedAtMap[visit.User][visit.Visited_at] = &visit

		visitsByUserMap[visit.User] = append(visitsByUserMap[visit.User], &visit)
	}
	if prevRef == nil || locationChanged {
		visitsByLocationMap[visit.Location] = append(visitsByLocationMap[visit.Location], &visit)
	}
}
