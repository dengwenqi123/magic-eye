### Build

# This can be unified later, here for easy demos
########################################
### Tools & dependencies

check_tools:
	cd tools && $(MAKE) check_tools

update_tools:
	cd tools && $(MAKE) update_tools

get_tools:
	cd tools && $(MAKE) get_tools

get_vendor_deps:
	@rm -rf vendor/
	@echo "--> Running dep ensure"
	@dep ensure -v

install:
	go install ./cmd/mgycli
	go install ./cmd/mgycoind
