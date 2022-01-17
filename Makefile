dir_bin=bin
dir_release=release
docker_image_name=dewanggasurya/fizz_buzz
os_list=linux windows
arch_list=amd64 386 arm64 arm
docker_arch_list=amd64 arm32v7 arm64v8

CMDS := server client
build: pre_build build_cmd post_build
pre_build:
	# start building ...
	@rm -rf $(dir_release)/*
post_build:		
	# done building ...
build_cmd: $(CMDS) # building all listed cmd
$(CMDS):
	@for GOOS in $(os_list); do \
		for GOARCH in $(arch_list); do \
			mkdir -p "$(dir_release)/$$GOOS/$$GOARCH/"; \
			EXT='' ; \
			TAGS='' ; \
			if [ "$$GOOS" = "windows"  ]; then\
			if [ "$$GOARCH" = "arm" ] || [ "$$GOARCH" = "arm64" ]; then\
			continue ; \
			fi; \
			EXT='.exe' ; \
			TAGS='release' ; \
			else \
			EXT='' ; \
			TAGS='release linux' ; \
			fi; \
			cd "./example/cmd/$@" && GOOS=$$GOOS GOARCH=$$GOARCH go build -o "$@$$EXT" -tags "$$TAGS"; \
			cd ../..; \
			echo "## building $@ <$$GOOS> <$$GOARCH>"; \
			mv "./example/cmd/$@/$@$$EXT" "$(dir_release)/$$GOOS/$$GOARCH/"; \
			# @chmod +x "$(dir_release)/$$GOOS/$$GOARCH/$@"; \
		done \
	done
test:
	# testing ...
	@go test -v ./...
clean:
	# cleaning ...
	@rm -rf "$(dir_release)" 
	@rm -rf "$(dir_bin)" 
init:
	# initializing ...
	@go mod download
	@go mod tidy
