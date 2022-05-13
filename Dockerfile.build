FROM golang:1.18 AS build-stage
RUN wget https://github.com/coredns/coredns/archive/refs/tags/v1.9.2.tar.gz
RUN wget https://github.com/vlcty/coredns-auto-ipv6-ptr/archive/refs/heads/master.tar.gz
RUN tar xzf v1.9.2.tar.gz
RUN tar xzf master.tar.gz
WORKDIR /go/coredns-1.9.2
RUN git apply /go/coredns-auto-ipv6-ptr-master/file-fallthrough.patch
RUN ln -s /go/coredns-auto-ipv6-ptr-master plugin/autoipv6ptr
RUN echo "autoipv6ptr:autoipv6ptr" >> plugin.cfg
RUN go generate coredns.go
RUN go get
RUN go build -o coredns .

FROM scratch AS export-stage
COPY --from=build-stage /go/coredns-1.9.2/coredns /
