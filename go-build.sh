#!/usr/bin/env bash
VERSION=$(git tag --sort=-version:refname | head -n 1)
echo "Building and zipping Linux binary ${VERSION}..."
GOOS=linux GOARCH=amd64 go build -ldflags "-X 'github.com/nschimek/nba-scraper/core.Version=${VERSION}'"
zip -r nba-scraper-linux-${VERSION}.zip nba-scraper config/sample.yaml
echo "Building and zipping Windows binary ${VERSION}..."
GOOS=windows GOARCH=amd64 go build -ldflags "-X 'github.com/nschimek/nba-scraper/core.Version=${VERSION}'"
zip -r nba-scraper-win-${VERSION}.zip nba-scraper.exe config/sample.yaml
git push origin ${VERSION}