.PHONY: build dev size tags tar test run mattermost update gotest ssh circle push

ORG=blacktop
NAME=scifgif
REPO=$(ORG)/$(NAME)
VERSION?=$(shell cat VERSION)

build: ## Build docker image
	docker build --build-arg IMAGE_XKCD_COUNT=100 --build-arg IMAGE_NUMBER=100 -t $(ORG)/$(NAME):$(VERSION) .

size: tags ## Update docker image size in README.md
	sed -i.bu 's/docker%20image-.*-blue/docker%20image-$(shell docker images --format "{{.Size}}" $(ORG)/$(NAME):$(VERSION)| cut -d' ' -f1)-blue/' README.md
	sed -i.bu '/latest/ s/[0-9.]\{3,5\}MB/$(shell docker images --format "{{.Size}}" $(ORG)/$(NAME):$(VERSION))/' README.md
	sed -i.bu '/$(VERSION)/ s/[0-9.]\{3,5\}MB/$(shell docker images --format "{{.Size}}" $(ORG)/$(NAME):$(VERSION))/' README.md

tags: ## Show all docker image tags
	docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}" $(ORG)/$(NAME)

tar: ## Export tar of docker image
	@docker save $(ORG)/$(NAME):$(VERSION) -o $(NAME).tar

ssh: ## SSH into docker image
	@docker run -it --rm -p 3993:3993 -p 9200:9200 --entrypoint=sh $(ORG)/$(NAME):$(VERSION)

push: build ## Push docker image to docker registry
	@echo "===> Pushing $(ORG)/$(NAME):$(VERSION) to docker hub..."
	@docker push $(ORG)/$(NAME):$(VERSION)

update: stop dbstop ## Update scifgif images
	@echo "===> Starting scifgif update..."
	@go run *.go -N 20 --xkcd-count 20 --date 2017-05-08 -V update
	@echo "===> Updating web deps..."
	@cd web; npm install

web: stop ## Start scifgif web-service
	@echo "===> Rebuilding web assets..."
	@cd web; npm run build
	@echo "===> Starting scifgif web service..."
	@open http://localhost:3993
	@go run *.go -V --host 127.0.0.1

.PHONY: export
export: stop ## Export scifgif DB
	docker run -d --name scifgif $(ORG)/$(NAME):$(VERSION) -V export; sleep 15
	docker cp scifgif:/mount/backups/snapshot ./database/snapshots/

run: stop ## Run scifgif
	@docker run -d --name scifgif -p 3993:3993 $(ORG)/$(NAME):$(VERSION) --host 127.0.0.1
	# @open http://localhost:8080/webpack-dev-server/
	# @cd public; npm run start

mattermost: ## Start mattermost
	git clone https://github.com/mattermost/mattermost-docker.git || true
	cd mattermost-docker;git checkout 5.7.0
	cp docker-compose.mattermost.yml mattermost-docker/docker-compose.yml
	cp config/mattermost/config.json mattermost-docker/config.json
	cd mattermost-docker;docker-compose up -d --build

test: ## Test build plugin
	@echo "===> Starting scifgif tests..."
	@docker run --init --rm -p 3993:3993 $(ORG)/$(NAME):$(VERSION)
	@http 127.0.0.1:3993/xkcd/city.jpg > city.jpg
	@ls -lah city.jpg
	@rm city.jpg

clean: ## Clean docker image and stop all running containers
	# docker-clean stop
	# docker rmi $(ORG)/$(NAME):$(VERSION) || true
	rm images/giphy/*.gif || true
	rm images/xkcd/*.jpg || true
	rm images/dilbert/*.jpg || true
	rm scifgif.db || true
	rm -r scifgif.bleve || true
	rm -rf mattermost-docker

stop: ## Kill running scifgif-plugin docker containers
	@docker rm -f scifgif || true

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
