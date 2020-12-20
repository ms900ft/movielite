APP=movielite
HELPER=vlc_helper
WATCH_FILES= find . -type f -not -path '*/\.*' | grep -i '.*[.]go\|html$$' 2> /dev/null

all:
	$(MAKE) staticbuild

test:
	go test --tags "fts5" $(test)

run:
	go run --tags "fts5" cmd/server/main.go start

build:
	go build -o ${APP}  --tags "fts5" ./cmd/server

helper:
	go build -o ${HELPER}   ./cmd/helper

static:
	$(MAKE) -C movieui
	statik  -src=$(shell pwd)/movieui/dist

serve:
	$(MAKE) -C movieui serve

clean:
	rm -f $(APP)
	rm -f statik/*
	$(MAKE) -C movieui clean


staticbuild:
	$(MAKE) static
	$(MAKE) build

runbuild:
	go build ./...; ./${APP}

lint:
	golint

install:
	go install  ./...

npminstall:
	$(MAKE) -C movieui install

deploy: build
	 			scp "./${APP}" nudel:/Users/github.com/ms900ft/movielite
	 			rm ./${APP}
				ssh nudel  launchctl unload ~/Library/LaunchAgents/org.local.movielite.plist
				ssh nudel  launchctl load ~/Library/LaunchAgents/org.local.movielite.plist

metalint:
	if command -v gometalinter > /dev/null; then echo ''; else go get -u github.com/alecthomas/gometalinter; fi
	gometalinter ./...

entr_warn:
	@echo "----------------------------------------------------------"
	@echo "     ! File watching functionality non-operational !      "
	@echo ""
	@echo "Install entr(1) to automatically run tasks on file change."
	@echo "See http://entrproject.org/"
	@echo "----------------------------------------------------------"

watch_echo:
	echo `${WATCH_FILES}`

watch_test:
	if command -v entr > /dev/null; then ${WATCH_FILES} | entr -rc $(MAKE) test; else $(MAKE) test entr_warn; fi

watch_run:
	if command -v entr > /dev/null; then ${WATCH_FILES} | entr -rc $(MAKE) run; else $(MAKE) run entr_warn; fi

watch_build:
	if command -v entr > /dev/null; then ${WATCH_FILES} | entr -rc $(MAKE) build; else $(MAKE) build entr_warn; fi

watch_runbuild:
	if command -v entr > /dev/null; then ${WATCH_FILES} | entr -rc $(MAKE) runbuild; else $(MAKE) runbuild entr_warn; fi

watch_lint:
	if command -v entr > /dev/null; then ${WATCH_FILES} | entr -rc $(MAKE) lint; else $(MAKE) lint entr_warn; fi

watch_metalint:
	if command -v entr > /dev/null; then ${WATCH_FILES} | entr -rc $(MAKE) metalint; else $(MAKE) metalint entr_warn; fi
