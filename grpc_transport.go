package lib

import (
	"context"
	"errors"
    "net"
	"io"
    "runtime/debug"

    grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
    grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
    grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
    grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
    grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
    log "github.com/sirupsen/logrus"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
)

// GRPCTransport properties
type GRPCTransport struct {
    Logger            *log.Entry
    Addr              string
    Server            *grpc.Server
    Closers           []io.Closer
}

// GRPCTransportOptions hold options for server
type GRPCTransportOptions struct {
	Addr string
}

// init prevents bootstrap race message
func init() {
    logger := log.WithField("module", "GRPCTransport")
    grpc_logrus.ReplaceGrpcLogger(logger)
}

// NewGRPCTransport
//func NewGRPCTransport(addr string, logPayload bool, decider grpc_logging.ServerPayloadLoggingDecider) (*GRPCTransport, error) {
func NewGRPCTransport(opts GRPCTransportOptions) (*GRPCTransport, error) {
    logger := log.WithField("module", "GRPCTransport")
    //l := logger.WithField("function", "NewGRPCTransport")

    grpcTransport := GRPCTransport{
        Logger:     logger,
        Addr:       opts.Addr,
        Closers:    []io.Closer{},
    }

    payloadLoggingDecider := grpcTransport.defaultPayloadLoggingDecider()

    logOpts := []grpc_logrus.Option{
            grpc_logrus.WithLevels(grpc_logrus.DefaultCodeToLevel),
    }

    recoverOpts := []grpc_recovery.Option{
        grpc_recovery.WithRecoveryHandlerContext(grpcTransport.recover()),
    }

    server := grpc.NewServer(
        grpc_middleware.WithUnaryServerChain(
            grpc_ctxtags.UnaryServerInterceptor(
                    grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor),
            ),
            grpc_logrus.UnaryServerInterceptor(logger, logOpts...),
            grpc_logrus.PayloadUnaryServerInterceptor(logger, payloadLoggingDecider),
            grpc_recovery.UnaryServerInterceptor(recoverOpts...),
        ),
        grpc_middleware.WithStreamServerChain(
            grpc_ctxtags.StreamServerInterceptor(
                    grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor),
            ),
            grpc_logrus.StreamServerInterceptor(logger, logOpts...),
            //grpc_logrus.PayloadStreamServerInterceptor(logger, payloadLoggingDecider),
            grpc_recovery.StreamServerInterceptor(recoverOpts...),
        ),
    )

    reflection.Register(server)

    grpcTransport.Server = server

    return &grpcTransport, nil
}

// defaultPayloadLoggingDecider
func (t *GRPCTransport) defaultPayloadLoggingDecider() grpc_logging.ServerPayloadLoggingDecider {
    return func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
        return log.GetLevel() > log.InfoLevel
    }
}

// GRPCServer
func (t *GRPCTransport) GRPCServer() *grpc.Server {
    return t.Server
}

// AddCloser
func (t *GRPCTransport) AddCloser(c io.Closer) {
    t.Closers = append(t.Closers, c)
}

// Run
func (t *GRPCTransport) Run() error {
    l := t.Logger.WithField("function", "Run")

    lis, err := net.Listen("tcp", t.Addr)
    if err != nil {
            l.Fatal(err)
    }
    l.WithField("addr", t.Addr).Info("listening")

    return t.Server.Serve(lis)
}

// Close
func (t *GRPCTransport) Close() error {
        l := t.Logger.WithField("function", "Close")

        if t.Server != nil {
            t.Server.GracefulStop()
        }

        for _, v := range t.Closers {
            if err := v.Close(); err != nil {
                l.WithField("err", err).Error()
            }
        }

        return nil
}

// recover
func (t *GRPCTransport) recover() grpc_recovery.RecoveryHandlerFuncContext {
    return func(ctx context.Context, p interface{}) error {
        //if log.GetLevel() > log.InfoLevel {
        //      debug.PrintStack()
		//}

		log.WithFields(log.Fields{
            "module":           "RECOVERY",
            "function":         "recover",
            "ctx":              ctx,
            "panic":            p,
		}).Error(debug.Stack())

		return errors.New("recovery error")
    }
}
