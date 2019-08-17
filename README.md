![logo](https://raw.githubusercontent.com/blacktop/scifgif/master/docs/imgs/logo.png)

[![Circle CI](https://circleci.com/gh/blacktop/scifgif.png?style=shield)](https://circleci.com/gh/blacktop/scifgif)
[![License](https://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org)
[![Docker Stars](https://img.shields.io/docker/stars/blacktop/scifgif.svg)](https://store.docker.com/community/images/blacktop/scifgif)
[![Docker Pulls](https://img.shields.io/docker/pulls/blacktop/scifgif.svg)](https://store.docker.com/community/images/blacktop/scifgif)
[![Docker Image](https://img.shields.io/badge/docker%20image-2GB-blue.svg)](https://store.docker.com/community/images/blacktop/scifgif)

> Humorous image microservice for isolated networks - `xkcd|dilbert|giphy` +
> **keyword/phrase** search API

---

### Dependencies

* [alpine:3.10](https://hub.docker.com/_/alpine/)

### Image Tags

```bash
REPOSITORY           TAG                 SIZE
blacktop/scifgif     latest              3GB
blacktop/scifgif     1.0                 3GB
blacktop/scifgif     0.3.0               2GB
blacktop/scifgif     0.2.0               2GB
blacktop/scifgif     0.1.0               2GB
```

> **NOTE:** the reason the docker image is so large is that it contains:
>
> * ~3000 `giphy` gifs _(1500 reactions, 250 futurama and 250 star wars)_
> * all of `xkcd`
> * bunch of `dilbert`

## Getting Started

```bash
$ docker run --init -d --name scifgif -p 3993:3993 blacktop/scifgif --host localhost
```

## Documentation

* [API Docs](http://docs.scifgif.apiary.io)
* [Use with Mattermost](https://github.com/blacktop/scifgif/blob/master/docs/mattermost.md)

![mattermost](https://raw.githubusercontent.com/blacktop/scifgif/master/docs/imgs/mattermost.png)

* [Image Search](https://github.com/blacktop/scifgif/blob/master/docs/image-search.md)

![mattermost](https://raw.githubusercontent.com/blacktop/scifgif/master/docs/imgs/image-search.png)

## TODO

* [x] Add meta-data DB for keyword text search (database)
* [x] Add docs for creating [Mattermost](https://github.com/mattermost/platform)
      slash command or integration
* [ ] Add ability to use expansion packs (use tag to control types of images
      available)
* [x] Add ascii art emojis (table flippers etc)
* [ ] Add ephemeral slash command help
* [x] remove non-alphanumerics so you can search for emojis :older_man:
* [x] remove `xkcd` details _(can they be a mouse over?)_
* [x] add ability to add GIFs and/or keywords to GIFs
* [x] add `dilbert`
* [ ] make search more "fuzzy"

## Issues

Find a bug? Want more features? Find something missing in the documentation? Let
me know! Please don't hesitate to
[file an issue](https://github.com/blacktop/scifgif/issues/new)

## License

MIT Copyright (c) 2017-2018 **blacktop**

![giphy](https://raw.githubusercontent.com/blacktop/scifgif/master/docs/PoweredBy_200_Horizontal_Light-Backgrounds_With_Logo.gif)
