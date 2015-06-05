package command

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	cf "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/dustin/go-humanize"
	"github.com/kreuzwerker/tacks/command/term"
	table "github.com/olekukonko/tablewriter"
)

const DefaultRefresh = 5 * time.Second

// Watch watches events for a given stack and summarizes them in a colorized table
type Watch struct {
	Stackname string
	Region    string
	Refresh   time.Duration
}

func (w *Watch) Run() error {

	const null = ""

	if w.Stackname == null {
		return errors.New("no stack name given")
	}

	if w.Refresh == 0 {
		w.Refresh = DefaultRefresh
	}

	name := aws.String(w.Stackname)

	client := cf.New(&aws.Config{
		Region: w.Region,
	})

	for {

		width, height, err := term.Size()

		if err != nil {
			return err
		}

		resp, err := client.DescribeStackEvents(&cf.DescribeStackEventsInput{
			StackName: name,
		})

		if err != nil {
			return err
		}

		tw := table.NewWriter(os.Stdout)
		tw.SetColWidth(width / 4)
		tw.SetColumnSeparator(" ")

		tw.SetHeader([]string{
			"Timestamp",
			"ID Logical",
			"ID Physical",
			"Status",
			"Reason",
		})

		var max = height / 2

		if l := len(resp.StackEvents); max > l {
			max = l
		}

		for _, e := range resp.StackEvents[0:max] {

			var reason string

			if r := e.ResourceStatusReason; r != nil {
				reason = strings.Replace(*r, "\n", " ", -1)
			}

			timestamp := humanize.Time(*e.Timestamp)

			data := []string{
				timestamp,
				*e.LogicalResourceID,
				*e.ResourceType,
				term.CloudFormation.Colorize(*e.ResourceStatus),
				reason,
			}

			tw.Append(data)
			fmt.Print("\033[H\033[2J")
			tw.Render()

		}

		time.Sleep(w.Refresh)

	}

}
