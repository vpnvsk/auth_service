package grpc

import (
	"auth_service/internal/handler/auth"
	"auth_service/internal/services"
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"log/slog"
	"net"
)

type App struct {
	log  *slog.Logger
	gRPC *grpc.Server
	port int
}

func New(log *slog.Logger, port int, services services.Services) *App {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			//logging.StartCall, logging.FinishCall,
			logging.PayloadReceived, logging.PayloadSent,
		),
		// Add any other option (check functions starting with logging.With).
	}
	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			log.Error("Recovered from panic", slog.Any("panic", p))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}
	gRPC := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...),
	))
	//gRPC := grpc.NewServer()
	auth.New(gRPC, services)
	return &App{log: log, gRPC: gRPC, port: port}
}

func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}
func (a *App) Run() error {
	const op = "internal.app.handler.Run"
	log := a.log.With(slog.String("op", op))
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s, %w", op, err)
	}
	log.Info("server is running", slog.String("on address", l.Addr().String()))
	if err := a.gRPC.Serve(l); err != nil {
		return fmt.Errorf("%s, %w", op, err)
	}
	return err
}
func (a *App) Stop() {
	const op = "app.handler.Stop"
	a.log.With(slog.String("op", op)).Info("shutting down the server...")
	a.gRPC.GracefulStop()
}
