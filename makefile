OUT_DIR := dist/mop
TS_CORE_SRC := $(shell find ui/core -name '*.ts' -type f)
ASSETS_INPUT := $(shell find assets/ -type f)
ASSETS := $(patsubst assets/%,$(OUT_DIR)/assets/%,$(ASSETS_INPUT))
# Recursive wildcard function. Needs to be '=' instead of ':=' because of recursion.
rwildcard = $(foreach d,$(wildcard $(1:=/*)),$(call rwildcard,$d,$2) $(filter $(subst *,%,$2),$d))
GOROOT := $(shell go env GOROOT)
UI_SRC := $(shell find ui -name '*.ts' -o -name '*.tsx' -o -name '*.scss' -o -name '*.html')
PAGE_INDECES := ui/death_knight/blood/index.html \
				ui/death_knight/frost/index.html \
				ui/death_knight/unholy/index.html \
				ui/druid/balance/index.html \
				ui/druid/feral/index.html \
				ui/druid/guardian/index.html \
				ui/druid/restoration/index.html \
				ui/hunter/beast_mastery/index.html \
				ui/hunter/marksmanship/index.html \
				ui/hunter/survival/index.html \
				ui/mage/arcane/index.html \
				ui/mage/fire/index.html \
				ui/mage/frost/index.html \
				ui/monk/brewmaster/index.html \
				ui/monk/mistweaver/index.html \
				ui/monk/windwalker/index.html \
				ui/paladin/holy/index.html \
				ui/paladin/protection/index.html \
				ui/paladin/retribution/index.html \
				ui/priest/discipline/index.html \
				ui/priest/holy/index.html \
				ui/priest/shadow/index.html \
				ui/rogue/assassination/index.html \
				ui/rogue/combat/index.html \
				ui/rogue/subtlety/index.html \
				ui/shaman/elemental/index.html \
				ui/shaman/enhancement/index.html \
				ui/shaman/restoration/index.html \
				ui/warlock/affliction/index.html \
				ui/warlock/demonology/index.html \
				ui/warlock/destruction/index.html \
				ui/warrior/arms/index.html \
				ui/warrior/fury/index.html \
				ui/warrior/protection/index.html \
				ui/raid/full/index.html \
				ui/results/detailed/index.html

$(OUT_DIR)/.dirstamp: \
  $(OUT_DIR)/lib.wasm \
  ui/core/proto/api.ts \
  $(ASSETS) \
  $(OUT_DIR)/bundle/.dirstamp
	touch $@

$(OUT_DIR)/bundle/.dirstamp: \
  $(UI_SRC) \
  $(PAGE_INDECES) \
  vite.config.mts \
  vite.build-workers.ts \
  node_modules \
  tsconfig.json \
  ui/core/index.ts \
  ui/core/proto/api.ts
	npx tsc --noEmit
	npx tsx vite.build-workers.ts
	npx vite build
	touch $@

ui/core/index.ts: $(TS_CORE_SRC)
	find ui/core -name '*.ts' | \
	  awk -F 'ui/core/' '{ print "import \x22./" $$2 "\x22;" }' | \
	  sed 's/\.ts";$$/";/' | \
	  grep -v 'import "./index";' > $@

.PHONY: clean
clean:
	rm -rf ui/core/proto/*.ts \
	  sim/core/proto/*.pb.go \
	  wowsimmop \
	  wowsimmop-windows.exe \
	  wowsimmop-amd64-darwin \
	  wowsimmop-arm64-darwin \
	  wowsimmop-amd64-linux \
	  dist \
	  binary_dist \
	  ui/core/index.ts \
	  ui/core/proto/*.ts \
	  node_modules \
	  $(PAGE_INDECES)
	find . -name "*.results.tmp" -type f -delete

ui/core/proto/api.ts: proto/*.proto node_modules
	npx protoc --ts_opt generate_dependencies --ts_out ui/core/proto --proto_path proto proto/api.proto
	npx protoc --ts_out ui/core/proto --proto_path proto proto/test.proto
	npx protoc --ts_out ui/core/proto --proto_path proto proto/ui.proto

ui/%/index.html: ui/index_template.html
	cat ui/index_template.html | sed -e 's/@@CLASS@@/$(shell dirname $(@D) | xargs basename)/g' -e 's/@@SPEC@@/$(shell basename $(@D))/g' > $@

.PHONY: package.json

package.json:
# Checks if the system is FreeBSD and jq is installed. This is due to the need to switch out the vite package for rollup on FreeBSD.
ifeq ($(shell uname -s), FreeBSD)
	@if ! command -v jq > /dev/null; then \
		echo "jq is not installed. Please install jq to proceed."; \
		exit 1; \
	fi; \
	\
	echo "Checking and updating package.json for FreeBSD..."; \
	\
	if ! grep -q '"overrides"' package.json; then \
		jq '. + { "overrides": { "vite": { "rollup": "npm:@rollup/wasm-node@4.13.0" } } }' package.json > package.json.tmp && mv package.json.tmp package.json && npm install; \
	else \
		jq '.overrides += { "vite": { "rollup": "npm:@rollup/wasm-node@4.13.0" } }' package.json > package.json.tmp && mv package.json.tmp package.json && npm install; \
	fi
endif

package-lock.json:
	npm install

node_modules: package-lock.json
	npm ci

# Generic rule for hosting any class directory
.PHONY: host_%
host_%: $(OUT_DIR) node_modules
	npx http-server $(OUT_DIR)/..

# Generic rule for building index.html for any class directory
$(OUT_DIR)/%/index.html: ui/index_template.html $(OUT_DIR)/assets
	$(eval title := $(shell echo $(shell basename $(@D)) | sed -r 's/(^|_)([a-z])/\U \2/g' | cut -c 2-))
	echo $(title)
	mkdir -p $(@D)
	cat ui/index_template.html | sed -e 's/@@CLASS@@/$(shell dirname $((@D)) | xargs basename)/g' -e 's/@@SPEC@@/$(shell basename $(@D))/g' > $@

.PHONY: wasm
wasm: $(OUT_DIR)/lib.wasm

# Builds the generic .wasm, with all items included.
$(OUT_DIR)/lib.wasm: sim/wasm/* sim/core/proto/api.pb.go $(filter-out sim/core/items/all_items.go, $(call rwildcard,sim,*.go))
	@echo "Starting webassembly compile now..."
	@if GOOS=js GOARCH=wasm go build -o ./$(OUT_DIR)/lib.wasm ./sim/wasm/; then \
		printf "\033[1;32mWASM compile successful.\033[0m\n"; \
	else \
		printf "\033[1;31mWASM COMPILE FAILED\033[0m\n"; \
		exit 1; \
	fi

$(OUT_DIR)/assets/%: assets/%
	mkdir -p $(@D)
	cp $< $@
	rm -rf $(OUT_DIR)/assets/db_inputs


binary_dist/dist.go: sim/web/dist.go.tmpl
	mkdir -p binary_dist/mop
	touch binary_dist/mop/embedded
	cp sim/web/dist.go.tmpl binary_dist/dist.go

binary_dist: $(OUT_DIR)/.dirstamp
	rm -rf binary_dist
	mkdir -p binary_dist
	cp -r $(OUT_DIR) binary_dist/
	rm binary_dist/mop/lib.wasm
	rm -rf binary_dist/mop/assets/db_inputs
	rm binary_dist/mop/assets/database/db.bin
	rm binary_dist/mop/assets/database/leftover_db.bin

# Rebuild the protobuf generated code.
.PHONY: proto
proto: sim/core/proto/api.pb.go ui/core/proto/api.ts

# Builds the web server with the compiled client.
.PHONY: wowsimmop
wowsimmop: binary_dist devserver

.PHONY: devserver
devserver: sim/core/proto/api.pb.go sim/web/main.go binary_dist/dist.go
	@echo "Starting server compile now..."
	@if go build -o wowsimmop ./sim/web/main.go ; then \
		printf "\033[1;32mBuild Completed Successfully\033[0m\n"; \
	else \
		printf "\033[1;31mBUILD FAILED\033[0m\n"; \
		exit 1; \
	fi

.PHONY: air
air:
ifeq ($(WATCH), 1)
	@if ! command -v air; then \
		echo "Missing air dependency. Please run \`make setup\`"; \
		exit 1; \
	fi
endif

rundevserver: air devserver
ifeq ($(WATCH), 1)
	npx tsx vite.build-workers.ts & npx vite build -m development --watch &
	ulimit -n 10240 && air -tmp_dir "/tmp" -build.include_ext "go,proto" -build.args_bin "--usefs=true --launch=false" -build.bin "./wowsimmop" -build.cmd "make devserver" -build.exclude_dir "assets,dist,node_modules,ui,tools"
else
	./wowsimmop --usefs=true --launch=false --host=":3333"
endif

wowsimmop-windows.exe: wowsimmop
# go build only considers syso files when invoked without specifying .go files: https://github.com/golang/go/issues/16090
	cp ./assets/favicon_io/icon-windows_amd64.syso ./sim/web/icon-windows_amd64.syso
	cd ./sim/web/ && GOOS=windows GOARCH=amd64 GOAMD64=v2 go build -o wowsimmop-windows.exe -ldflags="-X 'main.Version=$(VERSION)' -s -w"
	cd ./cmd/wowsimcli && GOOS=windows GOARCH=amd64 GOAMD64=v2 go build -o wowsimcli-windows.exe --tags=with_db -ldflags="-X 'main.Version=$(VERSION)' -s -w"
	rm ./sim/web/icon-windows_amd64.syso
	mv ./sim/web/wowsimmop-windows.exe ./wowsimmop-windows.exe
	mv ./cmd/wowsimcli/wowsimcli-windows.exe ./wowsimcli-windows.exe

release: wowsimmop wowsimmop-windows.exe
	GOOS=darwin GOARCH=amd64 GOAMD64=v2 go build -o wowsimmop-amd64-darwin -ldflags="-X 'main.Version=$(VERSION)' -s -w" ./sim/web/main.go
	GOOS=darwin GOARCH=arm64 go build -o wowsimmop-arm64-darwin -ldflags="-X 'main.Version=$(VERSION)' -s -w" ./sim/web/main.go
	GOOS=linux GOARCH=amd64 GOAMD64=v2 go build -o wowsimmop-amd64-linux   -ldflags="-X 'main.Version=$(VERSION)' -s -w" ./sim/web/main.go
	GOOS=linux GOARCH=amd64 GOAMD64=v2 go build -o wowsimcli-amd64-linux --tags=with_db -ldflags="-X 'main.Version=$(VERSION)' -s -w" ./cmd/wowsimcli/cli_main.go
# Now compress into a zip because the files are getting large.
	zip wowsimmop-windows.exe.zip wowsimmop-windows.exe
	zip wowsimmop-amd64-darwin.zip wowsimmop-amd64-darwin
	zip wowsimmop-arm64-darwin.zip wowsimmop-arm64-darwin
	zip wowsimmop-amd64-linux.zip wowsimmop-amd64-linux
	zip wowsimcli-amd64-linux.zip wowsimcli-amd64-linux
	zip wowsimcli-windows.exe.zip wowsimcli-windows.exe

sim/core/proto/api.pb.go: proto/*.proto
	protoc -I=./proto --go_out=./sim/core ./proto/*.proto

# Only useful for building the lib on a host platform that matches the target platform
.PHONY: locallib
locallib: sim/core/proto/api.pb.go
	go build -buildmode=c-shared -o wowsimmop.so --tags=with_db ./sim/lib/library.go

.PHONY: nixlib
nixlib: sim/core/proto/api.pb.go
	GOOS=linux GOARCH=amd64 GOAMD64=v2 go build -buildmode=c-shared -o wowsimmop-linux.so --tags=with_db ./sim/lib/library.go

.PHONY: winlib
winlib: sim/core/proto/api.pb.go
	GOOS=windows GOARCH=amd64 GOAMD64=v2 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -buildmode=c-shared -o wowsimmop-windows.dll --tags=with_db ./sim/lib/library.go

.PHONY: simdb
simdb: sim/core/items/all_items.go sim/core/proto/api.pb.go

CLIENTDATA_SETTINGS := $(shell realpath ./tools/database/generator-settings.json)
CLIENTDATAPTR_SETTINGS := $(shell realpath ./tools/database/ptr-generator-settings.json)
CLIENTDATA_OUTPUT   := $(shell realpath ./tools/database/wowsims.db)

.PHONY: db
db:
	@echo "Running DB2ToSqlite for clientdata"
	cd tools/DB2ToSqlite && dotnet run -- -s $(CLIENTDATA_SETTINGS) --output $(CLIENTDATA_OUTPUT)
	@echo "Running DBC generation tool"
	go run tools/database/gen_db/*.go -outDir=./assets -gen=db

.PHONY: ptrdb
ptrdb:
	@echo "Running DB2ToSqlite for clientdata"
	cd tools/DB2ToSqlite && dotnet run -- -s $(CLIENTDATAPTR_SETTINGS) --output $(CLIENTDATA_OUTPUT)
	@echo "Running DBC generation tool"
	go run tools/database/gen_db/*.go -outDir=./assets -gen=db

sim/core/items/all_items.go: $(call rwildcard,tools/database,*.go) $(call rwildcard,sim/core/proto,*.go)
	go run tools/database/gen_db/*.go -outDir=./assets -gen=db

.PHONY: test
test: $(OUT_DIR)/lib.wasm binary_dist/dist.go
	go test --tags=with_db ./sim/...

.PHONY: update-tests
update-tests:
	find . -name "*.results" -type f -delete
	find . -name "*.results.tmp" -exec bash -c 'cp "$$1" "$${1%.results.tmp}".results' _ {} \;

.PHONY: fmt
fmt: tsfmt
	gofmt -w ./sim
	gofmt -w ./tools

.PHONY: tsfmt
tsfmt:
	for dir in $$(find ./ui -maxdepth 1 -type d -not -path "./ui" -not -path "./ui/worker"); do \
		echo $$dir; \
		npx tsfmt -r --useTsfmt ./tsfmt.json --baseDir $$dir; \
	done

# one time setup to install pre-commit hook for gofmt and npm install needed packages
setup:
	cp pre-commit .git/hooks
	chmod +x .git/hooks/pre-commit
	! command -v air && curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin || true

# Host a local server, for dev testing
.PHONY: host
host: air $(OUT_DIR)/.dirstamp node_modules
ifeq ($(WATCH), 1)
	ulimit -n 10240 && air -tmp_dir "/tmp" -build.include_ext "go,ts,js,html" -build.bin "npx" -build.args_bin "http-server $(OUT_DIR)/.." -build.cmd "make" -build.exclude_dir "dist,node_modules,tools"
else
	# Intentionally serve one level up, so the local site has 'mop' as the first
	# directory just like github pages.
	npx http-server $(OUT_DIR)/..
endif

devmode: air devserver
ifeq ($(WATCH), 1)
	npx tsx vite.build-workers.ts & npx vite serve --host &
	air -tmp_dir "/tmp" -build.include_ext "go,proto" -build.args_bin "--usefs=true --launch=false --wasm=false" -build.bin "./wowsimmop" -build.cmd "make devserver" -build.exclude_dir "assets,dist,node_modules,ui,tools"
else
	./wowsimmop --usefs=true --launch=false --host=":3333"
endif

webworkers:
	npx tsx vite.build-workers.ts --watch=$(if $(WATCH),true,false)
