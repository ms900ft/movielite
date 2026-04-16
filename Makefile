APP=movielite
HELPER=vlc_helper
WATCH_FILES= find . -type f -not -path '*/\.*' -not -path '*movieui\/*' | grep -i '.*[.]go\|html$$' 2> /dev/null

all:
	$(MAKE) staticbuild

test:
	go test --tags "fts5" $(test)

run:
	go run --tags "fts5" cmd/server/main.go start

build:
	go build -o ${APP}  --tags "fts5" ./cmd/server

wails:
	cd new-ui && node node_modules/vite/bin/vite.js build
	cp -r new-ui/dist cmd/wails/frontend
	go build -o ${APP}-desktop --tags "fts5,wails" ./cmd/wails

app: wails
	mkdir -p Movielite.app/Contents/MacOS
	mkdir -p Movielite.app/Contents/Resources
	cp ${APP}-desktop Movielite.app/Contents/MacOS/Movielite
	chmod +x Movielite.app/Contents/MacOS/Movielite
	cp cmd/wails/icon.png Movielite.app/Contents/Resources/icon.png
	@echo '<?xml version="1.0" encoding="UTF-8"?><!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd"><plist version="1.0"><dict><key>CFBundleExecutable</key><string>Movielite</string><key>CFBundleIdentifier</key><string>com.movielite.app</string><key>CFBundleName</key><string>Movielite</string><key>CFBundleDisplayName</key><string>Movielite</string><key>CFBundleVersion</key><string>1.0</string><key>CFBundleShortVersionString</key><string>1.0</string><key>CFBundlePackageType</key><string>APPL</string><key>CFBundleSignature</key><string>????</string><key>CFBundleInfoDictionaryVersion</key><string>6.0</string><key>LSMinimumSystemVersion</key><string>10.13</string><key>LSApplicationCategoryType</key><string>public.app-category.entertainment</string><key>CFBundleIconFile</key><string>icon.png</string><key>NSHighResolutionCapable</key><true/><key>NSSupportsAutomaticTermination</key><true/><key>NSSupportsSuddenTermination</key><true/><key>NSHumanReadableCopyright</key><string>Copyright © 2026. All rights reserved.</string></dict></plist>' > Movielite.app/Contents/Info.plist
	@echo "App bundle created: Movielite.app"

helper:
	go build -o ${HELPER}   ./cmd/helper

static:
	cd new-ui && node node_modules/vite/bin/vite.js build
	statik -src=$(shell pwd)/new-ui/dist -dest=./statik -f
	mv statik/statik/statik.go statik/statik.go 2>/dev/null || true
	rm -rf statik/statik 2>/dev/null || true

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
	 			scp "./${APP}" nudel2:/Users/ms/movielite
	 			rm ./${APP}
				ssh nudel2  launchctl unload ~/Library/LaunchAgents/org.local.movielite.plist
				ssh nudel2  launchctl load ~/Library/LaunchAgents/org.local.movielite.plist

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
