.PHONY: all bin stdlib_release stdlib_debug clean

all:
	@$(MAKE) clear
	@$(MAKE) get_deps
	@$(MAKE) stdlib_release
	@$(MAKE) bin

get_deps:
	@echo
	@echo "\033[4m\033[1mGetting go deps\033[0m"
	@echo
	go get ./...

bin:
	@echo
	@echo "\033[4m\033[1mBuilding binary\033[0m"
	@echo
	go build -o bin/hummus

stdlib_release:
	@echo
	@echo "\033[4m\033[1mBuilding stdlib for release\033[0m"
	@echo
	chmod +x scripts/prepare_release.sh
	./scripts/prepare_release.sh

stdlib_debug:
	@echo
	@echo "\033[4m\033[1mBuilding stdlib for debug\033[0m"
	@echo
	chmod +x scripts/prepare_debug.sh
	./scripts/prepare_debug.sh

clear:
	@echo
	@echo "\033[4m\033[1mClearing output folder\033[0m"
	@echo
	rm -rf bin/ || true
