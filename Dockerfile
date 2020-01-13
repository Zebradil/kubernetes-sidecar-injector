FROM golang:1.13-alpine3.11 as build

RUN apk --no-cache update \
 && apk --no-cache add make git bash gcc libc-dev

RUN go get -u github.com/golang/dep/cmd/dep \
 && go get golang.org/x/tools/cmd/goimports \
 && go get -u golang.org/x/lint/golint

WORKDIR /go/src/github.com/expediadotcom/kubernetes-sidecar-injector

COPY . .

RUN dep ensure \
 && go vet ./... \
 && go list ./... | xargs golint -min_confidence 1.0 \
 && CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o /kubernetes-sidecar-injector


FROM scratch

COPY --from=build /kubernetes-sidecar-injector /kubernetes-sidecar-injector

ENTRYPOINT ["/kubernetes-sidecar-injector"]
