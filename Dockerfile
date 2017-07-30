FROM blacktop/elasticsearch:5.5

LABEL maintainer "https://github.com/blacktop"

COPY . /go/src/github.com/blacktop/scifgif
RUN apk --no-cache add ca-certificates
RUN apk --update add --no-cache -t .build-deps \
                                    build-base \
                                    mercurial \
                                    musl-dev \
                                    openssl \
                                    bash \
                                    wget \
                                    git \
                                    gcc \
                                    go \
  && echo "===> Building scifgif Go binary..." \
  && cd /go/src/github.com/blacktop/scifgif \
  && export GOPATH=/go \
  && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags "-X main.Version=$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/scifgif . \
  && rm -rf /go /usr/local/go /usr/lib/go /tmp/* \
  && apk del --purge .build-deps

COPY config/elasticsearch.yml /usr/share/elasticsearch/config/elasticsearch.yml

ENV IMAGE_NUMBER 25

RUN echo "===> Updating images..." \
  && mkdir -p images/xkcd \
  && mkdir -p images/giphy \
  && scifgif update

EXPOSE 9200

ENTRYPOINT ["scifgif"]
