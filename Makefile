default:
	go build -o t .
	du -h ./t
	cp ./t ~/go/bin/

release:
	go build -o t -ldflags="-s -w"
	du -h ./t
	cp ./t ~/go/bin/

small: release
	upx ./t
	du -h ./t
	cp ./t ~/go/bin/

clean:
	rm t
	rm tuesday