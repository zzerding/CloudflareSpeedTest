.PHONY: clean
dist: cftest_linux 
	mkdir -p dist
	echo $$(git tag | tail -n 1)>dist/version
cftest_linux cftest_mac cftest_win: *.go
	GOOS=linux GOARCH=amd64  go build  -o dist/cftest_linux .
	GOOS=darwin GOARCH=amd64 go build -o dist/cftest_mac .
	GOOS=windows GOARCH=amd64 go build -o dist/cftest_win .

clean:
	-rm dist/cftest_linux
	-rm dist/cftest_mac
