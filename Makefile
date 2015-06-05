.PHONY: install test

install:
	go get -u github.com/awslabs/aws-sdk-go/aws
	go get -u github.com/dustin/go-humanize
	go get -u github.com/fatih/color
	go get -u github.com/olekukonko/tablewriter
	go get -u github.com/Sirupsen/logrus
	go get -u github.com/spf13/cobra
	go get -u github.com/stretchr/testify/assert
	go get -u golang.org/x/crypto/ssh/terminal
	go get -u gopkg.in/yaml.v2

test:
	go test -test.parallel=2 ./...
