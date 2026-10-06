#!/bin/bash
CUR_DIR="$(dirname "$0")"
FLAG_FILE="./es_restart_request"
cd "$CUR_DIR/Grout" || exit 1

# Apply pending update
if [ -d "../.update" ]; then
    cp -rf ../.update/* ..
    rm -rf ../.update
fi

export CFW=EMUDECK
export LD_LIBRARY_PATH="$PWD/lib${LD_LIBRARY_PATH:+:${LD_LIBRARY_PATH}}"
chmod +x ./grout

./grout

if [ -f "$FLAG_FILE" ]; then
    rm -f "$FLAG_FILE"
    pkill -f es-de
fi

exit 0
