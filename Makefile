all build:
	make -C cmd/driver
	make -C cmd/provisioner
.PHONY: all build

container:
	make -C cmd/driver container
	make -C cmd/provisioner container
.PHONY: container

push:
	make -C cmd/driver push
	make -C cmd/provisioner push
.PHONY: push

clean:
	make -C cmd/driver clean
	make -C cmd/provisioner clean
.PHONY: clean
