run:
	go run .

format:
	go fmt ./...

install_latest:
	curl -L "https://github.com/dithmer/bookie/releases/download/latest/bookie" -o "$$HOME/.local/bin/bookie"
	chmod +x "$$HOME/.local/bin/bookie"
