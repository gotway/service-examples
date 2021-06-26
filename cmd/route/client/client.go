package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/gotway/gotway/pkg/log"
	"github.com/gotway/gotway/pkg/tlstest"
	route "github.com/gotway/service-examples/pkg/route"
)

var (
	server = flag.String(
		"server",
		"localhost:11000",
		"The server address in the format of host:port",
	)
	timeoutSeconds = flag.Int("timeout", 10, "Request timeout in seconds")
	timeout        = time.Duration(*timeoutSeconds) * time.Second
	tls            = flag.Bool("tls", false, "Enable TLS")
	tlsCa          = flag.String(
		"tls_ca",
		tlstest.CA(),
		"Certificate authority file for TLS",
	)
	tlsServerHost = flag.String(
		"tls_server_host",
		tlstest.Server(),
		"The server name used to verify the hostname returned by the TLS handshake",
	)
	ctx = context.Background()
)

func main() {
	flag.Parse()

	logger := log.NewLogger(log.Fields{
		"service": "route",
	}, "local", "debug", os.Stdout)
	logger.Info("starting...")

	opts := route.ClientOptions{
		Timeout: timeout,
		TLS: route.TLSOptions{
			Enabled:    *tls,
			CA:         *tlsCa,
			ServerHost: *tlsServerHost,
		},
	}
	logger.Debug(opts)

	logger.Infof("connecting to server at %s...", *server)
	client, err := route.NewClient(*server, opts)
	if err != nil {
		logger.Error(err)
	}
	defer client.Close()

	res, err := client.HealthCheck(ctx)
	if err != nil {
		logger.Error("health check failed ", err)
	}
	logger.Info("health check ", res)

	feature, err := client.GetFeature(ctx, validPoint)
	if err != nil {
		logger.Error("get feature failed ", err)
	}
	logger.Info("get feature ", feature)

	features, err := client.ListFeatures(ctx, rect)
	if err != nil {
		logger.Error("list features failed ", err)
	}
	logger.Info("list features ", features)

	summary, err := client.RecordRoute(ctx)
	if err != nil {
		logger.Error("record route failed ", err)
	}
	logger.Info("record route ", summary)

	notes, err := client.RouteChat(ctx)
	if err != nil {
		logger.Error("route chat failed ", err)
	}
	logger.Info("route chat ", notes)
}
