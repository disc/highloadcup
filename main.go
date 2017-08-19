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
	"math"
	"os"
	"strconv"
	"strings"
	"time"
	"sort"
)

var (
	addr     = flag.String("addr", ":80", "TCP address to listen to")
	compress = flag.Bool("compress", false, "Whether to enable transparent response compression")

	locationsMap = make(map[uint]Location)
	usersMap     = make(map[uint]User)
	visitsMap    = make(map[uint]Visit)

	locationsCountryMap = make(map[string][]uint)

	visitsByUserMap     = make(map[uint][]Visit)
	visitsByLocationMap = make(map[uint][]Visit)
)

type User struct {
	Id         uint   `json:"id"`
	Email      string `json:"email"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Gender     string `json:"gender"`
	Birth_date int    `json:"birth_date"`
}

type Location struct {
	Id       uint   `json:"id"`
	Place    string `json:"place"`
	Country  string `json:"country"`
	City     string `json:"city"`
	Distance uint   `json:"distance"`
}

type Visit struct {
	Id         uint `json:"id"`
	Location   uint `json:"location"`
	User       uint `json:"user"`
	Visited_at int  `json:"visited_at"`
	Mark       uint `json:"mark"`
}

type UserVisits struct {
	Visits []UserVisit `json:"visits"`
}

type UserVisit struct {
	Mark       uint   `json:"mark"`
	Visited_at int    `json:"visited_at"`
	Place      string `json:"place"`
}

type LocationAvg struct {
	Avg float64 `json:"avg"`
}

type UserVisitsFilter struct {
	fromDate   *int
	toDate     *int
	country    *string
	toDistance *uint
}

type LocationAvgFilter struct {
	fromDate *int
	toDate   *int
	fromAge  *int
	toAge    *int
	gender   *string
}

func parseLocations(fileBytes []byte) {
	type jsonKey struct {
		Locations []Location
	}

	var locations jsonKey
	json.Unmarshal(fileBytes, &locations)

	for _, loc := range locations.Locations {
		locationsMap[loc.Id] = loc
		locationsCountryMap[loc.Country] = append(locationsCountryMap[loc.Country], loc.Id)
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
		visitsByLocationMap[visit.Location] = append(visitsByLocationMap[visit.Location], visit)
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

func parseDatDir(dirPath string) {
	files, _ := ioutil.ReadDir(dirPath)
	for _, f := range files {
		parseFile(dirPath + f.Name())
	}
}

func main() {
	fmt.Println("Started")

	parseDatDir("./data/")

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

func getEntityId(path []byte) uint {
	from := bytes.IndexByte(path[1:], '/')
	to := bytes.IndexByte(path[from+2:], '/')

	if to == -1 {
		to = len(path)
	} else {
		to += from + 2
	}

	entityId, _ := strconv.Atoi(string(path[from+2 : to]))

	return uint(entityId)
}

func getUser(id uint) (User, error) {
	if user, ok := usersMap[id]; ok {
		return user, nil
	}
	return User{}, errors.New("user not found")
}

func getLocation(id uint) (Location, error) {
	if location, ok := locationsMap[id]; ok {
		return location, nil
	}
	return Location{}, errors.New("location not found")
}

func getVisits(id uint) (Visit, error) {
	if visit, ok := visitsMap[id]; ok {
		return visit, nil
	}
	return Visit{}, errors.New("visit not found")
}

func getUserVisits(userId uint, filters UserVisitsFilter /*fromDate *int*/) []UserVisit {
	var userVisits = make(map[int]UserVisit, 0)
	for _, visit := range visitsByUserMap[userId] {
		if filters.fromDate != nil && visit.Visited_at < *filters.fromDate {
			continue
		}
		if filters.toDate != nil && visit.Visited_at > *filters.toDate {
			continue
		}
		if filters.country != nil && locationsMap[visit.Location].Country != *filters.country {
			continue
		}
		if filters.toDistance != nil && locationsMap[visit.Location].Distance > *filters.toDistance {
			continue
		}
		userVisits[visit.Visited_at] = UserVisit{visit.Mark, visit.Visited_at, locationsMap[visit.Location].Place}
	}

	visitedAtList := make([]int, 0, len(userVisits))
	for visitedAt := range userVisits {
		visitedAtList = append(visitedAtList, visitedAt)
	}
	sort.Ints(visitedAtList) //sort by key

	sortedUserVisits := make([]UserVisit, 0, len(userVisits))
	for _, visitedAt := range visitedAtList {
		sortedUserVisits = append(sortedUserVisits, userVisits[visitedAt])
	}

	return sortedUserVisits
}

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	_div := math.Copysign(div, val)
	_roundOn := math.Copysign(roundOn, val)
	if _div >= _roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func getTimestampByAge(age *int) int {
	return int(time.Now().Unix()) - int((*age+1)*365*24*60*60)
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
		if filters.fromAge != nil || filters.toAge != nil {
			// timestamp -  ((age+1) * 365.24 * 24 * 60 * 60)

			if filters.fromAge != nil && usersMap[visit.User].Birth_date > getTimestampByAge(filters.fromAge) {
				continue
			}
			if filters.toAge != nil && usersMap[visit.User].Birth_date < getTimestampByAge(filters.toAge) {
				continue
			}
		}
		if filters.gender != nil && *filters.gender != usersMap[visit.User].Gender {
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

func detectRoute(path []byte) string {
	if path[1] == 'l' && path[len(path)-1] == 'g' {
		return "locationAvg"
	} else if path[1] == 'u' && path[len(path)-1] == 's' && len(path) >= 14 {
		return "userVisits"
	} else if path[1] == 'l' && path[9] == 's' {
		return "locations"
	} else if path[1] == 'u' && path[5] == 's' {
		return "users"
	} else if path[1] == 'v' && path[6] == 's' {
		return "visits"
	}

	return ""
}

func locationAvgRequestHandler(ctx *fasthttp.RequestCtx) []byte  {
	path := ctx.Path()
	query := ctx.QueryArgs()

	locationId := getEntityId(path)
	// get avg value
	if _, err := getLocation(locationId); err != nil {
		ctx.NotFound()
		return nil
	} else {
		var filters = LocationAvgFilter{}
		if fromDate := query.Has("fromDate"); fromDate {
			if fromDateInt, err := strconv.Atoi(string(query.Peek("fromDate"))); err != nil {
				ctx.Error("{}", 400)
				return nil
			} else {
				filters.fromDate = &fromDateInt
			}
		}
		if toDate := query.Has("toDate"); toDate {
			if toDateInt, err := strconv.Atoi(string(query.Peek("toDate"))); err != nil {
				ctx.Error("{}", 400)
				return nil
			} else {
				filters.toDate = &toDateInt
			}
		}
		if fromAge := query.Has("fromAge"); fromAge {
			//todo: validate + 400 if wrong param
			if fromAgeInt, err := strconv.Atoi(string(query.Peek("fromAge"))); err != nil {
				ctx.Error("{}", 400)
				return nil
			} else {
				filters.fromAge = &fromAgeInt
			}
		}
		if toAge := query.Has("toAge"); toAge {
			if toAgeInt, err := strconv.Atoi(string(query.Peek("toAge"))); err != nil {
				ctx.Error("{}", 400)
				return nil
			} else {
				filters.toAge = &toAgeInt
			}
		}
		if gender := query.Has("gender"); gender {
			if genderStr := string(query.Peek("gender")); len(genderStr) > 0 && (genderStr == "m" || genderStr == "f") {
				filters.gender = &genderStr
			} else {
				ctx.Error("{}", 400)
				return nil
			}
		}
		response, _ := json.Marshal(LocationAvg{Round(getLocationAvg(locationId, filters), .5, 5)})

		return response
	}
}

func userVisitsRequestHandler(ctx *fasthttp.RequestCtx) []byte {
	path := ctx.Path()
	query := ctx.QueryArgs()

	userId := getEntityId(path)
	// get user visits
	if _, err := getUser(userId); err != nil {
		ctx.NotFound()
		return nil
	} else {
		var filters = UserVisitsFilter{}
		if fromDate := query.Has("fromDate"); fromDate {
			//todo: validate + 400 if wrong param
			if fromDateInt, err := strconv.Atoi(string(query.Peek("fromDate"))); err != nil {
				ctx.Error("{}", 400)
				return nil
			} else {
				filters.fromDate = &fromDateInt
			}
		}
		if toDate := query.Has("toDate"); toDate {
			if toDateInt, err := strconv.Atoi(string(query.Peek("toDate"))); err != nil {
				ctx.Error("{}", 400)
				return nil
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
				return nil
			} else {
				distanceInt := uint(distanceInt)
				filters.toDistance = &distanceInt
			}
		}

		response, _ := json.Marshal(UserVisits{getUserVisits(userId, filters)})

		return response
	}
}

func locationRequestHandler(ctx *fasthttp.RequestCtx) []byte {
	path := ctx.Path()
	if location, err := getLocation(getEntityId(path)); err != nil {
		ctx.NotFound()
		return nil
	} else {
		if ctx.IsGet() {
			response, _ := json.Marshal(location)
			return response
		} else if ctx.IsPost() {
			// create or update
			return []byte("{}")
		}
	}
	return nil
}

func usersRequestHandler(ctx *fasthttp.RequestCtx) []byte {
	path := ctx.Path()
	if user, err := getUser(getEntityId(path)); err != nil {
		ctx.NotFound()
		return nil
	} else {
		if ctx.IsGet() {
			response, _ := json.Marshal(user)
			return response
		} else if ctx.IsPost() {
			// create or update
			return []byte("{}")
		}
	}
	return nil
}

func visitsRequestHandler(ctx *fasthttp.RequestCtx) []byte {
	path := ctx.Path()
	if visit, err := getVisits(getEntityId(path)); err != nil {
		ctx.NotFound()
		return nil
	} else {
		if ctx.IsGet() {
			response, _ := json.Marshal(visit)
			return response
		} else if ctx.IsPost() {
			// create or update
			return []byte("{}")
		}
	}
	return nil
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	path := ctx.Path()

	var response []byte

	switch detectRoute(path) {
	case "locationAvg":
		if ctx.IsGet() {
			response = locationAvgRequestHandler(ctx)
		} else {
			ctx.NotFound()
		}
	case "userVisits":
		if ctx.IsGet() {
			response = userVisitsRequestHandler(ctx)
		} else {
			ctx.NotFound()
		}
	case "users":
		response = usersRequestHandler(ctx)
	case "visits":
		response = visitsRequestHandler(ctx)
	case "locations":
		response = locationRequestHandler(ctx)
	default:
		ctx.NotFound()
		return
	}

	if len(response) > 0 {
		fmt.Fprintf(ctx, string(response))
		ctx.SetContentType("application/json; charset=utf8")
	}
	ctx.SetConnectionClose()
}
