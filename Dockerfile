##############################################
# API BUILDER                                #
##############################################
FROM golang:alpine as api

COPY . /go/src/github.com/blacktop/scifgif

RUN echo "===> Building scifgif binary..."
RUN apk add build-base

WORKDIR /go/src/github.com/blacktop/scifgif

RUN go build -ldflags "-X main.Version=$(git describe --tags) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/scifgif

##############################################
# WEB BUILDER                                #
##############################################
FROM node:12-alpine as web

COPY . /scifgif

RUN echo "===> Building scifgif Web UI..."

WORKDIR /scifgif/web

RUN yarn
RUN yarn build

##############################################
# SCIFGIF                                    #
##############################################
FROM alpine:3.14

LABEL maintainer "https://github.com/blacktop"

RUN apk --no-cache add ca-certificates

COPY --from=api /bin/scifgif /bin/scifgif

ARG IMAGE_XKCD_COUNT=-1
# ARG IMAGE_DILBERT_DATE=2016-04-28
ARG IMAGE_DILBERT_DATE=2019-01-01
ARG IMAGE_NUMBER

WORKDIR /scifgif

COPY ascii/emoji.json /scifgif/ascii/emoji.json
RUN echo "===> Updating images..." \
  && mkdir -p /scifgif/images/xkcd \
  && mkdir -p /scifgif/images/giphy \
  && mkdir -p /scifgif/images/contrib \
  && mkdir -p /scifgif/images/dilbert \
  && scifgif -V update

COPY images/default/giphy.gif /scifgif/images/default/giphy.gif
COPY images/default/xkcd.png /scifgif/images/default/xkcd.png
COPY images/default/dilbert.png /scifgif/images/default/dilbert.png
COPY images/icons/giphy-icon.png /scifgif/images/icons/giphy-icon.png
COPY images/icons/xkcd-icon.jpg /scifgif/images/icons/xkcd-icon.jpg
COPY images/icons/dilbert-icon.png /scifgif/images/icons/dilbert-icon.png

COPY --from=web /scifgif/web/build /scifgif/web/build

EXPOSE 3993

ENTRYPOINT ["scifgif"]
