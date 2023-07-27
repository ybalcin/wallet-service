.PHONY: run
run:
	MONGO_URI="mongodb+srv://wallet-service-user:EmGcfRCHL9tAoyYZ@cluster0.l1pmb.mongodb.net/?retryWrites=true&w=majority" DATABASE_NAME="wallet-service-db" go run main.go