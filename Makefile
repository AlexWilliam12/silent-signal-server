APP=silent-signal
SRC=cmd/app/main.go
TARGET=target/bin

.PHONK: build
build:
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