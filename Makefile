normal: clean
	go build

# Generate code for linux
linux: clean
	env GOOS=linux GOARCH=amd64 go build
	zip -9 bank_linux_v1.0.zip bank conf.yaml

# Generate code for windows
windows:
	env GOOS=windows GOARCH=amd64 go build
	mv bank.exe bank_amd64.exe
	zip -9 bank_win_v1.0.zip bank_amd64.exe conf.yaml

# Clean all old files
clean:
	rm -f bank
	rm -f *.exe
