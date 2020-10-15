default:
	go build -o t .
	cp ./t ~/go/bin/
clean:
	rm t
	rm tuesday