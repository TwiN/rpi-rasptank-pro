pbr: update build-and-run

build-and-run:
	go build -mod vendor && ./rpi-rasptank-pro

update:
	git pull