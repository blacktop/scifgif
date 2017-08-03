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

Add Slash Commands
------------------

`xkcd` slash  

![xkcd-slash](https://raw.githubusercontent.com/blacktop/scifgif/master/docs/imgs/xkcd-slash.png)

`scifgif` slash  

![scifgif-slash](https://raw.githubusercontent.com/blacktop/scifgif/master/docs/imgs/scifgif-slash.png)


Start **scifgif** microservice with new **token**

```sh
$ docker run --init -d --name scifgif -p 3993:3993 blacktop/scifgif --host HOST --port PORT --token sdqg4tm6jiy1zceyt6p7i8i6jr
```

---
