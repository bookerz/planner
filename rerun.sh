#!/bin/sh

# get the current path
CURPATH=`pwd`

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
       	./target/planner -config=./scripts/planner.cfg &
	fi
       
done