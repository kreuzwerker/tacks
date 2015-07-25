package main

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/kreuzwerker/tacks"
	"github.com/kreuzwerker/tacks/command"
	"github.com/kreuzwerker/tacks/command/term"
	"github.com/spf13/cobra"
)

const (
	DefaultRegion = "eu-west-1"
	null          = ""
)

var (
	build   string
	version string
	exitF   func(...interface{})
)

func main() {

	var logger = tacks.Logger()
	exitF = logger.Fatal

	var (
		dryRun      bool
		environment string
		region      string
		refresh     time.Duration
		verbose     bool
	)

	root := &cobra.Command{
		Use:   "tacks",
		Short: "Tacks provides executable CloudFormation stacks",
	}

	root.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	root.PersistentPreRun = func(_ *cobra.Command, _ []string) {

		if verbose {
			logger.Level = logrus.DebugLevel
		}

	}

	run := &cobra.Command{
		Use:   "run [tacks stack filename]",
		Short: "Run a tacks stack",
		Run: func(_ *cobra.Command, args []string) {

			var filename string

			if len(args) > 0 {
				filename = args[0]
			}

			run := &command.Run{
				DryRun:      dryRun,
				Environment: environment,
				Filename:    filename,
				Region:      region,
			}

			command.Foreground(run, exitF)

			if !dryRun {

				watch := &command.Watch{
					Stackname: run.Document().Environment.StackName,
					Region:    run.Region,
				}

				command.Background(watch, exitF)

				term.Wait()

			}

		},
	}

	runFlags := run.Flags()

	runFlags.BoolVarP(&dryRun, "dry-run", "d", false, "output stack to stderr instead of sending it to CloudFormation")
	runFlags.StringVarP(&environment, "environment", "e", "", "specify the environment")

	watch := &cobra.Command{
		Use:   "watch [stack-name]",
		Short: "Watch stack events",
		Run: func(_ *cobra.Command, args []string) {

			var stackname string

			if len(args) > 0 {
				stackname = args[0]
			}

			watch := &command.Watch{
				Stackname: stackname,
				Region:    region,
				Refresh:   refresh,
			}

			command.Background(watch, exitF)

			term.Wait()

		},
	}

	for _, value := range []*cobra.Command{run, watch} {

		value.Flags().StringVarP(&region, "region", "r", DefaultRegion, "AWS region")
		value.Flags().DurationVarP(&refresh, "refresh-interval", "i", command.DefaultRefresh, "refresh interval for watching stack events")

	}

	version := &cobra.Command{
		Use:   "version",
		Short: "Print the version information of tacks",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Tacks %s (%s)\n", version, build)
		},
	}

	root.AddCommand(run)
	root.AddCommand(watch)
	root.AddCommand(version)

	if err := root.Execute(); err != nil {
		logger.Fatal(err)
	}

}
