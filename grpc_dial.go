package ms

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// GRPCDial is a helper function that creates a grpc.ClientConn with
// logging and prometheus interecepters
func GRPCDial(url string, entry *logrus.Entry) (*grpc.ClientConn, error) {
	if entry == nil {
		entry = logrus.WithField("function", "GRPCDial")
	}

	opts := []grpc_logrus.Option{
		grpc_logrus.WithLevels(grpc_logrus.DefaultClientCodeToLevel),
	}

	return grpc.Dial(
		url,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			grpc_middleware.ChainUnaryClient(
				grpc_logrus.UnaryClientInterceptor(entry, opts...),
				grpc_prometheus.UnaryClientInterceptor,
			),
		),
	)
}
