FROM golang:alpine AS build

WORKDIR /go/src/github.com/cdrpl/granny/server

# Copy project to build dir and install
COPY . .
RUN go install .

# Empty alpine image
FROM alpine

WORKDIR /go/src/github.com/cdrpl/granny/server

# Expose port
EXPOSE 3000

# Copy bin to new image
COPY --from=build /go/bin /go/bin
ENTRYPOINT [ "/go/bin/server", "-e" ]
