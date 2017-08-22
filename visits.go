package main

import (
	"encoding/json"
	"errors"
	"github.com/valyala/fasthttp"
	"sync"
)

type Visit struct {
	Id         uint `json:"id"`
	Location   uint `json:"location"`
	User       uint `json:"user"`
	Visited_at int  `json:"visited_at"`
	Mark       uint `json:"mark"`
}

type VisitsMap struct {
	visits map[uint]*Visit
	sync.RWMutex
}

func (v *VisitsMap) Get(id uint) *Visit {
	v.RLock()
	defer v.RUnlock()

	return v.visits[id]
}

func (v *VisitsMap) Update(visit Visit) {
	v.Lock()
	v.visits[visit.Id] = &visit
	v.Unlock()
}

func getVisitRequestHandler(ctx *fasthttp.RequestCtx, entityId uint) {
	if visit := visitsMap.Get(entityId); visit != nil {
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
	if visit := visitsMap.Get(entityId); visit != nil {
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
	if visit := visitsMap.Get(visit.Id); visit != nil {
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
		userChanged     = false
		locationChanged = false
	)
	if prevRef == nil {
		visitsMap.Update(visit)
	} else {
		if prevVisit := visitsMap.Get(visit.Id); prevVisit != nil {
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
		visitsByUserMap[visit.User] = append(visitsByUserMap[visit.User], &visit)
	}
	if prevRef == nil || locationChanged {
		visitsByLocationMap[visit.Location] = append(visitsByLocationMap[visit.Location], &visit)
	}
}
