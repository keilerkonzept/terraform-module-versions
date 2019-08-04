FROM golang:1.10-alpine3.7 AS build
RUN apk add --no-cache make
WORKDIR /go/src/github.com/sgreben/versions/
COPY vendor vendor
COPY cmd cmd
COPY pkg pkg
COPY Makefile Makefile
ENV CGO_ENABLED=0
RUN make binaries/linux_x86_64/versions && mv binaries/linux_x86_64/versions /app

FROM scratch
COPY --from=build /app /versions
ENTRYPOINT [ "/versions" ]
