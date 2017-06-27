#!/usr/bin/env bash

start(){
    echo $"Start..."
    nohup ./mailservice-linux &
}

run(){
    echo $"Run..."
    ./mailservice-linux
}

stop(){
    echo $"Stop..."
    ps aux | grep mailservice-linux | grep -v grep | awk '{print $2}' | xargs kill -9
}

build_linux(){
    echo $"Build for linux..."
    GOOS=linux GOARCH=amd64 go build -o ./mailservice-linux
}

case "$1" in
   start)
        start
        exit 1
        ;;
   stop)
        stop
        exit 1
        ;;
   restart)
        echo $"Restart..."
        build_linux
        stop
        start
        exit 1
        ;;
   build-linux)
        build_linux
        exit 1
        ;;
   publish)
        build_linux
        git add mailservice-linux
        git commit -m "新版本发布"
        git pull
        git push
        exit 1
        ;;
   *)
        echo $"Usage: $0 {start|run|stop|restart|build-linux|publish}"
        exit 1
        ;;
esac