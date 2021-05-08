linux-amd64:
	GOOS=linux GOARCH=amd64 go build \
		 -ldflags "-X 'main.commitID="$$(git rev-parse HEAD)"'" \
		 .

darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build \
		 -ldflags "-X 'main.commitID="$$(git rev-parse HEAD)"'" \

windows-amd64:
	GOOS=windows GOARCH=amd64 go build \
		 -ldflags "-X 'main.commitID="$$(git rev-parse HEAD)"'" \
		 .

linux-amd64-release:
	GOOS=linux GOARCH=amd64 go build \
		 -o mdt-linux-amd64-release \
		 -ldflags "-w -s \
		 	-X 'main.commitID="$$(git rev-parse HEAD)"' \
			-X 'main.version="$$(git describe --tags)"'" \
		 .

darwin-amd64-release:
	GOOS=darwin GOARCH=amd64 go build \
		 -o mdt-darwin-amd64-release \
		 -ldflags "-w -s \
		 	-X 'main.commitID="$$(git rev-parse HEAD)"' \
			-X 'main.version="$$(git describe --tags)"'" \
		 .

windows-amd64-release:
	GOOS=windows GOARCH=amd64 go build \
		 -o mdt-windows-amd64-release \
		 -ldflags "-w -s \
		 	-X 'main.commitID="$$(git rev-parse HEAD)"' \
			-X 'main.version="$$(git describe --tags)"'" \
		 .

zipped-releases: linux-amd64-release darwin-amd64-release windows-amd64-release
	# Linux
	mv mdt-linux-amd64-release mdt
	tar -czvf "mdt-$$(git describe --tags)-linux-amd-64.tar.gz" mdt
	rm mdt
	# Darwin
	mv mdt-darwin-amd64-release mdt
	tar -czvf "mdt-$$(git describe --tags)-darwin-amd-64.tar.gz" mdt
	rm mdt
	# Windows
	mv mdt-windows-amd64-release "mdt.exe"
	zip -r "mdt-$$(git describe --tags)-windows-amd-64.zip" "mdt.exe"
	rm "mdt.exe"
