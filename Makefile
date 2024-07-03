clean:
	rm -f db/knowledgebase.ddb db/knowledgebase.ddb.wal
	rm -f kdb
run:
	go run *.go

build:
	go build -o kdb