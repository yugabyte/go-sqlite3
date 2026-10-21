#!/bin/sh
# Syncs sqlite-vec's source from github.com/yugabyte/sqlite-vec into sqlite-vec/, for the
# sqlite_vec build tag (sqlite3_opt_vec.go). Usage: upgrade/sqlite-vec.sh <tag or commit>
#
# sqlite-vec.h is rendered from sqlite-vec.h.tmpl the way sqlite-vec's own Makefile does, with
# DATE and SOURCE taken from the synced commit.

set -e

REF=$1
if [ -z "$REF" ]; then
  echo "usage: $0 <sqlite-vec tag or commit>" >&2
  exit 2
fi
REPO=${SQLITE_VEC_REPO:-https://github.com/yugabyte/sqlite-vec}

cd "$(dirname "$0")/.."
DEST=$PWD/sqlite-vec
TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

git -C "$TMP" init -q
git -C "$TMP" fetch -q --depth 1 "$REPO" "$REF"
git -C "$TMP" checkout -q FETCH_HEAD

VERSION=$(cat "$TMP/VERSION")
SOURCE=$(git -C "$TMP" rev-parse HEAD)
DATE=$(TZ=UTC0 git -C "$TMP" log -1 --date=format-local:'%Y-%m-%dT%H:%M:%SZ+0000' --format=%cd)
MAJOR=$(echo "$VERSION" | cut -d. -f1)
MINOR=$(echo "$VERSION" | cut -d. -f2)
PATCH=$(echo "$VERSION" | cut -d. -f3 | cut -d- -f1)

# Replace everything but the shim (sqlite3.h), so files sqlite-vec drops don't linger.
find "$DEST" -type f ! -name sqlite3.h -exec rm -f {} +
cp "$TMP"/sqlite-vec.c "$TMP"/LICENSE-MIT "$TMP"/LICENSE-APACHE "$DEST"/
for f in "$TMP"/sqlite-vec-*.c; do
  [ -e "$f" ] && cp "$f" "$DEST"/
done
sed -e "s|\${VERSION}|$VERSION|g" \
    -e "s|\${DATE}|$DATE|g" \
    -e "s|\${SOURCE}|$SOURCE|g" \
    -e "s|\${VERSION_MAJOR}|$MAJOR|g" \
    -e "s|\${VERSION_MINOR}|$MINOR|g" \
    -e "s|\${VERSION_PATCH}|$PATCH|g" \
    "$TMP/sqlite-vec.h.tmpl" > "$DEST/sqlite-vec.h"
printf 'sqlite-vec %s\nsource: %s %s (%s)\n' "$VERSION" "$REPO" "$REF" "$SOURCE" > "$DEST/VERSION"

echo "sqlite-vec $VERSION ($SOURCE) synced into sqlite-vec/"
