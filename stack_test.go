package tacks

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func readStack(a *assert.Assertions, file string) Stack {

	var stack Stack

	out, err := ioutil.ReadFile(file)
	a.NoError(err)

	fmt.Println(string(out))

	err = yaml.Unmarshal(out, &stack)
	a.NoError(err)

	return stack

}

func TestStackIamRequired(t *testing.T) {

	assert := assert.New(t)

	stack := readStack(assert, "test/opsworks-vpc-elb.json")

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

	assert.Equal(types, stack.Types())
	assert.Equal(true, stack.IsIamCapabilitiesRequired())

}

func TestStackIamNotRequired(t *testing.T) {

	assert := assert.New(t)

	stack := readStack(assert, "test/redshift-vpc.json")

	types := []string{
		"AWS::EC2::SecurityGroup",
		"AWS::EC2::Subnet",
		"AWS::EC2::VPC",
		"AWS::Redshift::Cluster",
		"AWS::Redshift::ClusterParameterGroup",
		"AWS::Redshift::ClusterSubnetGroup",
	}

	assert.Equal(types, stack.Types())
	assert.Equal(false, stack.IsIamCapabilitiesRequired())

}

func TestStackMarshalText(t *testing.T) {

	assert := assert.New(t)

	in := []byte(`{
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
}`)

	var stack Stack

	if err := yaml.Unmarshal(in, &stack); err != nil {
		assert.NoError(err)
	}

	out, err := stack.MarshalText()
	assert.NoError(err)

	assert.Equal(in, out)

}
