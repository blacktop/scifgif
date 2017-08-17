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
  && echo "===> Building scifgif binary..." \
  && cd /go/src/github.com/blacktop/scifgif \
  && export GOPATH=/go \
  && go build -ldflags "-X main.Version=$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/scifgif \
  && rm -rf /go /usr/local/go /usr/lib/go /tmp/* \
  && apk del --purge .build-deps

COPY config/elasticsearch.yml /usr/share/elasticsearch/config/elasticsearch.yml

ARG IMAGE_XKCD_COUNT=-1
ARG IMAGE_NUMBER

WORKDIR /scifgif

RUN echo "===> Create elasticsearch data directory..." \
  && mkdir -p /scifgif/elasticsearch/data \
  && chown -R elasticsearch:elasticsearch /scifgif/elasticsearch/data

RUN echo "===> Updating images..." \
  && mkdir -p /scifgif/images/xkcd \
  && mkdir -p /scifgif/images/giphy \
  && scifgif update \
  && echo "===> Stopping elasticsearch PID: $(cat /tmp/epid)" \
  && sleep 10; kill $(cat /tmp/epid) \
  && wait $(cat /tmp/epid); exit 0;

COPY images/default/giphy.gif /scifgif/images/default/giphy.gif
COPY images/default/xkcd.png /scifgif/images/default/xkcd.png
COPY images/icons/giphy-icon.png /scifgif/images/icons/giphy-icon.png
COPY images/icons/xkcd-icon.jpg /scifgif/images/icons/xkcd-icon.jpg

EXPOSE 3993

ENTRYPOINT ["scifgif"]
