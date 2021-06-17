package route

import (
	"context"
	"fmt"
	"io"
	"time"

	pb "github.com/gotway/service-examples/pkg/route/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type TLSOptions struct {
	Enabled    bool
	CA         string
	ServerHost string
}

type ClientOptions struct {
	Timeout time.Duration
	TLS     TLSOptions
}

type Client struct {
	conn         *grpc.ClientConn
	healthClient healthpb.HealthClient
	routeClient  pb.RouteClient
	options      ClientOptions
}

func (c *Client) HealthCheck(ctx context.Context) (*healthpb.HealthCheckResponse, error) {
	return c.healthClient.Check(ctx, &healthpb.HealthCheckRequest{})
}

func (c *Client) GetFeature(ctx context.Context, point *pb.Point) (*pb.Feature, error) {
	return c.routeClient.GetFeature(ctx, point)
}

func (c *Client) ListFeatures(ctx context.Context, rect *pb.Rectangle) ([]*pb.Feature, error) {
	stream, err := c.routeClient.ListFeatures(ctx, rect)
	if err != nil {
		return nil, err
	}
	var features []*pb.Feature
	for {
		feature, err := stream.Recv()
		if err == io.EOF {
			return features, nil
		}
		if err != nil {
			return features, err
		}
		features = append(features, feature)
	}
}

func (c *Client) RecordRoute(ctx context.Context) (*pb.RouteSummary, error) {
	stream, err := c.routeClient.RecordRoute(ctx)
	if err != nil {
		return nil, err
	}
	for _, point := range randomPoints() {
		if err := stream.Send(point); err != nil {
			return nil, err
		}
	}
	summary, err := stream.CloseAndRecv()
	if err != nil {
		return nil, err
	}
	return summary, nil
}

func (c *Client) RouteChat(ctx context.Context) ([]*pb.RouteNote, error) {
	stream, err := c.routeClient.RouteChat(ctx)
	if err != nil {
		return nil, err
	}
	done := make(chan struct{})
	var recvNotes []*pb.RouteNote
	go func() {
		for {
			note, err := stream.Recv()
			if err == io.EOF {
				close(done)
				return
			}
			recvNotes = append(recvNotes, note)
		}
	}()
	for _, note := range notes() {
		if err := stream.Send(note); err != nil {
			return nil, err
		}
	}
	stream.CloseSend()
	<-done

	return recvNotes, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func getConn(server string, clientOpts ClientOptions) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTimeout(clientOpts.Timeout),
	}
	if clientOpts.TLS.Enabled {
		creds, err := credentials.NewClientTLSFromFile(clientOpts.TLS.CA, clientOpts.TLS.ServerHost)
		if err != nil {
			return nil, fmt.Errorf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	return grpc.Dial(server, opts...)
}

func NewClient(server string, opts ClientOptions) (*Client, error) {
	conn, err := getConn(server, opts)
	if err != nil {
		return nil, err
	}
	client := pb.NewRouteClient(conn)
	healthClient := healthpb.NewHealthClient(conn)
	return &Client{conn, healthClient, client, opts}, nil
}
