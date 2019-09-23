package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/buoyantio/emojivoto/emojivoto-emoji-svc/api"
	"github.com/buoyantio/emojivoto/emojivoto-emoji-svc/emoji"
	"google.golang.org/grpc"
	"contrib.go.opencensus.io/exporter/ocagent"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/trace"
)

var (
	grpcPort    = os.Getenv("GRPC_PORT")
	ocagentHost = os.Getenv("OC_AGENT_HOST")
)

func main() {

	if grpcPort == "" {
		log.Fatalf("GRPC_PORT (currently [%s]) environment variable must me set to run the server.", grpcPort)
	}

	oce, err := ocagent.NewExporter(
		ocagent.WithInsecure(),
		ocagent.WithReconnectionPeriod(5 * time.Second),
		ocagent.WithAddress(ocagentHost),
		ocagent.WithServiceName("emoji"))
	if err != nil {
		log.Fatalf("Failed to create ocagent-exporter: %v", err)
	}
	trace.RegisterExporter(oce)

	allEmoji := emoji.NewAllEmoji()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{}))
	api.NewGrpServer(grpcServer, allEmoji)
	log.Printf("Starting grpc server on GRPC_PORT=[%s]", grpcPort)
	grpcServer.Serve(lis)
}
