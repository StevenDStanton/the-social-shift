# Makefile for building the game for Web, PC (Windows & Linux), Android, or All

# Variables
PROJECT_PATH = github.com/StevenDStanton/game-name
DIST_DIR = dist
VERSION = 0.0.10
GOROOT_WASM = $(shell go env GOROOT)/misc/wasm/wasm_exec.js
BUTLER_USER = ApocalypseTheory
BUTLER_PROJECT = game-name
BUTLER_PLATFORM = HTML5

# Default build target
.PHONY: build build-web build-pc-windows build-pc-linux build-android build-all clean setup-web setup-pc-windows setup-pc-linux setup-android compile-web compile-pc-windows compile-pc-linux compile-android create-html deploy-web dist-pc-windows dist-pc-linux dist-android cleanup-web cleanup-pc-windows cleanup-pc-linux cleanup-android depend status demo test help

# Help Target: Display available targets
help:
	@echo "Available targets:"
	@echo "  make build TYPE=<web|pc-windows|pc-linux|android|all> - Build the game for specified platform(s)"
	@echo "  make build-web                                    - Build and deploy for Web (HTML5)"
	@echo "  make build-pc-windows                             - Build for PC (Windows)"
	@echo "  make build-pc-linux                               - Build for PC (Linux)"
	@echo "  make build-android                                - Build for Android"
	@echo "  make build-all                                    - Build for all platforms"
	@echo "  make clean                                        - Clean all build artifacts"
	@echo "  make depend                                       - Update dependencies"
	@echo "  make demo                                         - Run the Linux build locally"
	@echo "  make test                                         - Run tests"
	@echo "  make status                                       - Check deployment status on itch.io"

# Main build target with TYPE parameter
build:
	@$(if $(TYPE),,$(error Please specify the build TYPE. Example: make build TYPE=web))
	@$(if $(filter web,$(TYPE)), make build-web)
	@$(if $(filter pc-windows,$(TYPE)), make build-pc-windows)
	@$(if $(filter pc-linux,$(TYPE)), make build-pc-linux)
	@$(if $(filter android,$(TYPE)), make build-android)
	@$(if $(filter all,$(TYPE)), make build-all)
	@$(if $(filter-out web pc-windows pc-linux android all,$(TYPE)), \
		$(error Unknown TYPE "$(TYPE)". Available TYPEs: web, pc-windows, pc-linux, android, all))

# Build all platforms
build-all: build-web build-pc-windows build-pc-linux build-android

# Build Web
build-web: depend setup-web compile-web create-html deploy-web cleanup-web

# Setup for Web
setup-web:
	@echo "=== Setting Up Web Build ==="
	@echo "Creating build/web folder."
	@mkdir -p build/web
	@echo "Copying wasm_exec.js to build/web folder."
	@cp $(GOROOT_WASM) build/web/
	@echo "Copying assets to build/web folder."
	@cp -r assets build/web/
	@echo "Web setup completed."

# Compile Web (WASM)
compile-web:
	@echo "=== Compiling Game to WebAssembly ==="
	@env GOOS=js GOARCH=wasm go build -o build/web/game.wasm $(PROJECT_PATH)
	@echo "WebAssembly compilation completed."

# Create HTML for Web
create-html:
	@echo "=== Creating HTML File for Web ==="
	@echo '<!DOCTYPE html><html><head><meta charset="UTF-8"><title>Alchemist of the Shadow Bureau</title></head><body><script src="wasm_exec.js"></script><script>if(!WebAssembly.instantiateStreaming){WebAssembly.instantiateStreaming=async(resp,importObject)=>{const source=await(resp).arrayBuffer();return await WebAssembly.instantiate(source,importObject);};}const go=new Go();WebAssembly.instantiateStreaming(fetch("game.wasm"),go.importObject).then(result=>{go.run(result.instance);});</script></body></html>' > build/web/index.html
	@echo "HTML file created."

# Deploy Web (upload to itch.io)
deploy-web:
	@echo "=== Deploying Web Build to itch.io ==="
	@butler push build/web $(BUTLER_USER)/$(BUTLER_PROJECT):$(BUTLER_PLATFORM) --userversion $(VERSION)
	@echo "Web build deployed successfully."

