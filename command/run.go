package command

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	cf "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/kreuzwerker/tacks"
	"github.com/kreuzwerker/tacks/command/term"
)

type Run struct {
	document    *tacks.Document
	Ask         bool
	DryRun      bool
	Environment string
	Filename    string
	Region      string
	Timeout     int
}

func (r *Run) Document() *tacks.Document {
	return r.document
}

func (r *Run) Run() error {

	data, err := r.readData()

	if err != nil {
		return err
	}

	tpl, err := tacks.NewTemplateFromReader(data)

	if err != nil {
		return err
	}

	var action = r.run

	if r.DryRun {
		action = r.runDry
	}

	if err := tpl.Evaluate(r.Environment, action); err != nil {
		return err
	}

	return nil

}

func (r *Run) readData() (io.ReadCloser, error) {

	const null = ""

	context := &tacks.Context{
		Args0:    os.Args[0],
		Filename: r.Filename,
	}

	if err := context.DetectHashbang(); err != nil {
		return nil, err
	}

	if context.Filename == null {
		return nil, errors.New("Neither running in hashbang-mode nor file argument given")
	}

	if data, err := context.Data(); err != nil {
		return nil, err
	} else {
		return data, nil
	}

}

func (r *Run) runDry(d tacks.Document) error {

	stack, err := d.Parse()

	if err != nil {
		return err
	}

	r.document = &d

	fmt.Fprintf(os.Stderr, stack)

	return nil

}

func (r *Run) run(d tacks.Document) error {

	const null = ""

	e := d.Environment

	if e.Ask && !term.Confirm("Are you sure you to deploy to %q? ", r.Environment) {
		return errors.New("no confirmation received")
	}

	stack, err := d.Parse()

	if err != nil {
		return err
	}

	r.document = &d

	config := &aws.Config{
		Region: r.Region,
	}

	// prefer region from tacks stack
	if e.Region != null {
		config.Region = e.Region
	}

	client := cf.New(config)

	switch e.Mode {
	case "create":
		return r.runCreate(client, d, stack)
	case "upsert":
		return r.runUpsert(client, d, stack)
	default:
		return fmt.Errorf(`unknown mode %q, supported are "create" and "upsert"`)
	}

}

func (r *Run) runCreate(client *cf.CloudFormation, d tacks.Document, stack string) error { // TODO: maybe move document & stack into run-state, also environment

	e := d.Environment

	tacks.Logger().Infof("Creating stack %s", e.StackName)

	var (
		capabilities     []*string
		onFailure        = "DO_NOTHING"
		tags             []*cf.Tag
		timeoutInMinutes uint8 = 15
	)

	if d.IsIamCapabilitiesRequired() {
		capabilities = append(capabilities, aws.String("CAPABILITY_IAM"))
	}

	if e.DeleteOnFailure {
		onFailure = "DELETE"
	}

	if e.Timeout > 0 {
		timeoutInMinutes = e.Timeout
	}

	for key, value := range e.Tags {
		tags = append(tags, &cf.Tag{
			Key:   aws.String(key),
			Value: aws.String(value),
		})
	}

	_, err := client.CreateStack(&cf.CreateStackInput{
		Capabilities:     capabilities,
		OnFailure:        aws.String(onFailure),
		StackName:        aws.String(e.StackName),
		Tags:             tags,
		TemplateBody:     aws.String(stack),
		TimeoutInMinutes: aws.Long(int64(timeoutInMinutes)),
	})

	if err != nil {
		return err
	}

	return nil

}

func (r *Run) runUpdate(client *cf.CloudFormation, d tacks.Document, stack string) error {
	return nil
}

func (r *Run) runUpsert(client *cf.CloudFormation, d tacks.Document, stack string) error {
	return nil
}
