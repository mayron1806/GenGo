build:
	go build -o bin/main main.go

run:
	go build -o bin/gengo && ./bin/gengo $(ARGS)


windows-build:
	env GOOS=windows GOARCH=amd64 go build -o bin/gengo.exe main.go