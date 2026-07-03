
VERSION := v1.2.1-$(shell git rev-parse --short HEAD 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-X main.VERSION=$(VERSION)"

.PHONY: build
build:
	go build $(LDFLAGS) -o zd .

.PHONY: install
install: build
	mkdir -p ~/scripts || true
	cp ./zd ~/scripts
	cp -r ./assets ~/scripts
	# seed script dirs without clobbering an existing library (zdscripts.git etc.)
	cp -rn ./lua ~/scripts 2>/dev/null || true
	cp -rn ./python ~/scripts 2>/dev/null || true
	cp -rn ./sh ~/scripts 2>/dev/null || true

.PHONY: release
release:
	@echo building $(VERSION)...
	go build $(LDFLAGS) -o zd .
	@echo packaging...
	tar -cvJf zd-amd64-unix.tar.xz zd lua/ python/ sh/ assets/
