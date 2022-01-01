image_name := myapp

.PHONY: build
build:
	@docker build --rm -t $(image_name) -f docker/Dockerfile . 
	@docker image prune -f

.PHONY: run
run: build
	@docker run -it --rm -v /var/run/docker.sock:/var/run/docker.sock $(image_name)

.PHONY: rund
rund: build
	@docker run -d --rm -v /var/run/docker.sock:/var/run/docker.sock $(image_name)