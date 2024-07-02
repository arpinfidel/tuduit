FROM golang:1.22

ENV GIT_TERMINAL_PROMPT=1

RUN apt update
RUN apt install -y parallel
RUN apt-get install zip unzip

RUN git config --global --add safe.directory '*'

RUN go install github.com/air-verse/air@latest
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1
RUN go install github.com/arpinfidel/autodoc@latest
# RUN go install -v github.com/go-delve/delve/cmd/dlv@latest
RUN go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN echo "[url \"git@github.com:\"]\n\tinsteadOf = https://github.com/" >> /root/.gitconfig

ENV GO111MODULE=on

WORKDIR /go/src/github.com/arpinfidel/tuduit
