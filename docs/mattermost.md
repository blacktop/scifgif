Use with Mattermost
===================

```sh
$ git clone https://github.com/mattermost/mattermost-docker.git
$ cd mattermost-docker
$ docker-compose up -d
```

## Add an Integration

![add-integration](docs/imgs/add-integration.png)

## Add an Outgoing Webhook  

![outgoing-integration](docs/imgs/outgoing-integration.png)

### Take note of the new **token**  

![outgoing-token](docs/imgs/outgoing-token.png)

```sh
$ docker run --init -d --name scifgif blacktop/scifgif --token
```

---

## Add a Slash Command  

### xkcd slash command  

![xkcd-slash](docs/imgs/xkcd-slash.png)

### giphy slash command  

![giphy-slash](docs/imgs/giphy-slash.png)

Registered Slash Commands

![show-slashes](docs/imgs/show-slashes.png)
