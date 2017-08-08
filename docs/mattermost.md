Use with Mattermost
===================

To run a dev mattermost setup in docker

```sh
$ git clone https://github.com/mattermost/mattermost-docker.git
$ cd mattermost-docker
$ docker-compose up -d
```

Open [http://localhost](http://localhost)

Add an Integration
------------------

![add-integration](https://raw.githubusercontent.com/blacktop/scifgif/master/docs/imgs/add-integration.png)

Add Slash Commands
------------------

`xkcd` slash

![xkcd-slash](https://raw.githubusercontent.com/blacktop/scifgif/master/docs/imgs/xkcd-slash.png)

`scifgif` slash

![scifgif-slash](https://raw.githubusercontent.com/blacktop/scifgif/master/docs/imgs/scifgif-slash.png)

Start **scifgif** microservice with new **token**

```sh
$ docker run --init -d \
             -p 3993:3993 \
             --name scifgif \
             blacktop/scifgif --host HOST --port 3993 --token MATTERMOST_INTEGRATION_TOKEN
```

> **NOTE:** token auth is currently disabled because `scifgif` hosts two microservices `xkcd` and `giphy`

Now in **mattermost** you can type `/xkcd physics` or `/scifgif thumbs up` to activate the slash commands

---
