PREFIX=/usr/local

format:
	gofmt -s -w .

build: format
	GOOS=linux GOARCH=mipsle go build -trimpath -mod=readonly -modcacherw -ldflags="-s -w"

install: build
	mkdir -p $(PREFIX)/bin
	cp -f modem-exporter $(PREFIX)/bin

install-service:
	cp -f assets/prometheus-edgemax-exporter.service /etc/systemd/system/
