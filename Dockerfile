FROM golang:1.8

RUN apt-get update && apt-get install -y --no-install-recommends \
		unzip \
	&& rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/github.com/disc/highloadcup
COPY . .
RUN rm -rf ./data/

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

EXPOSE 80

CMD ["make","app-run"]