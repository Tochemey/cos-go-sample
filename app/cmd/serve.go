package cmd

import (
	"context"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tochemey/cos-go-sample/app/cos"
	"github.com/tochemey/cos-go-sample/app/service"
	gopack "github.com/tochemey/gopack/grpc"
	"github.com/tochemey/gopack/log/zapl"
)

// serveCmd represents the runApi command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the accounts api service",
	Run: func(cmd *cobra.Command, args []string) {
		// create the base context
		ctx := context.Background()
		// load the service config
		config := service.LoadConfig()
		// create the cos client
		cosClient, err := cos.NewClient(ctx, config.CosHost, config.CosPort)
		// handle the error
		if err != nil {
			zapl.Panic(errors.Wrap(err, "failed to create the CoS client"))
		}
		// create an instance of the apis service
		apisService := service.NewService(cosClient)
		// create the grpc server
		grpcServer, err := gopack.
			NewServerBuilderFromConfig(config.GRPCConfig.GetGrpcConfig()).
			WithService(apisService).Build()
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
	rootCmd.AddCommand(serveCmd)
}
