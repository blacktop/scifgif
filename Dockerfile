##############################################
# BUILDER                                    # 
##############################################
FROM golang:alpine as builder

COPY . /go/src/github.com/blacktop/scifgif

RUN echo "===> Building scifgif binary..."
WORKDIR /go/src/github.com/blacktop/scifgif
RUN go build -ldflags "-X main.Version=$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/scifgif

##############################################
# SCIFGIF                                    # 
##############################################
FROM blacktop/elasticsearch:5.6

LABEL maintainer "https://github.com/blacktop"

RUN apk --no-cache add ca-certificates

COPY --from=builder /bin/scifgif /bin/scifgif

COPY config/elasticsearch.yml /usr/share/elasticsearch/config/elasticsearch.yml

ARG IMAGE_XKCD_COUNT=-1
ARG IMAGE_DILBERT_DATE=2016-04-28
ARG IMAGE_NUMBER

WORKDIR /scifgif

RUN echo "===> Create elasticsearch data/repo directories..." \
  && mkdir -p /scifgif/elasticsearch/data /mount/backups \
  && chown -R elasticsearch:elasticsearch /scifgif/elasticsearch/data /mount/backups

COPY ascii/emoji.json /scifgif/ascii/emoji.json
RUN echo "===> Updating images..." \
  && mkdir -p /scifgif/images/xkcd \
  && mkdir -p /scifgif/images/giphy \
  && mkdir -p /scifgif/images/contrib \
  && mkdir -p /scifgif/images/dilbert \
  && scifgif update \
  && echo "===> Stopping elasticsearch PID: $(cat /tmp/epid)" \
  && sleep 10; kill $(cat /tmp/epid) \
  && wait $(cat /tmp/epid); exit 0;

COPY images/default/giphy.gif /scifgif/images/default/giphy.gif
COPY images/default/xkcd.png /scifgif/images/default/xkcd.png
COPY images/default/dilbert.png /scifgif/images/default/dilbert.png
COPY images/icons/giphy-icon.png /scifgif/images/icons/giphy-icon.png
COPY images/icons/xkcd-icon.jpg /scifgif/images/icons/xkcd-icon.jpg
COPY images/icons/dilbert-icon.png /scifgif/images/icons/dilbert-icon.png

# Add web app resources
COPY public/index.html /scifgif/public/index.html
COPY public/bundle.js /scifgif/public/bundle.js
COPY public/style/ /scifgif/public/style/
COPY public/etc/passwd /public/public/etc/passwd

EXPOSE 3993

ENTRYPOINT ["scifgif"]
