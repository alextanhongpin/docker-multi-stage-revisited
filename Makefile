VERSION := $(shell cat VERSION)
export

run:
	@docker run -p 8080:8080 alextanhongpin/app:$(VERSION)


build:
	@DOCKER_BUILDKIT=1 docker build --progress=plain -t alextanhongpin/app:latest -t alextanhongpin/app:$(VERSION) .


build-distroless:
	@docker build -f Dockerfile.distroless -t alextanhongpin/app:latest -t alextanhongpin/app:$(VERSION) .


debug:
	@docker exec -it `docker ps -q` ps -ef


clean:
	@docker system prune --volumes --force


up:
	@docker-compose up -d

down:
	@docker-compose down

load:
	#apr_socket_connect(): Invalid argument (22)
	# Use 127.0.0.1 instead of localhost.
	ab -n 1000 -c 10 http://127.0.0.1:8080/
