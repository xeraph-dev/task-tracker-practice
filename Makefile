build:
	go build -o build/task main.go

run:
	go run main.go

clean:
	rm -rf build

release: windows darwin linux

windows: windows-arm64 windows-amd64

windows-arm64:
	GOOS=windows GOARCH=arm64 go build -ldflags '-w -s' -trimpath -o build/task-windows-arm64.exe main.go

windows-amd64:
	GOOS=windows GOARCH=amd64 go build -ldflags '-w -s' -trimpath -o build/task-windows-amd64.exe main.go

darwin: darwin-arm64 darwin-amd64

darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -ldflags '-w -s' -trimpath -o build/task-darwin-arm64 main.go

darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -ldflags '-w -s' -trimpath -o build/task-darwin-amd64 main.go

linux: linux-arm64 linux-amd64

linux-arm64:
	GOOS=linux GOARCH=arm64 go build -ldflags '-w -s' -trimpath -o build/task-linux-arm64 main.go

linux-amd64:
	GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -trimpath -o build/task-linux-amd64 main.go
