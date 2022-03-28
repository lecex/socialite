FROM bigrocs/golang-gcc:1.16 as builder

ARG ACCES_STOKEN
RUN apk add git
RUN go env -w GOPRIVATE=github.com/lecex
RUN git config --global url."https://bigrocs:${ACCES_STOKEN}@github.com".insteadOf "https://github.com"


WORKDIR /go/src/github.com/lecex/socialite
COPY . .

ENV GO111MODULE=on CGO_ENABLED=1 GOOS=linux GOARCH=amd64
RUN go build -a -installsuffix cgo -o bin/socialite

FROM bigrocs/alpine:ca-data

RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

COPY --from=builder /go/src/github.com/lecex/socialite/bin/socialite /usr/local/bin/
CMD ["socialite"]
EXPOSE 8080