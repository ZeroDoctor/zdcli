
.PHONY: build
build:
	go build -o zd .

.PHONY: install
install:
	mkdir -p ~/scripts || true
	cp ./zd ~/scripts
	cp -r ./assets ~/scripts

.PHONY: release
release:
	
	@echo update commit version
	COMMIT=$$(git rev-parse --short HEAD) && \
				 sed -i "s/{{ .PRE-RELEASE }}/$${COMMIT}/g" version.go && \
				 go build -o zd . && \
				 sed -i "s/$${COMMIT}/{{ .PRE-RELEASE }}/g" version.go

	@echo packaging...
	tar -cvJf zd-amd64-unix.tar.xz zd lua/ assets/
