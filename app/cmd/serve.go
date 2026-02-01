package cmd

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	gopack "github.com/tochemey/gopack/grpc"

	"github.com/tochemey/cos-go-sample/app/cos"
	"github.com/tochemey/cos-go-sample/app/log"
	"github.com/tochemey/cos-go-sample/app/service"
	"github.com/tochemey/cos-go-sample/app/subscription"
)

// serveCmd represents the runApi command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the accounts api service",
	Run: func(cmd *cobra.Command, _ []string) {
		// create the base context
		ctx := cmd.Context()
		// load the service config
		config := service.LoadConfig()
		// create the cos client
		cosClient, err := cos.NewClient(config.CosHost, config.CosPort)
		// handle the error
		if err != nil {
			log.Fatal(errors.Wrap(err, "failed to create the CoS client"))
		}

		// create an instance of the apis service
		apisService := service.NewService(cosClient)
		// get the grpc config
		grpcConfig := config.GRPCConfig.GetGrpcConfig()

		// create subscription handler and manager for CoS event streaming
		subHandler := subscription.NewSubscriptionHandler(nil)
		subManager, err := subscription.NewManager(config.CosHost, config.CosPort, subHandler)
		if err != nil {
			log.Fatal(errors.Wrap(err, "failed to create subscription manager"))
		}

		if err := subManager.Start(ctx); err != nil {
			log.Fatal(errors.Wrap(err, "failed to start CoS subscription"))
		}

		log.Info("CoS subscription started: subscribeAll")

		// create the grpc server with shutdown hook to unsubscribe on stop
		grpcServer, err := gopack.
			NewServerBuilderFromConfig(grpcConfig).
			WithService(apisService).
			WithShutdownHook(func(ctx context.Context) error {
				return subManager.Stop(ctx)
			}).
			Build()
		// log the error in case there is one and panic
		if err != nil {
			log.Fatal(errors.Wrap(err, "failed to build a grpc server"))
		}
		// start the service
		if err := grpcServer.Start(ctx); err != nil {
			log.Fatal(errors.Wrap(err, "failed to create a grpc service"))
		}

		log.Infof("accounts service started on (%s)", fmt.Sprintf("%s:%d", grpcConfig.GrpcHost, grpcConfig.GrpcPort))

		// await for termination
		grpcServer.AwaitTermination(ctx)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
