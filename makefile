run:
	go run server.go


build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o spurtcms-graphql
