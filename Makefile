APP_NAME = processor

CMD = ./CMD

run:
	echo "Running $(APP_NAME)..."
	go run $(CMD)/main.go