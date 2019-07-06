#!/bin/bash -e

REMOTE_ADDR=lenslocked.myslidekit.com
REMOTE_GOPATH=/root/go
REMOTE_SRC_DIR=$REMOTE_GOPATH/src/github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com
REMOTE_APP_DIR=/root/lenslocked
REMOTE_GO_BIN=/usr/local/go/bin

cd "$GOPATH/src/github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com"

echo "==== Releasing lenslocked.com ===="
echo "  - Deleting the local binary if it exists..."
rm server || true
echo "  - Local binary deleted successfully!"

echo "  - Deleting existing source code..."
ssh root@$REMOTE_ADDR "rm -rf $REMOTE_SRC_DIR"
echo "  - Code deleted successfully!"

echo "  - Uploading code..."
ssh root@$REMOTE_ADDR "mkdir -p $REMOTE_SRC_DIR"
rsync -avr \
    --exclude '.git/*' --exclude '.gitignore' --exclude 'tmp/*' \
    --exclude 'images/*' --exclude 'exp/' --exclude '.config' \
    --exclude 'scripts/' --exclude 'tmp/' --exclude 'runner.conf' \
    ./ root@$REMOTE_ADDR:$REMOTE_SRC_DIR
echo "  - Code uploaded successfully!"

echo "  - Go getting dependencies..."
ssh root@$REMOTE_ADDR "export GOPATH=$REMOTE_GOPATH; \
    $REMOTE_GO_BIN/go get golang.org/x/crypto/bcrypt"
ssh root@$REMOTE_ADDR "export GOPATH=$REMOTE_GOPATH; \
    $REMOTE_GO_BIN/go get github.com/jinzhu/gorm"
ssh root@$REMOTE_ADDR "export GOPATH=$REMOTE_GOPATH; \
    $REMOTE_GO_BIN/go get github.com/lib/pq"
ssh root@$REMOTE_ADDR "export GOPATH=$REMOTE_GOPATH; \
    $REMOTE_GO_BIN/go get github.com/gorilla/schema"
ssh root@$REMOTE_ADDR "export GOPATH=$REMOTE_GOPATH; \
    $REMOTE_GO_BIN/go get github.com/gorilla/mux"
ssh root@$REMOTE_ADDR "export GOPATH=$REMOTE_GOPATH; \
    $REMOTE_GO_BIN/go get github.com/gorilla/csrf"
ssh root@$REMOTE_ADDR "export GOPATH=$REMOTE_GOPATH; \
    $REMOTE_GO_BIN/go get gopkg.in/mailgun/mailgun-go.v3"

echo "  - Building source code on remote server..."
ssh root@$REMOTE_ADDR "export GOPATH=$REMOTE_GOPATH; \
    cd $REMOTE_APP_DIR; \
    $REMOTE_GO_BIN/go build -o ./server $REMOTE_SRC_DIR/main.go"
echo "  - Code built successfully!"

echo "  - Moving assets..."
ssh root@$REMOTE_ADDR "cd $REMOTE_APP_DIR; \
    cp -R $REMOTE_SRC_DIR/assets ."
echo "  - Assets moved successfully!"

echo "  - Moving views..."
ssh root@$REMOTE_ADDR "cd $REMOTE_APP_DIR; \
    cp -R $REMOTE_SRC_DIR/views ."
echo "  - Views moved successfully!"

echo "  - Moving Caddyfile..."
ssh root@$REMOTE_ADDR "cd $REMOTE_APP_DIR; \
    cp -R $REMOTE_SRC_DIR/Caddyfile ."
echo "  - Caddyfile moved successfully!"

echo "  - Restarting the server..."
ssh root@$REMOTE_ADDR "sudo service lenslocked.com restart"
echo "  - Server restarted successfully!"

echo "  - Restarting Caddy server..."
ssh root@$REMOTE_ADDR "sudo service caddy restart"
echo "  - Caddy server restarted successfully!"

echo "==== Done releasing lenslocked.com ===="
