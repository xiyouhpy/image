#!/bin/bash

BIN="image"
BINPID="run.pid"
[[ -d bin/ ]] || mkdir -p bin

start() {
    stop
    sleep 1
    go build -o ./bin/$BIN main.go
    ./bin/$BIN </dev/null &>/dev/null &
    ps aux | grep "/bin/$BIN" | grep -v "grep" | awk '{print $2}' > ./bin/$BINPID
}

stop(){
    [[ -f $BINPID ]] || kill -9 $(cat ./bin/$BINPID)
}

case RUN"$1" in
    RUN)
    start
    echo "elseDone!"
        ;;
    RUNstart)
        start
        echo "Start Done!"
        ;;
    RUNstop)
        stop
        echo "Stop Done!"
        ;;
    RUN*)
        echo "Usage: $0 {start|stop}"
        ;;
esac
