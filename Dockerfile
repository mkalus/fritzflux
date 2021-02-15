# Build image
FROM golang:1.15-alpine AS build
WORKDIR /go/src/
COPY . github.com/mkalus/fritzflux
RUN cd github.com/mkalus/fritzflux && CGO_ENABLED=0 go build github.com/mkalus/fritzflux/cmd/fritzflux

# Mini image running binary
FROM scratch
COPY --from=build /go/src/github.com/mkalus/fritzflux/fritzflux /bin/fritzflux
ENTRYPOINT ["/bin/fritzflux"]
