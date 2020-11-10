pbr: update build-and-run
br: build-and-run

build-and-run:
	go build -mod vendor && sudo ./rpi-rasptank-pro

update:
	git pull

i2c:
	i2cdetect -y 1

gpio:
	raspi-gpio get