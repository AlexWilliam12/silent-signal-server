APP=silent-signal
SRC=main.go
TARGET=target/bin

.PHONK: build
build:
	@if [ ! -d "$(TARGET)" ]; then \
		mkdir -p "$(TARGET)"; \
	fi
	@if [ ! -d "uploads" ]; then \
		mkdir -p "uploads"; \
	fi
	@echo "Building the app $(APP)..."
	go build -o $(TARGET)/$(APP) $(SRC)
	@echo "Build has finished"

.PHONK: build
run: build
	@./$(TARGET)/$(APP)

.PHONK: clear
clear:
	@echo "Cleaning binaries..."
	@rm $(TARGET)/$(APP)
	@echo "Clear has finished"
