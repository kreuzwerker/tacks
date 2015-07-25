VERSION := 0.1.0

export GITHUB_REPO := tacks
export GITHUB_USER := kreuzwerker
export TOKEN = `cat .token`

FLAGS := "-X main.build `git rev-parse --short HEAD` -X main.version $(VERSION)"

install:
	go get -u github.com/aws/aws-sdk-go/aws
	go get -u github.com/aws/aws-sdk-go/service/cloudformation
	go get -u github.com/fatih/color
	go get -u github.com/olekukonko/tablewriter
	go get -u github.com/Sirupsen/logrus
	go get -u github.com/spf13/cobra
	go get -u github.com/stretchr/testify/assert
	go get -u golang.org/x/crypto/ssh/terminal
	go get -u gopkg.in/yaml.v2
	go get -u github.com/dustin/go-humanize

build: build/darwin build/linux

build/darwin:
	GOOS=darwin go build -o build/darwin/tacks -ldflags $(FLAGS) bin/tacks.go

build/linux:
	GOOS=linux go build -o build/linux/tacks -ldflags $(FLAGS) bin/tacks.go

clean:
	rm -rf build

release: clean build
	git tag $(VERSION) -f && git push --tags -f
	github-release release --tag $(VERSION) -s $(TOKEN)
	github-release upload --tag $(VERSION) -s $(TOKEN) --name tacks-osx --file build/darwin/tacks
	github-release upload --tag $(VERSION) -s $(TOKEN) --name tacks-linux --file build/linux/tacks

retract:
	github-release delete --tag $(VERSION) -s $(TOKEN)

test:
	go test

.PHONY: test
