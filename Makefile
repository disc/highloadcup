run-local: build
	docker kill $$(docker ps -q); docker run --rm -p 8080:80 -t golang-app
run:
	./app
build:
	docker build -t golang-app .