run: build-assets
	@echo "Starting Go server..."
	@go run ./cmd/server/main.go

build-assets:
	@echo "Building frontend assets..."
	@reset
	@./build-assets.sh

clean:
	@echo "Cleaning up generated files..."
	@rm -f web/static/css/style.css
	@rm -f web/static/js/bundle.js

