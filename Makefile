TAR?=tar
VERSION=$(shell git describe --tags --long --abbrev=4 --dirty=-D)
DISTDIR=gvweb-$(VERSION)

.PHONY: gvweb dist gvweb-$(VERSION)
gvweb:
	go build --ldflags "-X main.gVersion=$(VERSION)" $@

gvweb-$(VERSION):
	mkdir -p $(DISTDIR)
	mkdir -p $(DISTDIR)/etc/
	cp gvweb $(DISTDIR)
	mkdir -p $(DISTDIR)/data/
	mkdir -p $(DISTDIR)/static/
	cp -ap etc/gvweb.service $(DISTDIR)/etc/
	cp -ap static/* $(DISTDIR)/static/
	${TAR} --owner=nobody --group=nobody -cvzf $(DISTDIR).tar.gz $(DISTDIR)
	rm -rf $(DISTDIR)

dist: gvweb gvweb-$(VERSION)
