# React frontend build configuration
.PHONY: all install build clean serve cache-clear cache-verify reinstall

all: build

# Install dependencies
install:
	@echo "Installing dependencies..."
	npm ci

# Build React application
build: install
	@echo "Building React application..."
	npm run build

# Clean build artifacts
clean:
	@echo "Cleaning build directory..."
	rm -rf build
	@echo "Clean complete!"

# Clean dependencies and cache
clean-all: clean
	@echo "Cleaning dependencies and cache..."
	npm run clean
	@echo "Clean complete!"

# Verify npm cache
cache-verify:
	@echo "Verifying npm cache..."
	npm run cache:verify

# Clear npm cache
cache-clear:
	@echo "Clearing npm cache..."
	npm run cache:clear

# Reinstall dependencies
reinstall:
	@echo "Reinstalling dependencies..."
	npm run reinstall

# Serve the application for local development
serve: install
	@echo "Starting development server at http://localhost:3000"
	npm start
