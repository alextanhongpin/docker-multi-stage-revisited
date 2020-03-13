start:
	@docker run -p 8080:8080 alextanhongpin/app

docker:
	@docker build -t alextanhongpin/app .
