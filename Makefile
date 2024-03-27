run:
	@docker run -p 8080:8080 alextanhongpin/app:$(shell cat VERSION)


build:
	@DOCKER_BUILDKIT=1 docker build --progress=plain -t alextanhongpin/app:latest -t alextanhongpin/app:$(shell cat VERSION) .


build-distroless:
	@docker build -f Dockerfile.distroless -t alextanhongpin/app:latest -t alextanhongpin/app:$(shell cat VERSION) .


debug:
	@docker exec -it `docker ps -q` ps -ef


clean:
	@docker system prune --volumes --force
