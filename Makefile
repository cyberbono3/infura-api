
install:
	go install

run:
	go run main.go

test:
	go test -v ./...

docker:
	docker build -t cyberbono3/infura:latest .

docker-run:
	docker run -d -p 8080:8080 cyberbono3/infura:latest

