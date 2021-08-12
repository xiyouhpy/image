#!/bin/bash

BIN="image"
BINPID="run.pid"
[[ -d bin/ ]] || mkdir -p bin

restart() {
    stop
    sleep 1
    start
}

start() {
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
    RUNrestart)
        restart
        echo "Restart Done!"
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
        echo "Usage: $0 {restart|start|stop}"
        ;;
esac