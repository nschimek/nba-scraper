#!/usr/bin/env bash
if [ -z "$0" ]; then 
  echo "Usage: $0 <version>"
  exit 1
fi
VERSION=$0
echo "Creating git tag ${VERSION}..."
git tag ${VERSION}
echo "Building and zipping Linux binary ${VERSION}..."
GOOS=linux GOARCH=amd64 go build -ldflags "-X 'github.com/nschimek/nba-scraper/core.Version=${VERSION}'"
zip -r nba-scraper-linux-${VERSION}.zip nba-scraper config/sample.yaml
echo "Building and zipping Windows binary ${VERSION}..."
GOOS=windows GOARCH=amd64 go build -ldflags "-X 'github.com/nschimek/nba-scraper/core.Version=${VERSION}'"
zip -r nba-scraper-win-${VERSION}.zip nba-scraper.exe config/sample.yaml
git push origin ${VERSION}