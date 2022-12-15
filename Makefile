run:
	@docker run -p 8080:8080 alextanhongpin/app:$(shell cat VERSION)

build:
	@docker build -t alextanhongpin/app:latest -t alextanhongpin/app:$(shell cat VERSION) .


build-distroless:
	@docker build -f Dockerfile.distroless -t alextanhongpin/app:latest -t alextanhongpin/app:$(shell cat VERSION) .
