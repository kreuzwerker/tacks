.PHONY: install

install:
	go get -u github.com/awslabs/aws-sdk-go/aws
	go get -u github.com/Sirupsen/logrus
	go get -u github.com/spf13/cobra
	go get -u github.com/stretchr/testify/assert
	go get -u gopkg.in/yaml.v2
