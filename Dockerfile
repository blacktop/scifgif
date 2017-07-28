FROM golang:1.8.3 as builder

COPY . /go/src/github.com/blacktop/scifgif
WORKDIR /go/src/github.com/blacktop/scifgif

# RUN go get -u github.com/golang/dep/cmd/dep
# RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
  -ldflags "-X main.Version=$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o app .

FROM blacktop/elasticsearch:5.5
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /go/src/github.com/blacktop/scifgif/app .

RUN mkdir -p images/xkcd \
  && mkdir -p images/giphy \
  && ./app update

# COPY images /root/images

CMD ["./app"]
