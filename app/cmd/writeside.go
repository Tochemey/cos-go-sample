package cmd

import (
	"context"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tochemey/cos-go-sample/app/grpconfig"
	"github.com/tochemey/cos-go-sample/app/writeside"
	"github.com/tochemey/cos-go-sample/app/writeside/commands"
	"github.com/tochemey/cos-go-sample/app/writeside/events"
	gopack "github.com/tochemey/gopack/grpc"
	"github.com/tochemey/gopack/log/zapl"
)

// writesideCmd represents the runWriteside command
var writesideCmd = &cobra.Command{
	Use:   "writeside",
	Short: "Run the commands and events handler service",
	Run: func(cmd *cobra.Command, args []string) {
		// create the base context
		ctx := context.Background()
		// load the grpc config
		config := grpconfig.LoadConfig()
		// create the commands dispatcher
		commandsDispatcher := commands.NewDispatcher()
		// create the events dispatcher
		eventsDispatcher := events.NewDispatcher()
		// create the instance of the service
		service := writeside.NewHandlerService(commandsDispatcher, eventsDispatcher)
		// create the grpc server
		grpcServer, err := gopack.
			NewServerBuilderFromConfig(config).
			WithService(service).Build()
		// log the error in case there is one and panic
		if err != nil {
			zapl.Panic(errors.Wrap(err, "failed to build a grpc server"))
		}
		// start the service
		if err := grpcServer.Start(ctx); err != nil {
			zapl.Panic(errors.Wrap(err, "failed to create a grpc service"))
		}
		// await for termination
		grpcServer.AwaitTermination(ctx)
	},
}

func init() { rootCmd.AddCommand(writesideCmd) }
