package main

import (
	"encoding/json"
	"errors"
	"github.com/valyala/fasthttp"
	"sync"
)

type Location struct {
	Id       uint   `json:"id"`
	Place    string `json:"place"`
	Country  string `json:"country"`
	City     string `json:"city"`
	Distance uint   `json:"distance"`
}

type LocationsMap struct {
	locations map[uint]*Location
	sync.RWMutex
}

func (l *LocationsMap) Get(id uint) *Location {
	l.RLock()
	defer l.RUnlock()

	return l.locations[id]
}

func (l *LocationsMap) Update(location Location) {
	l.Lock()
	l.locations[location.Id] = &location
	l.Unlock()
}

func getLocationRequestHandler(ctx *fasthttp.RequestCtx, entityId uint) {
	if location := locationsMap.Get(entityId); location != nil {
		response, _ := json.Marshal(location)
		ctx.Success("application/json", response)
		return
	}
	ctx.NotFound()
}

func createLocationRequestHandler(ctx *fasthttp.RequestCtx) {
	if location, err := createLocation(ctx.PostBody()); err == nil {
		ctx.SetConnectionClose()
		ctx.Success("application/json", []byte("{}"))

		go func() {
			locationsMap.Update(*location)
		}()

		return
	}
	ctx.Error("{}", 400)
	return

}

func updateLocationRequestHandler(ctx *fasthttp.RequestCtx, entityId uint) {
	if location := locationsMap.Get(entityId); location != nil {
		if updatedLocation, err := updateLocation(ctx.PostBody(), location); err == nil {
			ctx.SetConnectionClose()
			ctx.Success("application/json", []byte("{}"))

			go func() {
				locationsMap.Update(*updatedLocation)
			}()

			return
		}
		ctx.Error("{}", 400)
		return
	}
	ctx.NotFound()
}

func createLocation(postBody []byte) (*Location, error) {
	location := Location{}
	if err := json.Unmarshal(postBody, &location); err != nil {
		return nil, err
	}
	if location.Id == 0 || len(location.Place) == 0 || len(location.Country) == 0 ||
		len(location.City) == 0 {
		return nil, errors.New("Validation error")
	}
	if location := locationsMap.Get(location.Id); location != nil {
		return nil, errors.New("Location already exists")
	}

	return &location, nil
}

func updateLocation(postBody []byte, location *Location) (*Location, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(postBody, &data); err != nil {
		return nil, err
	}

	var updatedLocation Location
	updatedLocation = *location

	//updatedLocation := location

	if place, ok := data["place"]; ok {
		if place != nil {
			updatedLocation.Place = place.(string)
		} else {
			return nil, errors.New("Field validation error")
		}
	}
	if country, ok := data["country"]; ok {
		if country != nil {
			updatedLocation.Country = country.(string)
		} else {
			return nil, errors.New("Field validation error")
		}
	}
	if city, ok := data["city"]; ok {
		if city != nil {
			updatedLocation.City = city.(string)
		} else {
			return nil, errors.New("Field validation error")
		}
	}
	if distance, ok := data["distance"]; ok {
		if distance != nil {
			updatedLocation.Distance = uint(distance.(float64))
		} else {
			return nil, errors.New("Field validation error")
		}
	}

	return &updatedLocation, nil
}
