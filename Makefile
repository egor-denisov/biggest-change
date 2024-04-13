.PHONY:

build:
	docker-compose up --build --force-recreate

up:
	docker-compose up -d

run:
	go run cmd/app/main.go

get:
	go get -d -v ./...

test:
	go test -cover ./...   

swag:
	swag init -dir internal/controller/http/v1/ -generalInfo router.go --parseDependency internal/entity/ 
