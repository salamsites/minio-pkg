# Run tests
sminio_run:
	go run test/server/server.go

sminio_run_and_test:
	go run test/server/server.go & sleep 5
	go test -v -count=1 ./...

sminio_test:
	go test -v -count=1 ./...
