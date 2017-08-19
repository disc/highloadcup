package main

import (
	"flag"
	"fmt"
	"log"

	"bytes"
	"encoding/json"
	"errors"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var (
	addr     = flag.String("addr", ":80", "TCP address to listen to")
	compress = flag.Bool("compress", false, "Whether to enable transparent response compression")

	locationsMap = make(map[uint32]Location)
	usersMap     = make(map[uint32]User)
	visitsMap    = make(map[uint32]Visit)

	visitsByUserMap = make(map[uint32][]Visit)
)

type User struct {
	Id         uint32 `json:"id"`
	Email      string `json:"email"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Gender     string `json:"gender"`
	Birth_date int32  `json:"birth_date"`
}

type Location struct {
	Id       uint32 `json:"id"`
	Place    string `json:"place"`
	Country  string `json:"country"`
	City     string `json:"city"`
	Distance uint   `json:"distance"`
}

type Visit struct {
	Id         uint32 `json:"id"`
	Location   uint32 `json:"location"`
	User       uint32 `json:"user"`
	Visited_at int32  `json:"visited_at"`
	Mark       uint   `json:"mark"`
}

type UserVisits struct {
	Visits []UserVisit `json:"visits"`
}

type UserVisit struct {
	Mark       uint   `json:"mark"`
	Visited_at int32  `json:"visited_at"`
	Place      string `json:"place"`
}

type LocationAvg struct {
	Avg float32 `json:"avg"`
}

func parseLocations(fileBytes []byte) {
	type jsonKey struct {
		Locations []Location
	}

	var locations jsonKey
	json.Unmarshal(fileBytes, &locations)

	for _, loc := range locations.Locations {
		locationsMap[loc.Id] = loc
	}
}

func parseVisits(fileBytes []byte) {
	type jsonKey struct {
		Visits []Visit
	}

	var visits jsonKey
	json.Unmarshal(fileBytes, &visits)

	for _, visit := range visits.Visits {
		visitsMap[visit.Id] = visit

		visitsByUserMap[visit.User] = append(visitsByUserMap[visit.User], visit)
	}
}

func parseUsers(fileBytes []byte) {
	type jsonKey struct {
		Users []User
	}

	var users jsonKey
	json.Unmarshal(fileBytes, &users)

	for _, user := range users.Users {
		usersMap[user.Id] = user
	}
}

func parseFile(filename string) {
	rawData, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if strings.LastIndex(filename, "users_") != -1 {
		parseUsers(rawData)
	} else if strings.LastIndex(filename, "locations_") != -1 {
		parseLocations(rawData)
	} else if strings.LastIndex(filename, "visits_") != -1 {
		parseVisits(rawData)
	}
}

func main() {
	fmt.Println("Started")

	parseFile("./locations_1.json")
	parseFile("./users_1.json")
	parseFile("./visits_1.json")

	fmt.Println("Parsing completed")

	flag.Parse()

	h := requestHandler
	if *compress {
		h = fasthttp.CompressHandler(h)
	}

	if err := fasthttp.ListenAndServe(*addr, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func getEntityId(path []byte) uint32 {
	from := bytes.IndexByte(path[1:], '/')
	to := bytes.IndexByte(path[from+2:], '/')

	if to == -1 {
		to = len(path)
	} else {
		to += from + 2
	}

	entityId, _ := strconv.Atoi(string(path[from+2 : to]))

	return uint32(entityId)
}

func getUser(id uint32) (User, error) {
	if user, ok := usersMap[id]; ok {
		return user, nil
	}
	return User{}, errors.New("user not found")
}

func getLocation(id uint32) (Location, error) {
	if location, ok := locationsMap[id]; ok {
		return location, nil
	}
	return Location{}, errors.New("location not found")
}

func getVisits(id uint32) (Visit, error) {
	if visit, ok := visitsMap[id]; ok {
		return visit, nil
	}
	return Visit{}, errors.New("visit not found")
}

func getUserVisits(userId uint32) []UserVisit {
	var userVisits []UserVisit
	for _, visit := range visitsByUserMap[userId] {
		userVisits = append(userVisits, UserVisit{visit.Mark, visit.Visited_at, "Moscow"})
	}

	return userVisits
}

func getLocationAvg(locationId uint32) float32 {
	return 1.23
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetConnectionClose()

	path := ctx.Request.URI().Path()

	var response []byte

	if path[1] == 'l' && path[len(path)-1] == 'g' {
		locationId := getEntityId(path)
		// get avg value
		if _, err := getLocation(locationId); err != nil {
			ctx.NotFound()
			return
		} else {
			response, _ = json.Marshal(LocationAvg{getLocationAvg(locationId)})
		}

	} else if path[1] == 'u' && path[len(path)-1] == 's' && len(path) >= 14 {
		userId := getEntityId(path)
		// get user visits
		if _, err := getUser(userId); err != nil {
			ctx.NotFound()
			return
		} else {
			response, _ = json.Marshal(UserVisits{getUserVisits(userId)})
		}
	} else if path[1] == 'l' && path[9] == 's' {
		// get location
		if location, err := getLocation(getEntityId(path)); err != nil {
			ctx.NotFound()
			return
		} else {
			response, _ = json.Marshal(location)
		}
	} else if path[1] == 'u' && path[5] == 's' {
		// get user
		if user, err := getUser(getEntityId(path)); err != nil {
			ctx.NotFound()
			return
		} else {
			response, _ = json.Marshal(user)
		}
	} else if path[1] == 'v' && path[6] == 's' {
		// get visit
		if visit, err := getVisits(getEntityId(path)); err != nil {
			ctx.NotFound()
			return
		} else {
			response, _ = json.Marshal(visit)
		}
	} else {
		ctx.NotFound()
		return
	}

	if len(response) > 0 {
		ctx.SetBody(response)
		ctx.SetContentType("application/json; charset=utf8")
	}
}
