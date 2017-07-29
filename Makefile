.PHONY: build dev size tags tar test run ssh circle push

ORG=blacktop
NAME=scifgif
REPO=$(ORG)/$(NAME)
VERSION?=$(shell cat VERSION)

build: ## Build docker image
	docker build -t $(ORG)/$(NAME):$(VERSION) .

dev: base ## Build docker dev image
	docker build -f Dockerfile.dev -t $(ORG)/$(NAME):$(VERSION) .

size: tags ## Update docker image size in README.md
	sed -i.bu 's/docker%20image-.*-blue/docker%20image-$(shell docker images --format "{{.Size}}" $(ORG)/$(NAME):$(VERSION)| cut -d' ' -f1)-blue/' README.md
	sed -i.bu '/latest/ s/[0-9.]\{3,5\}MB/$(shell docker images --format "{{.Size}}" $(ORG)/$(NAME):$(VERSION))/' README.md
	sed -i.bu '/$(VERSION)/ s/[0-9.]\{3,5\}MB/$(shell docker images --format "{{.Size}}" $(ORG)/$(NAME):$(VERSION))/' README.md

tags: ## Show all docker image tags
	docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}" $(ORG)/$(NAME)

run: stop ## Run scifgif web-service
	@echo "===> Starting scifgif..."
	@docker run --init -d --name scifgif -p 3993:3993 $(ORG)/$(NAME):$(VERSION)

ssh: ## SSH into docker image
	@docker run --init -it --rm --entrypoint=sh $(ORG)/$(NAME):$(VERSION)

tar: ## Export tar of docker image
	@docker save $(ORG)/$(NAME):$(VERSION) -o $(NAME).tar

test: ## Test build plugin
	@echo "===> Starting scifgif tests..."
	@docker run --init --rm -p 3993:3993 $(ORG)/$(NAME):$(VERSION)
	@http 127.0.0.1:3993/xkcd/city.jpg > city.jpg
	@ls -lah city.jpg
	@rm city.jpg

gotest: stop ## Test go code
	@echo "===> Starting gotests..."
	@echo " - Starting elasticsearch"
	@docker run -d --name elasticsearch -p 9200:9200 blacktop/elasticsearch:5.5
	@echo " - Starting kibana"
	@sleep 10;docker run -d --name kibana --link elasticsearch -p 5601:5601 blacktop/kibana:5.5
	@go run scifgif.go update

push: build ## Push docker image to docker registry
	@echo "===> Pushing $(ORG)/$(NAME):$(VERSION) to docker hub..."
	@docker push $(ORG)/$(NAME):$(VERSION)

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

stop: ## Kill running scifgif-plugin docker containers
	@docker rm -f scifgif || true
	@docker rm -f elasticsearch || true
	@docker rm -f kibana || true

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
