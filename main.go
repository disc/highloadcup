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

	locationsMap = make(map[uint]Location)
	usersMap     = make(map[uint]User)
	visitsMap    = make(map[uint]Visit)

	locationsCountryMap = make(map[string][]uint)

	visitsByUserMap     = make(map[uint][]Visit)
	visitsByLocationMap = make(map[uint][]Visit)
)

type validator interface {
	validate(isNew bool) bool
}

type User struct {
	Id         uint   `json:"id"`
	Email      string `json:"email"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Gender     string `json:"gender"`
	Birth_date int    `json:"birth_date"`
}

type UserRequestJson struct {
	Id         *uint
	Email      *string
	First_name *string
	Last_name  *string
	Gender     *string
	Birth_date *int
}

func (r *UserRequestJson) validate(isNew bool) bool {
	if isNew {
		return r.Id != nil && r.Email != nil && r.First_name != nil && r.Last_name != nil &&
			r.Gender != nil && r.Birth_date != nil
	} else {
		return r.Id == nil
	}
}

type Location struct {
	Id       uint   `json:"id"`
	Place    string `json:"place"`
	Country  string `json:"country"`
	City     string `json:"city"`
	Distance uint   `json:"distance"`
}

//type Nstring string
//
//func (n *Nstring) UnmarshalJSON(b []byte) (err error) {
//	if string(b) == "null" {
//		return nil
//	}
//	return json.Unmarshal(b, (*string)(n))
//}

type LocationRequestJson struct {
	Id       *uint
	Place    string
	Country  string
	City     string
	Distance uint
}

func (r *LocationRequestJson) validate(isNew bool) bool {
	var condition bool
	if isNew {
		condition = r.Id != nil && &r.Place != nil && &r.Country != nil && &r.City != nil &&
			&r.Distance != nil
	} else {
		condition = r.Id == nil
	}

	//REQUEST URI:/locations/308?query_id=999 BODY:{"city": null, "place": "\u0414\u043e\u043c"}
	//RESPONSE STATUS 200 != 400. BODY {} /
	//	46 requests (4.60%) failed

	//if &r.Place != nil {
	//	condition = condition && r.Place != ""
	//}
	//if &r.Country != nil {
	//	condition = condition && r.Country != ""
	//}

	if &r.City != nil /*&& &r.City != nil*/ {
		// {"city": null, "place": "\u0414\u043e\u043c"}
		// {"distance": 65, "country": "\u0421\u0428\u0410"}
		condition = condition && r.City != ""
	}
	//if &r.Distance != nil {
	//	condition = condition && r.Distance > 0
	//}

	return condition
}

type Visit struct {
	Id         uint `json:"id"`
	Location   uint `json:"location"`
	User       uint `json:"user"`
	Visited_at int  `json:"visited_at"`
	Mark       uint `json:"mark"`
}

type VisitRequestJson struct {
	Id         *uint
	Location   *uint
	User       *uint
	Visited_at *int
	Mark       *uint
}

func (r *VisitRequestJson) validate(isNew bool) bool {
	if isNew {
		return r.Id != nil && r.Location != nil && r.User != nil && r.Visited_at != nil &&
			r.Mark != nil
	} else {
		return r.Id == nil
	}
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

func updateUser(ctx *fasthttp.RequestCtx, isNew bool) error {
	user := User{}
	if err := json.Unmarshal(ctx.PostBody(), &user); err != nil {
		return err
	}
	if isNew {
		if user.Id == 0 || len(user.First_name) == 0 || len(user.Last_name) == 0 ||
			(user.Gender != "m" && user.Gender != "f")  {
			return errors.New("Validation error")
		}
		if _, ok := usersMap[user.Id]; isNew && ok {
			return errors.New("User already exists")
		}
		// todo; move to single method
		usersMap[user.Id] = user
	} else {
		user := usersMap[user.Id]
		var data map[string]interface{}
		if err := json.Unmarshal(ctx.PostBody(), &data); err != nil {
			return err
		}
		if email, ok := data["email"]; ok {
			if email != nil {
				user.Email = email.(string)
			} else {
				return errors.New("Field validation error")
			}
		}
		if firstName, ok := data["first_name"]; ok {
			if firstName != nil {
				user.First_name = firstName.(string)
			} else {
				return errors.New("Field validation error")
			}
		}
		if lastName, ok := data["last_name"]; ok {
			if lastName != nil {
				user.Last_name = lastName.(string)
			} else {
				return errors.New("Field validation error")
			}
		}
		if gender, ok := data["gender"]; ok {
			if gender != nil && (gender.(string) == "m" || gender.(string) == "f") {
				user.Gender = gender.(string)
			} else {
				return errors.New("Field validation error")
			}
		}
		if birthDate, ok := data["birth_date"]; ok {
			if birthDate != nil {
				user.Birth_date = int(birthDate.(float64))
			} else {
				return errors.New("Field validation error")
			}
		}
	}

	return nil
}

func updateVisit(ctx *fasthttp.RequestCtx, isNew bool) error {
	visit := Visit{}
	if err := json.Unmarshal(ctx.PostBody(), &visit); err != nil {
		return err
	}

	if isNew {
		if visit.Id == 0 || visit.Location == 0 || visit.User == 0 || visit.Mark > 5 {
			return errors.New("Validation error")
		}
		if _, ok := visitsMap[visit.Id]; isNew && ok {
			return errors.New("Visit already exists")
		}
		// todo; move to single method
		visitsMap[visit.Id] = visit
	} else {
		visit := visitsMap[visit.Id]
		var data map[string]interface{}
		if err := json.Unmarshal(ctx.PostBody(), &data); err != nil {
			return err
		}
		if location, ok := data["location"]; ok {
			if location != nil {
				visit.Location = uint(location.(float64))
			} else {
				return errors.New("Field validation error")
			}
		}
		if user, ok := data["user"]; ok {
			if user != nil {
				visit.User = uint(user.(float64))
			} else {
				return errors.New("Field validation error")
			}
		}
		if visitedAt, ok := data["visited_at"]; ok {
			if visitedAt != nil {
				visit.Visited_at = int(visitedAt.(float64))
			} else {
				return errors.New("Field validation error")
			}
		}
		if mark, ok := data["mark"]; ok {
			if mark != nil && uint(mark.(float64)) <= 5 {
				visit.Mark = uint(mark.(float64))
			} else {
				return errors.New("Field validation error")
			}
		}
	}

	return nil
}

func updateLocation(ctx *fasthttp.RequestCtx, isNew bool) error {
	location := Location{}
	if err := json.Unmarshal(ctx.PostBody(), &location); err != nil {
		return err
	}
	if isNew {
		if location.Id == 0 || len(location.Place) == 0 || len(location.Country) == 0 ||
			len(location.City) == 0 {
			return errors.New("Validation error")
		}
		if _, ok := locationsMap[location.Id]; ok {
			return errors.New("Location already exists")
		}
		// todo; move to single method
		locationsMap[location.Id] = location
	} else {
		location := locationsMap[location.Id]
		var data map[string]interface{}
		if err := json.Unmarshal(ctx.PostBody(), &data); err != nil {
			return err
		}
		if place, ok := data["place"]; ok {
			if place != nil {
				location.Place = place.(string)
			} else {
				return errors.New("Field validation error")
			}
		}
		if country, ok := data["country"]; ok {
			if country != nil {
				location.Country = country.(string)
			} else {
				return errors.New("Field validation error")
			}
		}
		if city, ok := data["city"]; ok {
			if city != nil {
				location.City = city.(string)
			} else {
				return errors.New("Field validation error")
			}
		}
		if distance, ok := data["distance"]; ok {
			if distance != nil {
				//fmt.Println(distance.(float64))
				location.Distance = uint(distance.(float64))
			} else {
				return errors.New("Field validation error")
			}
		}
	}

	return nil
}

func locationRequestHandler(ctx *fasthttp.RequestCtx) []byte {
	path := ctx.Path()
	isNew := path[len(path)-1] == 'w'

	location := Location{}
	if !isNew {
		var err error
		if location, err = getLocation(getEntityId(path)); err != nil {
			ctx.NotFound()
			return nil
		}
	}

	if ctx.IsGet() {
		response, _ := json.Marshal(location)
		return response
	} else if ctx.IsPost() {
		// create or update
		if err := updateLocation(ctx, isNew); err != nil {
			ctx.Error("{}", 400)
		} else {
			return []byte("{}")
		}
	}

	return nil
}

func usersRequestHandler(ctx *fasthttp.RequestCtx) []byte {
	path := ctx.Path()
	isNew := path[len(path)-1] == 'w'

	user := User{}
	if !isNew {
		var err error
		if user, err = getUser(getEntityId(path)); err != nil {
			ctx.NotFound()
			return nil
		}
	}

	if ctx.IsGet() {
		response, _ := json.Marshal(user)
		return response
	} else if ctx.IsPost() {
		// create or update
		if err := updateUser(ctx, isNew); err != nil {
			ctx.Error("{}", 400)
		} else {
			return []byte("{}")
		}
	}

	return nil
}

func visitsRequestHandler(ctx *fasthttp.RequestCtx) []byte {
	path := ctx.Path()
	isNew := path[len(path)-1] == 'w'

	visit := Visit{}
	if !isNew {
		var err error
		if visit, err = getVisits(getEntityId(path)); err != nil {
			ctx.NotFound()
			return nil
		}
	}

	if ctx.IsGet() {
		response, _ := json.Marshal(visit)
		return response
	} else if ctx.IsPost() {
		// create or update
		if err := updateVisit(ctx, isNew); err != nil {
			ctx.Error("{}", 400)
		} else {
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
