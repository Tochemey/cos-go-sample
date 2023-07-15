package cmd

import (
	"context"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tochemey/cos-go-sample/app/dbwriter"
	"github.com/tochemey/cos-go-sample/app/grpconfig"
	"github.com/tochemey/cos-go-sample/app/storage"
	gopack "github.com/tochemey/gopack/grpc"
	"github.com/tochemey/gopack/log/zapl"
)

// dbWriterCmd represents the dbwriter command
var dbWriterCmd = &cobra.Command{
	Use:   "dbwriter",
	Short: "Run the read side db writer",
	Run: func(cmd *cobra.Command, args []string) {
		// create the base context
		ctx := context.Background()
		// load the grpc config
		config := grpconfig.LoadConfig()
		// get the dataStore
		dataStore := storage.New(ctx)
		// create the service
		service, err := dbwriter.NewService(dataStore)
		// log the error in case there is one and panic
		if err != nil {
			zapl.Panic(errors.Wrap(err, "failed to create db writer service"))
		}
		// create the grpc server
		grpcServer, err := gopack.
			NewServerBuilderFromConfig(config).
			WithService(service).
			WithShutdownHook(dataStore.Shutdown(ctx)).
			Build()
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

func init() {
	rootCmd.AddCommand(dbWriterCmd)
}
