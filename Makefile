build:
	mkdir -p target
	go get github.com/julienschmidt/httprouter
	go get github.com/lib/pq
	go get launchpad.net/gocheck
	go get github.com/golang/groupcache
	go build -o target/planner .
	go test -v 

refresh-deps:
	go get github.com/julienschmidt/httprouter
	go get -u github.com/lib/pq
	go get -u launchpad.net/gocheck
	go get -u github.com/golang/groupcache

clean:
	rm -rf target

