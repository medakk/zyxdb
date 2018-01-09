compose-build:
	docker-compose build

compose-up: compose-build
	docker-compose up
