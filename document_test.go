package tacks

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func readDocument(a *assert.Assertions, file string) Document {

	var document Document

	out, err := ioutil.ReadFile(file)
	a.NoError(err)

	err = yaml.Unmarshal(out, &document.stack)
	a.NoError(err)
	a.NotEmpty(document.stack)

	return document

}

func TestDocumentIamRequired(t *testing.T) {

	assert := assert.New(t)

	document := readDocument(assert, "test/opsworks-vpc-elb.json")

	types := []string{
		"AWS::EC2::EIP",
		"AWS::EC2::Instance",
		"AWS::EC2::InternetGateway",
		"AWS::EC2::NetworkAcl",
		"AWS::EC2::NetworkAclEntry",
		"AWS::EC2::Route",
		"AWS::EC2::RouteTable",
		"AWS::EC2::SecurityGroup",
		"AWS::EC2::Subnet",
		"AWS::EC2::SubnetNetworkAclAssociation",
		"AWS::EC2::SubnetRouteTableAssociation",
		"AWS::EC2::VPC",
		"AWS::EC2::VPCGatewayAttachment",
		"AWS::ElasticLoadBalancing::LoadBalancer",
		"AWS::IAM::InstanceProfile",
		"AWS::IAM::Role",
		"AWS::OpsWorks::App",
		"AWS::OpsWorks::ElasticLoadBalancerAttachment",
		"AWS::OpsWorks::Instance",
		"AWS::OpsWorks::Layer",
		"AWS::OpsWorks::Stack",
	}

	assert.Equal(types, document.Types())
	assert.Equal(true, document.IsIamCapabilitiesRequired())

}

func TestDocumentIamNotRequired(t *testing.T) {

	assert := assert.New(t)

	document := readDocument(assert, "test/redshift-vpc.json")

	types := []string{
		"AWS::EC2::SecurityGroup",
		"AWS::EC2::Subnet",
		"AWS::EC2::VPC",
		"AWS::Redshift::Cluster",
		"AWS::Redshift::ClusterParameterGroup",
		"AWS::Redshift::ClusterSubnetGroup",
	}

	assert.Equal(types, document.Types())
	assert.Equal(false, document.IsIamCapabilitiesRequired())

}

func TestDocumentParse(t *testing.T) {

	assert := assert.New(t)

	in := `{
  "a": "b",
  "c": [
    0,
    1,
    {
      "d": {
        "e": "f"
      }
    },
    3
  ]
}`

	var document Document

	if err := yaml.Unmarshal([]byte(in), &document.stack); err != nil {
		assert.NoError(err)
	}

	out, err := document.Parse()
	assert.NoError(err)

	assert.Equal(in, out)

}
