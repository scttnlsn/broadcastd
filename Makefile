PREFIX=/usr/local
BINDIR=${PREFIX}/bin

all: build/broadcastd

build:
	mkdir build

build/broadcastd: build
	go build -o build/broadcastd

clean:
	rm -rf build

install: build/broadcastd
	install -m 755 -d ${BINDIR}
	install -m 755 build/broadcastd ${BINDIR}/broadcastd

uninstall:
	rm ${BINDIR}/broadcastd

test:
	cd broadcastd; go test

.PHONY: install uninstall clean all test