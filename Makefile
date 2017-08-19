run-local: build
#	if $$(docker ps -q); then
#		docker kill $$(docker ps -q) &> /dev/null;
#	fi
	docker run --rm -p 8080:80 -v ~/Downloads/data5.zip:/tmp/data/data.zip -t golang-app
run:
	app
build:
	docker build -t golang-app .

app-run: app-unzip run

app-unzip:
	mkdir -p $$(pwd)/data/ > /dev/null
	unzip -oq /tmp/data/data.zip -d $$(pwd)/data/