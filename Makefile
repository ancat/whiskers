main:
	go build .

fmt-branch:
	gofmt -w $$(git diff --name-only main '*.go')

release:
	GOOS=linux GOARCH=amd64 go build -o dist/whiskers.linux.amd64
	GOOS=linux GOARCH=arm64 go build -o dist/whiskers.linux.arm64
	GOOS=darwin GOARCH=amd64 go build -o dist/whiskers.darwin.amd64
	GOOS=darwin GOARCH=arm64 go build -o dist/whiskers.darwin.arm64
