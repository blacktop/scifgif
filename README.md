![logo](https://raw.githubusercontent.com/blacktop/scifgif/master/docs/logo.png)

scifgif
=======

[![Circle CI](https://circleci.com/gh/blacktop/scifgif.png?style=shield)](https://circleci.com/gh/blacktop/scifgif) [![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org) [![Docker Stars](https://img.shields.io/docker/stars/blacktop/scifgif.svg)](https://store.docker.com/community/images/blacktop/scifgif) [![Docker Pulls](https://img.shields.io/docker/pulls/blacktop/scifgif.svg)](https://store.docker.com/community/images/blacktop/scifgif) [![Docker Image](https://img.shields.io/badge/docker%20image-1.2GB-blue.svg)](https://store.docker.com/community/images/blacktop/scifgif)

> Humorous image microservice for isolated networks - xkcd and giphy full text search API

---

### Dependencies

-	[alpine:3.6](https://hub.docker.com/_/alpine/)

### Image Tags

```bash
REPOSITORY           TAG                 SIZE
blacktop/scifgif     latest              1.17GB
blacktop/scifgif     0.1.0               1.17GB
```

Getting Started
---------------

```bash
$ docker run --init -d --name scifgif -p 3993:3993 blacktop/scifgif --host localhost
```

Documentation
-------------

-	[Use with Mattermost](https://github.com/blacktop/scifgif/blob/master/docs/mattermost.md)

![mattermost](https://raw.githubusercontent.com/blacktop/scifgif/master/docs/imgs/mattermost.png)

### Web Routes

```apib
FORMAT: 1A

# SCIFgif  

## xkcd [/xkcd]                           
### Get Random Xkcd Comic [GET]

## xkcd [/xkcd/number/{number}]                           
### Get Xkcd by Number [GET]  

+ Parameters

    + number: `1319` (number, required) - The xkcd comic ID

## xkcd [/xkcd/search]                           
### Get Xkcd Search [GET]  

+ ```query```: Query (string, required) - Search terms

## xkcd [/xkcd/slash]                           
### Send xkcd slash query [POST]

+ ```text```: Query (string, required) - Search terms

## xkcd [/xkcd/new_post]                           
### Send xkcd outgoing-webhook query [POST]

+ ```token```: Token (string, required) - Integration token
+ ```trigger_word```: Query (string, required) - Search

## Giphy [/giphy]
### Get Random Giphy Gif [GET]

## Giphy [/giphy/search]                           
### Get Giphy Search [GET]  

+ ```query```: Query (string, required) - Search terms

## Giphy [/giphy/slash]                           
### Send Giphy slash query [POST]

+ ```text```: Query (string, required) - Search terms

## Giphy [/giphy/new_post]                           
### Send Giphy outgoing-webhook query [POST]

+ ```token```: Token (string, required) - Integration token
+ ```trigger_word```: Query (string, required) - Search
```

### TODO

-	[x] Add meta-data DB for keyword text search (elasticsearch)
-	[x] Add docs for creating [Mattermost](https://github.com/mattermost/platform) slash command or integration
- [ ] Add ability to use expansion packs (use tag to control types of images available)
- [ ] Add ascii art emojis (table flippers etc)
- [ ] Add ephemeral slash command help

### Issues

Find a bug? Want more features? Find something missing in the documentation? Let me know! Please don't hesitate to [file an issue](https://github.com/blacktop/scifgif/issues/new)

### License

MIT Copyright (c) 2017 **blacktop**

![giphy](https://raw.githubusercontent.com/blacktop/scifgif/master/docs/PoweredBy_200_Horizontal_Light-Backgrounds_With_Logo.gif)
