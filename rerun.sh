#!/bin/sh

# get the current path
CURPATH=`pwd`

mkdir /tmp/planner_logs

inotifywait -m --timefmt '%d/%m/%y %H:%M' --format '%T %w %f' \
 -e moved_to . | while read date time dir file; do

	if [[ $file == *.go ]]; then
		echo $file
		FILECHANGE=${dir}${file}
       	# convert absolute path to relative
       	FILECHANGEREL=`echo "$FILECHANGE" | sed 's_'$CURPATH'/__'`

       	make
       	make test
       	killall planner
       	./target/planner -v=0 -log_dir="/tmp/planner_logs" -config=./scripts/planner.cfg &
	fi
       
done