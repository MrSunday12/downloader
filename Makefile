image_name := myapp

build:
	@docker build --rm -t $(image_name) -f docker/Dockerfile . 
	@docker image prune -f

run: build
	@docker run -it --rm -v /var/run/docker.sock:/var/run/docker.sock $(image_name)

rund: build
	@docker run -d $(image_name)