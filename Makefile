test-coverage:
	go test -v ./... -coverprofile cover.out
	go tool cover -html=cover.out

test:
	go test -v ./...