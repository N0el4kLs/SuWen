DIR=data
$(shell mkdir -p ${DIR})

# Go build flags
LDFLAGS=-ldflags "-s -w"

default:
	go build ${LDFLAGS} -o ${DIR}/SuWen main.go

# Compile Server - Windows x64
windows:
	export GOOS=windows;export GOARCH=amd64;go build ${LDFLAGS} -o ${DIR}/SuWen.exe main.go

# Compile Server - Linux x64
linux:
	export GOOS=linux;export GOARCH=amd64;go build ${LDFLAGS} -o ${DIR}/SuWen-linux main.go

# Compile Server - Darwin x64
darwin:
	export GOOS=darwin;export GOARCH=amd64;go build ${LDFLAGS} -o ${DIR}/SuWen-darwin main.go

# clean
clean:
	rm -rf ${DIR}