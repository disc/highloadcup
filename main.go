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
	"time"
	"bufio"
	"runtime/debug"
)

var (
	addr = flag.String("addr", ":80", "TCP address to listen to")

	locationsMap = LocationsMap{locations: make(map[uint]*Location)}
	usersMap     = UsersMap{users: make(map[uint]*User)}
	visitsMap    = VisitsMap{visits: make(map[uint]*Visit)}

	visitsByUserMap     = make(map[uint][]*Visit)
	visitsByLocationMap = make(map[uint][]*Visit)

	now = int(time.Now().Unix())

	needToDisableGC = true
	postsCount = 0
	maxPostsCount = 0
	isTestMode = true
)

func parseLocations(fileBytes []byte) {
	type jsonKey struct {
		Locations []Location
	}

	var locations jsonKey
	json.Unmarshal(fileBytes, &locations)

	for _, location := range locations.Locations {
		locationsMap.Update(location)
	}
}

func parseVisits(fileBytes []byte) {
	type jsonKey struct {
		Visits []Visit
	}

	var visits jsonKey
	json.Unmarshal(fileBytes, &visits)

	for _, visit := range visits.Visits {
		visitsMap.Update(visit, nil)
	}
}

func parseUsers(fileBytes []byte) {
	type jsonKey struct {
		Users []User
	}

	var users jsonKey
	json.Unmarshal(fileBytes, &users)

	for _, user := range users.Users {
		usersMap.Update(user)
	}
}

func parseOptions(filename string) {
	if file, err := os.OpenFile(filename, os.O_RDONLY, 0644); err == nil {
		reader := bufio.NewReader(file)
		if line, _, err := reader.ReadLine(); err == nil {
			now, _ = strconv.Atoi(string(line))
			fmt.Println("`Now` was updated from options.txt", now)
		}
		if line, _, err := reader.ReadLine(); err == nil {
			mode, _ := strconv.Atoi(string(line))
			if mode == 1 {
				isTestMode = false
			}
		}
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
	} else if strings.LastIndex(filename, "options.txt") != -1 {
		parseOptions(filename)
	}
}

func parseDataDir(dirPath string) {
	files, _ := ioutil.ReadDir(dirPath)
	for _, f := range files {
		parseFile(dirPath + f.Name())
	}
}

func main() {
	start := time.Now()
	fmt.Println("Started")

	parseDataDir("./data/")

	fmt.Println("Parsing completed at " + time.Since(start).String())

	var runMode string
	if isTestMode {
		runMode = "train"
		maxPostsCount = 3000
	} else {
		runMode = "full"
		maxPostsCount = 39999
	}
	fmt.Println("Running mode: " + runMode)

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

	isGetRequest := ctx.IsGet()

	if !isGetRequest {
		postsCount++
		if postsCount == maxPostsCount {
			go disableGC()
		}
	}

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
			if isGetRequest {
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
			if isGetRequest {
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
			if isGetRequest {
				getLocationRequestHandler(ctx, id)
			} else {
				updateLocationRequestHandler(ctx, id)
			}
		}
		return
	}
}

func disableGC() {
	if needToDisableGC {
		debug.SetGCPercent(-1)
		fmt.Println("GC was disabled")
		needToDisableGC = false
	}
}