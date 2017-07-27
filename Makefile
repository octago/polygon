server: ulimitimage
	go run main.go ./templates

client:
	go run client/client.go

ulimitimage: ulimitabuser
	docker build -t ulimit images/ulimit

ulimitabuser:
	GOOS=linux go build  -o images/ulimit/ulimitabuser ./ulimitabuser

.PHONY: ulimitimage ulimitabuser server client
