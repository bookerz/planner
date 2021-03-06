#!/bin/bash

# get the current path
CURPATH=`pwd`

mkdir /tmp/planner_logs

OS=`uname`

function act_on_event {
       file=$1
       if [[ $file == *.go || $file == *.html ]]; then
              echo $file
              FILECHANGE=${dir}${file}
              # convert absolute path to relative
              FILECHANGEREL=`echo "$FILECHANGE" | sed 's_'$CURPATH'/__'`
              make
              make test
              killall planner
              ./target/planner -v=2 -log_dir="/tmp/planner_logs" -config=./scripts/planner.cfg &
              echo ""
              echo "============================================================================="
              echo ""
       fi
}

if [[ $OS == "Linux" ]]; then
       inotifywait -m -r --timefmt '%d/%m/%y %H:%M' --format '%T %w %f' \
              -e moved_to . | while read date time dir file; do
              act_on_event $file
       done
fi

if [[ $OS == "Darwin" ]]; then
       fswatch 
fi
