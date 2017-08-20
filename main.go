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
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	addr     = flag.String("addr", ":80", "TCP address to listen to")
	compress = flag.Bool("compress", false, "Whether to enable transparent response compression")

	locationsMap = make(map[uint]*Location)
	usersMap     = make(map[uint]*User)
	visitsMap    = make(map[uint]*Visit)

	visitsByUserMap     = make(map[uint][]*Visit)
	visitsByLocationMap = make(map[uint][]*Visit)
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
		updateLocationMaps(loc)
	}
}

func updateLocationMaps(location Location) {
	locationsMap[location.Id] = &location
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
		visitsByUserMap[visit.User] = append(visitsByUserMap[visit.User], &visit)
	}
	if prevRef == nil || locationChanged {
		visitsByLocationMap[visit.Location] = append(visitsByLocationMap[visit.Location], &visit)
	}
}

func updateUsersMaps(user User) {
	usersMap[user.Id] = &user
}

func parseVisits(fileBytes []byte) {
	type jsonKey struct {
		Visits []Visit
	}

	var visits jsonKey
	json.Unmarshal(fileBytes, &visits)

	for _, visit := range visits.Visits {
		updateVisitsMaps(visit, nil)
	}
}

func parseUsers(fileBytes []byte) {
	type jsonKey struct {
		Users []User
	}

	var users jsonKey
	json.Unmarshal(fileBytes, &users)

	for _, user := range users.Users {
		updateUsersMaps(user)
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

func getUser(id uint) (*User, error) {
	if user, ok := usersMap[id]; ok {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func getLocation(id uint) (*Location, error) {
	if location, ok := locationsMap[id]; ok {
		return location, nil
	}
	return nil, errors.New("location not found")
}

func getVisit(id uint) (*Visit, error) {
	if visit, ok := visitsMap[id]; ok {
		return visit, nil
	}
	return nil, errors.New("visit not found")
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
		if filters.country != nil && locationsMap[visit.Location].Country != *filters.country {
			continue
		}
		if filters.toDistance != nil && locationsMap[visit.Location].Distance >= *filters.toDistance {
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
	return int(time.Now().Unix()) - (*age)*int(math.Floor(365.24*24*60*60))
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

func locationAvgRequestHandler(ctx *fasthttp.RequestCtx) []byte {
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
		//TODO: Fix toDistance results
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

func createUser(ctx *fasthttp.RequestCtx) error {
	user := User{}
	if err := json.Unmarshal(ctx.PostBody(), &user); err != nil {
		return err
	}

	if user.Id == 0 || len(user.First_name) == 0 || len(user.Last_name) == 0 ||
		(user.Gender != "m" && user.Gender != "f") {
		return errors.New("Validation error")
	}
	if _, ok := usersMap[user.Id]; ok {
		return errors.New("User already exists")
	}

	updateUsersMaps(user)

	return nil
}

func updateUser(ctx *fasthttp.RequestCtx, user *User) error {
	var data map[string]interface{}
	if err := json.Unmarshal(ctx.PostBody(), &data); err != nil {
		return err
	}

	var updatedUser User
	updatedUser = *user

	if email, ok := data["email"]; ok {
		if email != nil {
			updatedUser.Email = email.(string)
		} else {
			return errors.New("Field validation error")
		}
	}
	if firstName, ok := data["first_name"]; ok {
		if firstName != nil {
			updatedUser.First_name = firstName.(string)
		} else {
			return errors.New("Field validation error")
		}
	}
	if lastName, ok := data["last_name"]; ok {
		if lastName != nil {
			updatedUser.Last_name = lastName.(string)
		} else {
			return errors.New("Field validation error")
		}
	}
	if gender, ok := data["gender"]; ok {
		if gender != nil && (gender.(string) == "m" || gender.(string) == "f") {
			updatedUser.Gender = gender.(string)
		} else {
			return errors.New("Field validation error")
		}
	}
	if birthDate, ok := data["birth_date"]; ok {
		if birthDate != nil {
			updatedUser.Birth_date = int(birthDate.(float64))
		} else {
			return errors.New("Field validation error")
		}
	}

	updateUsersMaps(updatedUser)

	return nil
}

func createVisit(ctx *fasthttp.RequestCtx) error {
	visit := Visit{}
	if err := json.Unmarshal(ctx.PostBody(), &visit); err != nil {
		return err
	}

	if visit.Id == 0 || visit.Location == 0 || visit.User == 0 || visit.Mark > 5 {
		return errors.New("Validation error")
	}
	if _, ok := visitsMap[visit.Id]; ok {
		return errors.New("Visit already exists")
	}

	updateVisitsMaps(visit, nil)

	return nil
}

func updateVisit(ctx *fasthttp.RequestCtx, visit *Visit) error {
	var data map[string]interface{}
	if err := json.Unmarshal(ctx.PostBody(), &data); err != nil {
		return err
	}
	var updatedVisit Visit
	updatedVisit = *visit

	if location, ok := data["location"]; ok {
		if location != nil {
			updatedVisit.Location = uint(location.(float64))
		} else {
			return errors.New("Field validation error")
		}
	}
	if user, ok := data["user"]; ok {
		if user != nil {
			updatedVisit.User = uint(user.(float64))
		} else {
			return errors.New("Field validation error")
		}
	}
	if visitedAt, ok := data["visited_at"]; ok {
		if visitedAt != nil {
			updatedVisit.Visited_at = int(visitedAt.(float64))
		} else {
			return errors.New("Field validation error")
		}
	}
	if mark, ok := data["mark"]; ok {
		if mark != nil && uint(mark.(float64)) <= 5 {
			updatedVisit.Mark = uint(mark.(float64))
		} else {
			return errors.New("Field validation error")
		}
	}

	updateVisitsMaps(updatedVisit, visit)

	return nil
}

func createLocation(ctx *fasthttp.RequestCtx) error {
	location := Location{}
	if err := json.Unmarshal(ctx.PostBody(), &location); err != nil {
		return err
	}
	if location.Id == 0 || len(location.Place) == 0 || len(location.Country) == 0 ||
		len(location.City) == 0 {
		return errors.New("Validation error")
	}
	if _, ok := locationsMap[location.Id]; ok {
		return errors.New("Location already exists")
	}

	updateLocationMaps(location)

	return nil
}

func updateLocation(ctx *fasthttp.RequestCtx, location *Location) error {
	var data map[string]interface{}
	if err := json.Unmarshal(ctx.PostBody(), &data); err != nil {
		return err
	}

	var updatedLocation Location
	updatedLocation = *location

	if place, ok := data["place"]; ok {
		if place != nil {
			updatedLocation.Place = place.(string)
		} else {
			return errors.New("Field validation error")
		}
	}
	if country, ok := data["country"]; ok {
		if country != nil {
			updatedLocation.Country = country.(string)
		} else {
			return errors.New("Field validation error")
		}
	}
	if city, ok := data["city"]; ok {
		if city != nil {
			updatedLocation.City = city.(string)
		} else {
			return errors.New("Field validation error")
		}
	}
	if distance, ok := data["distance"]; ok {
		if distance != nil {
			//fmt.Println(distance.(float64))
			updatedLocation.Distance = uint(distance.(float64))
		} else {
			return errors.New("Field validation error")
		}
	}
	updateLocationMaps(updatedLocation)

	return nil
}

func locationRequestHandler(ctx *fasthttp.RequestCtx) []byte {
	path := ctx.Path()
	isNew := path[len(path)-1] == 'w'

	entityId := getEntityId(path)

	var (
		location *Location
		err      error
	)

	if location, err = getLocation(entityId); !isNew && err != nil {
		ctx.NotFound()
		return nil
	}

	if ctx.IsGet() {
		response, _ := json.Marshal(location)
		return response
	}

	if !isNew {
		if location, err := getLocation(entityId); err != nil {
			ctx.NotFound()
			return nil
		} else if ctx.IsGet() {
			response, _ := json.Marshal(location)
			return response
		}
	}

	if ctx.IsPost() {
		if isNew {
			if err := createLocation(ctx); err != nil {
				ctx.Error("{}", 400)
			} else {
				return []byte("{}")
			}
		} else {
			if err := updateLocation(ctx, location); err != nil {
				ctx.Error("{}", 400)
			} else {
				return []byte("{}")
			}
		}
	}

	return nil
}

func usersRequestHandler(ctx *fasthttp.RequestCtx) []byte {
	path := ctx.Path()
	isNew := path[len(path)-1] == 'w'

	entityId := getEntityId(path)

	var (
		user *User
		err  error
	)

	if user, err = getUser(entityId); !isNew && err != nil {
		ctx.NotFound()
		return nil
	}

	if ctx.IsGet() {
		response, _ := json.Marshal(user)
		return response
	}

	if !isNew {
		if user, err := getUser(entityId); err != nil {
			ctx.NotFound()
			return nil
		} else if ctx.IsGet() {
			response, _ := json.Marshal(user)
			return response
		}
	}

	if ctx.IsPost() {
		if isNew {
			if err := createUser(ctx); err != nil {
				ctx.Error("{}", 400)
			} else {
				return []byte("{}")
			}
		} else {
			if err := updateUser(ctx, user); err != nil {
				ctx.Error("{}", 400)
			} else {
				return []byte("{}")
			}
		}
	}

	return nil
}

func visitsRequestHandler(ctx *fasthttp.RequestCtx) []byte {
	path := ctx.Path()
	isNew := path[len(path)-1] == 'w'

	entityId := getEntityId(path)

	var (
		visit *Visit
		err   error
	)

	if visit, err = getVisit(entityId); !isNew && err != nil {
		ctx.NotFound()
		return nil
	}

	if ctx.IsGet() {
		response, _ := json.Marshal(visit)
		return response
	}

	if !isNew {
		if visit, err := getVisit(entityId); err != nil {
			ctx.NotFound()
			return nil
		} else if ctx.IsGet() {
			response, _ := json.Marshal(visit)
			return response
		}
	}

	if ctx.IsPost() {
		if isNew {
			if err := createVisit(ctx); err != nil {
				ctx.Error("{}", 400)
			} else {
				return []byte("{}")
			}
		} else {
			if err := updateVisit(ctx, visit); err != nil {
				ctx.Error("{}", 400)
			} else {
				return []byte("{}")
			}
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
