#!/bin/bash
CUR_DIR="$(realpath "$(dirname "$0")")"
FLAG_FILE="./es_restart_request"
cd "$CUR_DIR/Grout" || exit 1

# Apply pending update
if [ -d "../.update" ]; then
    cp -rf ../.update/* ..
    rm -rf ../.update
fi

export CFW=EMUDECK
export LD_LIBRARY_PATH="$CUR_DIR/Grout/lib${LD_LIBRARY_PATH:+:${LD_LIBRARY_PATH}}"
chmod +x ./grout

./grout

# ES-DE does not watch this flag, so nothing else will clear it.
rm -f "$FLAG_FILE"

exit 0
