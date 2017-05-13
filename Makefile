deploy: linux zipem
	
zipem:
	zip -r commaai.zip ./bin/commaai_linux static/

all: linux windows darwin
darwin:
	CGO_ENABLED=0 GOOS=darwin go build -o ./bin/commaai_darwin
linux:
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/commaai_linux
windows:
	CGO_ENABLED=0 GOOS=windows go build -o ./bin/commaai_windows
