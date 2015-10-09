#!/bin/bash
DESTDIR=${DESTDIR:-/var/lib/gvweb/}
set -e
go build gvweb

mkdir -p "$DESTDIR"
mkdir -p "$DESTDIR"/data/
chown nobody "$DESTDIR"/data/


install  -p gvweb "$DESTDIR"
cp -ap static/ "$DESTDIR"


