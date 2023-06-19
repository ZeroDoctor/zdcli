
.PHONY: build
build:
	go build -o zd .

.PHONY: install
install:
	mkdir -p ~/scripts || true
	cp ./zd ~/scripts

.PHONY: release
release:
	tar -cvJf zd-amd64-unix.tar.xz zd lua/
