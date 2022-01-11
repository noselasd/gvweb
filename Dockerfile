FROM golang:1.17-stretch

ENV SRC_PATH ${SRC_PATH:-./}

# Copy the app to the image
COPY ${SRC_PATH} /go/src/gvweb

RUN set -ex && \
	useradd -r gvweb &&\
	apt-get update &&\
	apt-get -y install graphviz &&\
	cd /go/src/gvweb &&\
    	make &&\
	mkdir -p /app/data &&\
	chown -R gvweb:gvweb /app &&\
	cp -a gvweb static /app/ &&\
	rm -rf /go/src/*



EXPOSE 8080
WORKDIR /app
USER gvweb:gvweb

CMD ["./gvweb", "-port", "8080"]

