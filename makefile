BINARY_NAME=vcs
INSTALL_PATH=/usr/local/bin

build:
	go build -o $(BINARY_NAME) main.go

install: build
	@echo "Requires to enter sudo password"
	sudo cp $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)

uninstall:
	rm -f $(INSTALL_PATH)/$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)

