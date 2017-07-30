Use with Mattermost
===================

```sh
$ git clone https://github.com/mattermost/mattermost-docker.git
$ cd mattermost-docker
$ docker-compose up -d
```

## Add an Integration

![add-integration](https://raw.githubusercontent.com/blacktop/scifgif/master/docs/imgs/add-integration.png)

## Add an Outgoing Webhook  

![outgoing-integration](https://raw.githubusercontent.com/blacktop/scifgif/master/docs/imgs/outgoing-integration.png)

### Take note of the new **token**  

![outgoing-token](https://raw.githubusercontent.com/blacktop/scifgif/master/docs/imgs/outgoing-token.png)

Start **scifgif** microservice with new **token**

```sh
$ docker run --init -d --name scifgif blacktop/scifgif --host HOST --token sdqg4tm6jiy1zceyt6p7i8i6jr
```

---

## Add a Slash Command  

### `xkcd` slash command  

![xkcd-slash](https://raw.githubusercontent.com/blacktop/scifgif/master/docs/imgs/xkcd-slash.png)

### `scifgif` slash command  

![giphy-slash](https://raw.githubusercontent.com/blacktop/scifgif/master/docs/imgs/giphy-slash.png)

Registered Slash Commands

![show-slashes](https://raw.githubusercontent.com/blacktop/scifgif/master/docs/imgs/show-slashes.png)
