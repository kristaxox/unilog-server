package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/kristaxox/unilog-server/collector"
	"github.com/kristaxox/unilog-server/pb"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	grpcprom "github.com/grpc-ecosystem/go-grpc-prometheus"
	promthttp "github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	grpcPort              = kingpin.Flag("grpc-port", "grpc port").Envar("GRPC_PORT").Default("9123").Int()
	prometheusMetricsPort = kingpin.Flag("prometheus-metrics-port", "Prometheus metrics port").Envar("PROMETHEUS_METRICS_PORT").Default("9124").Int()
)

func main() {
	kingpin.Parse()
	logrus.SetFormatter(&logrus.JSONFormatter{})

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *grpcPort))
	if err != nil {
		logrus.WithError(err).Fatalln("unable to create tcp listener for grpc")
	}

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpcprom.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpcprom.UnaryServerInterceptor),
	)

	pb.RegisterLogCollectorServer(grpcServer, collector.NewServer())

	grpcprom.Register(grpcServer)
	http.Handle("/metrics", promthttp.Handler())
	go func() {
		logrus.Info("starting metrics server")
		logrus.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", *prometheusMetricsPort), nil))
	}()

	logrus.Info("starting grpc server")
	grpcServer.Serve(lis)
}
