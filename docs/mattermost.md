Use with Mattermost
===================

```sh
$ git clone https://github.com/mattermost/mattermost-docker.git
$ cd mattermost-docker
$ docker-compose up -d
```

Add an Integration
------------------

![add-integration](https://raw.githubusercontent.com/blacktop/scifgif/master/docs/imgs/add-integration.png)

Add an Outgoing Webhook
-----------------------

![outgoing-webhook](https://raw.githubusercontent.com/blacktop/scifgif/master/docs/imgs/outgoing-webhook.png)

### Take note of the new **token**

![webhook-token](https://raw.githubusercontent.com/blacktop/scifgif/master/docs/imgs/webhook-token.png)

Start **scifgif** microservice with new **token**

```sh
$ docker run --init -d --name scifgif blacktop/scifgif --host localhost --token sdqg4tm6jiy1zceyt6p7i8i6jr
```

---
