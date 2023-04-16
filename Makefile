.PHONY: all lint lint-ci integration-tests swagger-tool

gobuild:
	@GOOS=linux GOARCH=amd64 go build -o ./build/bin/tx-builder

gobuild-dbg:
	CGO_ENABLED=1 go build -gcflags=all="-N -l" -i -o ./build/bin/tx-builder

run: gobuild
	@./build/bin/tx-builder run

run-dbg: gobuild-dbg
	@dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./build/bin/tx-builderr run

run-server-dev:
	REST_PORT=8001 go run main.go api run
