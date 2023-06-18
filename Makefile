
.PHONY: build
build:
	go build -o zd .

.PHONY: install
install:
	mkdir -p ~/scripts || true
	cp ./zd ~/scripts
