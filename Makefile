build:
	docker build -t golang-app .
run-local: build
	docker run --rm -p 8080:80 -v ~/Downloads/data-train.zip:/tmp/data/data.zip -t golang-app
deploy: build
	docker tag golang-app stor.highloadcup.ru/travels/little_eagle
	docker push stor.highloadcup.ru/travels/little_eagle
tank:
	cd $$HOME/workspace/hlcupdocs/ && $$HOME/workspace/hlcupdocs/start.sh
kill:
	docker kill $$(docker ps -q) &> /dev/null;

test-phase-1:
	$$GOBIN/highloadcup_tester -addr http://127.0.0.1:8080 -hlcupdocs $$HOME/workspace/hlcupdocs/data/TRAIN/ -test -phase 1
test-phase-2:
	$$GOBIN/highloadcup_tester -addr http://127.0.0.1:8080 -hlcupdocs $$HOME/workspace/hlcupdocs/data/TRAIN/ -test -phase 2
test-phase-3:
	$$GOBIN/highloadcup_tester -addr http://127.0.0.1:8080 -hlcupdocs $$HOME/workspace/hlcupdocs/data/TRAIN/ -test -phase 3
tests: test-phase-1 test-phase-2 test-phase-3

app-run: app-unzip
	/go/bin/highloadcup
app-unzip:
	mkdir -p $$(pwd)/data/ > /dev/null
	unzip -oq /tmp/data/data.zip -d $$(pwd)/data/

bench:
	go test -bench=.