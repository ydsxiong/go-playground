TAG?=$(shell echo "prod-$(shell git rev-list HEAD --max-count=1 --abbrev-commit)")
export TAG

test:
	go test ./...

build:
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(TAG)" -o boxoffice .


pack: build
	docker build -t xman2019/boxoffice-service:$(TAG) .

upload:
	docker push xman2019/boxoffice-service:$(TAG)

deploy:
	envsubst < deployment.yml | kubectl apply -f -

ship: test pack upload deploy
