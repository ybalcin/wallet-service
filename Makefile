.PHONY: run
run:
	MONGO_URI="mongodb+srv://wallet-service-user:EmGcfRCHL9tAoyYZ@cluster0.l1pmb.mongodb.net/?retryWrites=true&w=majority" DATABASE_NAME="wallet-service-db" PORT="8080" go run main.go

.PHONY: run-test
run-test:
	go test ./...