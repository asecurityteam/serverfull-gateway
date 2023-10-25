TAG := $(shell git rev-parse --short HEAD)
DIR := $(shell pwd -L)

# SDCLI
SDCLI_VERSION=v1.5
SDCLI=docker run -ti \
	--mount src="$(DIR)",target="$(DIR)",type="bind" \
	-w "$(DIR)" \
	-e "GOFLAGS= -buildvcs=false" \
	asecurityteam/sdcli:$(SDCLI_VERSION)

dep:
	$(SDCLI) go dep

lint:
	$(SDCLI) go lint

test:
	$(SDCLI) go test

integration:
	DIR=$(DIR) \
	docker-compose \
		-f docker-compose.it.yml \
		up \
			--abort-on-container-exit \
			--build \
			--exit-code-from test

coverage:
	$(SDCLI) go coverage

doc: ;

build-dev: ;

build: ;

run: ;

deploy-dev: ;

deploy: ;
