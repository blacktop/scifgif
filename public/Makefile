SCIFGIF?=0.3.0

.PHONY: dev
dev:
	npm install

.PHONY: start
start: stop dev
	@echo "===> Starting scifgif..."
	@docker run --init -d --name scifgif -p 3993:3993 -p 9200:9200 blacktop/scifgif:$(SCIFGIF) --host localhost;sleep 15
	@open -a Brave http://localhost:8080/webpack-dev-server/
	@npm run start

.PHONY: stop
stop:
	@echo "===> Stopping scifgif..."
	@docker rm -f scifgif || true
