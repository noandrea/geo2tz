#!/bin/bash

VERSION=2023b
WORKDIR=dist

pushd $WORKDIR
echo download timezones
curl -LO https://github.com/evansiroky/timezone-boundary-builder/releases/download/$VERSION/timezones-with-oceans.geojson.zip
unzip timezones-with-oceans.geojson.zip
./geo2tz build --json "combined-with-oceans.json" --db=timezone
echo "https://github.com/evansiroky/timezone-boundary-builder" > SOURCE
echo "tzdata v$VERSION" >> SOURCE

popd
mv dist/timezone.snap.json tzdata/timezone.snap.json
mv dist/SOURCE tzdata/SOURCE