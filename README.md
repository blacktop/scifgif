![logo](https://raw.githubusercontent.com/blacktop/scifgif/master/docs/logo.png)

scifgif
=======

[![Circle CI](https://circleci.com/gh/blacktop/scifgif.png?style=shield)](https://circleci.com/gh/blacktop/scifgif) [![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org) [![Docker Stars](https://img.shields.io/docker/stars/blacktop/scifgif.svg)](https://store.docker.com/community/images/blacktop/scifgif) [![Docker Pulls](https://img.shields.io/docker/pulls/blacktop/scifgif.svg)](https://store.docker.com/community/images/blacktop/scifgif) [![Docker Image](https://img.shields.io/badge/docker%20image-1.2GB-blue.svg)](https://store.docker.com/community/images/blacktop/scifgif)

> Humorous Image microservice for isolated networks.

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
$ docker run --init -d --name scifgif -p 3993:3993 blacktop/scifgif
```

### Web Routes

```apib
GET /xkcd (random xkcd comic)
GET /xkcd/{number} (xkcd comic by number)
GET /giphy (random Giphy gif)
GET /giphy/{search} (Giphy gif matching search term)
```

### TODO

-	[x] Add meta-data DB for keyword text search (elasticsearch)
-	[ ] Add Giphy (use tag to control types of images available)
-	[ ] Add docs for creating [Mattermost](https://github.com/mattermost/platform) slash command or integration

### Issues

Find a bug? Want more features? Find something missing in the documentation? Let me know! Please don't hesitate to [file an issue](https://github.com/blacktop/scifgif/issues/new)

### License

MIT Copyright (c) 2017 **blacktop**

![giphy](https://raw.githubusercontent.com/blacktop/scifgif/master/docs/PoweredBy_200_Horizontal_Light-Backgrounds_With_Logo.gif)
