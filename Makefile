all build:
	make -C pkg/driver
	make -C pkg/provisioner
.PHONY: all build

container:
	make -C pkg/driver container
	make -C pkg/provisioner container
.PHONY: container

push:
	make -C pkg/driver push
	make -C pkg/provisioner push
.PHONY: push

clean:
	make -C pkg/driver clean
	make -C pkg/provisioner clean
.PHONY: clean
