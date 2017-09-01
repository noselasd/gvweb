FROM golang:1.9-stretch

ENV SRC_PATH ${SRC_PATH:-./}

# Copy the app to the image
COPY ${SRC_PATH} /go/src/gvweb

RUN\
	apt-get update &&\
	apt-get -y install graphviz &&\
	cd /go/src/gvweb &&\
    	make

ENV PORT=80
EXPOSE $PORT
WORKDIR /go/src/gvweb

CMD ["./gvweb", "-port", "80"]
