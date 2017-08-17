.PHONY: build dev size tags tar test run mattermost update gotest ssh circle push

ORG=blacktop
NAME=scifgif
REPO=$(ORG)/$(NAME)
VERSION?=$(shell cat VERSION)

build: ## Build docker image
	docker build --build-arg IMAGE_XKCD_COUNT=100 --build-arg IMAGE_NUMBER=100 -t $(ORG)/$(NAME):$(VERSION) .

dev: base ## Build docker dev image
	docker build -f Dockerfile.dev -t $(ORG)/$(NAME):$(VERSION) .

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
	@echo " - Starting elasticsearch"
	@docker run -d --name elasticsearch -p 9200:9200 blacktop/elasticsearch:5.5
	@echo " - Starting kibana"
	@sleep 10;docker run -d --name kibana --link elasticsearch -p 5601:5601 blacktop/kibana:5.5
	@go run scifgif.go -N 50 -V update

web: stop ## Start scifgif web-service
	@echo "===> Starting scifgif web service..."
	@go run scifgif.go --host 127.0.0.1

run: stop ## Run scifgif
	docker run -d --name scifgif -p 3993:3993 -p 9200:9200 $(ORG)/$(NAME):$(VERSION)

mattermost: ## Start mattermost
	git clone https://github.com/mattermost/mattermost-docker.git || true
	cp docker-compose.mattermost.yml mattermost-docker/docker-compose.yml
	cd mattermost-docker;docker-compose up -d --build

test: ## Test build plugin
	@echo "===> Starting scifgif tests..."
	@docker run --init --rm -p 3993:3993 $(ORG)/$(NAME):$(VERSION)
	@http 127.0.0.1:3993/xkcd/city.jpg > city.jpg
	@ls -lah city.jpg
	@rm city.jpg

circle: ci-size ## Get docker image size from CircleCI
	@sed -i.bu 's/docker%20image-.*-blue/docker%20image-$(shell cat .circleci/SIZE)-blue/' README.md
	@echo "===> Image size is: $(shell cat .circleci/SIZE)"

ci-build:
	@echo "===> Getting CircleCI build number"
	@http https://circleci.com/api/v1.1/project/github/${REPO} | jq '.[0].build_num' > .circleci/build_num

ci-size: ci-build
	@echo "===> Getting image build size from CircleCI"
	@http "$(shell http https://circleci.com/api/v1.1/project/github/${REPO}/$(shell cat .circleci/build_num)/artifacts circle-token==${CIRCLE_TOKEN} | jq '.[].url')" > .circleci/SIZE

clean: ## Clean docker image and stop all running containers
	docker-clean stop
	docker rmi $(ORG)/$(NAME):$(VERSION) || true
	rm images/giphy/*.gif || true
	rm images/xkcd/*.png || true
	rm -rf mattermost-docker

stop: ## Kill running scifgif-plugin docker containers
	@docker rm -f scifgif || true

dbstop: ## Kill DB containers
		@docker rm -f elasticsearch || true
		@docker rm -f kibana || true

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
