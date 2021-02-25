FROM golang:1.15-alpine3.12 AS build

WORKDIR /go/src/github.com/cdrpl/server

# Copy project to build dir and install
COPY . .
RUN go install .

# Empty alpine image
FROM alpine:3.12

WORKDIR /go/src/github.com/cdrpl/server

# Default env
ENV ENV=production
ENV PORT=3000

# Expose port
EXPOSE 3000

# Copy the migrations and html
COPY --from=build /go/src/github.com/cdrpl/server/migration ./migration
COPY --from=build /go/src/github.com/cdrpl/server/html ./html

# Copy bin to new image
COPY --from=build /go/bin /go/bin
ENTRYPOINT [ "/go/bin/server" ]
