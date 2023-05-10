FROM golang:1.18-alpine

# Install git and mercurial
RUN apk add bash gcc g++ libc-dev

# Define environment variables
ENV APPNAME=task-manager-api\
    CGO_ENABLED=1 \
    WRKDIR=/app \
    GO111MODULE=on \
    USER=appuser \
    UID=10001
ENV TZ=Asia/Bangkok

# Never run a process as root in a container.
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR $WRKDIR

# copy and fetch external module
COPY go.mod go.sum ./

RUN go mod download && \
    go mod verify

# copy and build application
# removing debug informations and compile only for linux target and disabling cross compilation
COPY . .

# RUN GOOS=linux GOARCH=amd64 go build -tags musl -ldflags="-w -s" -o $APPNAME
RUN GOOS=linux go build -tags musl -ldflags="-w -s" -o $APPNAME 

USER root

RUN  echo "#!/bin/sh" >> $WRKDIR/start_service.sh && \
     echo "set -m" >> $WRKDIR/start_service.sh && \
     echo "exec $WRKDIR/$APPNAME" >> $WRKDIR/start_service.sh \
     chmod -R 777 $WRKDIR && \
     chown -R 1001:1001 $WRKDIR

RUN ["chmod", "+x", "./start_service.sh"]

USER 1001

EXPOSE 3000

CMD ["./start_service.sh"]