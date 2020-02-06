backend-test:
	chmod +x ./start_test_env.sh
	./start_test_env.sh
	go test ./...

backend-generate:
	go generate ./...