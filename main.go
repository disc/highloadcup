package main

import (
	"flag"
	"fmt"
	"log"

	"bytes"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var (
	addr = flag.String("addr", ":80", "TCP address to listen to")

	locationsMap = make(map[uint]*Location)
	usersMap     = make(map[uint]*User)
	visitsMap    = make(map[uint]*Visit)

	visitsByUserMap     = make(map[uint][]*Visit)
	visitsByLocationMap = make(map[uint][]*Visit)

	//usersVisitsByVisitedAtMap = make(map[uint]map[int]*Visit)
)

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

func parseDataDir(dirPath string) {
	files, _ := ioutil.ReadDir(dirPath)
	for _, f := range files {
		parseFile(dirPath + f.Name())
	}
}

func main() {
	fmt.Println("Started")

	parseDataDir("./data/")

	fmt.Println("Parsing completed")

	flag.Parse()

	h := requestHandler

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

	entityId, _ := strconv.ParseUint(string(path[from+2:to]), 0, 32)

	return uint(entityId)
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	path := ctx.Path()

	if path[1] == 'l' && path[len(path)-1] == 'g' {
		locationAvgRequestHandler(ctx, getEntityId(path), ctx.QueryArgs())
		return
	}

	if path[1] == 'u' && path[len(path)-1] == 's' && len(path) >= 14 {
		userVisitsRequestHandler(ctx, getEntityId(path), ctx.QueryArgs())
		return
	}

	if path[1] == 'v' && path[6] == 's' {
		if path[len(path)-1] == 'w' {
			createVisitRequestHandler(ctx)
		} else {
			if ctx.IsGet() {
				getVisitRequestHandler(ctx, getEntityId(path))
			} else {
				updateVisitRequestHandler(ctx, getEntityId(path))
			}
		}
		return
	}

	if path[1] == 'u' && path[5] == 's' {
		if path[len(path)-1] == 'w' {
			createUserRequestHandler(ctx)
		} else {
			id := getEntityId(path)
			if ctx.IsGet() {
				getUserRequestHandler(ctx, id)
			} else {
				updateUserRequestHandler(ctx, id)
			}
		}
		return
	}

	if path[1] == 'l' && path[9] == 's' {
		if path[len(path)-1] == 'w' {
			createLocationRequestHandler(ctx)
		} else {
			id := getEntityId(path)
			if ctx.IsGet() {
				getLocationRequestHandler(ctx, id)
			} else {
				updateLocationRequestHandler(ctx, id)
			}
		}
		return
	}
}