# Clean up Web Build
cleanup-web:
	@echo "=== Cleaning Up Web Build ==="
	@rm -rf build/web
	@echo "Web build cleanup completed."

# Build PC Windows
build-pc-windows: depend setup-pc-windows compile-pc-windows dist-pc-windows cleanup-pc-windows

# Setup for PC Windows
setup-pc-windows:
	@echo "=== Setting Up PC (Windows) Build ==="
	@echo "Creating dist/pc/windows folder."
	@mkdir -p $(DIST_DIR)/pc/windows
	@echo "PC (Windows) setup completed."

# Compile PC Windows
compile-pc-windows:
	@echo "=== Compiling Game for PC (Windows) ==="
	@env GOOS=windows GOARCH=amd64 go build -o $(DIST_DIR)/pc/windows/game-name.exe $(PROJECT_PATH)
	@echo "PC (Windows) compilation completed."

# Create Distribution for PC Windows
dist-pc-windows:
	@echo "=== PC (Windows) Build is Located at $(DIST_DIR)/pc/windows/game-name.exe ==="

# Clean up PC Windows Build
cleanup-pc-windows:
	@echo "=== Cleaning Up PC (Windows) Build ==="
	@rm -rf build/pc_windows_temp # Adjust if temporary files are created
	@echo "PC (Windows) build cleanup completed."

# Build PC Linux
build-pc-linux: depend setup-pc-linux compile-pc-linux dist-pc-linux cleanup-pc-linux

# Setup for PC Linux
setup-pc-linux:
	@echo "=== Setting Up PC (Linux) Build ==="
	@echo "Creating dist/pc/linux folder."
	@mkdir -p $(DIST_DIR)/pc/linux
	@echo "PC (Linux) setup completed."

# Compile PC Linux
compile-pc-linux:
	@echo "=== Compiling Game for PC (Linux) ==="
	@env GOOS=linux GOARCH=amd64 go build -o $(DIST_DIR)/pc/linux/game-name $(PROJECT_PATH)
	@echo "PC (Linux) compilation completed."

# Create Distribution for PC Linux
dist-pc-linux:
	@echo "=== PC (Linux) Build is Located at $(DIST_DIR)/pc/linux/game-name ==="

# Clean up PC Linux Build
cleanup-pc-linux:
	@echo "=== Cleaning Up PC (Linux) Build ==="
	@rm -rf build/pc_linux_temp # Adjust if temporary files are created
	@echo "PC (Linux) build cleanup completed."

# Build Android
build-android: depend setup-android compile-android dist-android cleanup-android

# Setup for Android
setup-android:
	@echo "=== Setting Up Android Build ==="
	@echo "Creating dist/android folder."
	@mkdir -p $(DIST_DIR)/android
	@echo "Android setup completed."

# Compile Android (APK using gomobile)
compile-android:
	@echo "=== Compiling Game for Android ==="
	@gomobile init
	@gomobile build -target=android -o $(DIST_DIR)/android/game-name.apk $(PROJECT_PATH)
	@echo "Android compilation completed."

# Create Distribution for Android
dist-android:
	@echo "=== Android Build is Located at $(DIST_DIR)/android/game-name.apk ==="

# Clean up Android Build
cleanup-android:
	@echo "=== Cleaning Up Android Build ==="
	@rm -rf build/android_temp # Adjust if temporary files are created
	@echo "Android build cleanup completed."

# Clean all build artifacts
clean:
	@echo "=== Cleaning All Build Artifacts ==="
	@rm -rf build
	@rm -rf $(DIST_DIR)
	@echo "All build artifacts cleaned."

# Update dependencies
depend:
	@echo "=== Updating Dependencies ==="
	@go mod tidy
	@echo "Dependencies updated."

# Check deployment status on itch.io
status:
	@echo "=== Checking Deployment Status on itch.io ==="
	@butler status $(BUTLER_USER)/$(BUTLER_PROJECT):$(BUTLER_PLATFORM)

# Run the Linux build locally
demo:
	@echo "=== Running the Game Locally (Linux Build) ==="
	@./$(DIST_DIR)/pc/linux/game-name

# Run tests
test:
	@echo "=== Running Tests ==="
	@go test ./...
