run-local: build
#	if $$(docker ps -q); then
#		docker kill $$(docker ps -q) &> /dev/null;
#	fi
	docker run --rm -p 8080:80 -v ~/Downloads/data5.zip:/tmp/data/data.zip -t golang-app
run:
	app
build:
	docker build -t golang-app .

kill:
	docker kill $$(docker ps -q) &> /dev/null;
deploy: build
	docker tag golang-app stor.highloadcup.ru/travels/little_eagle
	docker push stor.highloadcup.ru/travels/little_eagle

tests:
#	$$GOBIN/highloadcup_tester -addr http://127.0.0.1:8080 -hlcupdocs $$HOME/workspace/hlcupdocs/data/FULL/ -test -phase 1
	$$GOBIN/highloadcup_tester -addr http://127.0.0.1:8080 -hlcupdocs $$HOME/workspace/hlcupdocs/data/FULL/ -test -phase 2
#	$$GOBIN/highloadcup_tester -addr http://127.0.0.1:8080 -hlcupdocs $$HOME/workspace/hlcupdocs/data/FULL/ -test -phase 3

app-run: app-unzip run

app-unzip:
	mkdir -p $$(pwd)/data/ > /dev/null
	unzip -oq /tmp/data/data.zip -d $$(pwd)/data/