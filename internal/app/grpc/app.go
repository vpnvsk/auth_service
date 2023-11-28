package grpc

import (
	"auth_service/internal/handler/auth"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log  *slog.Logger
	gRPC *grpc.Server
	port int
}

func New(log *slog.Logger, port int) *App {
	gRPC := grpc.NewServer()
	auth.Register(gRPC)
	return &App{log: log, gRPC: gRPC, port: port}
}
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}
func (a *App) Run() error {
	const op = "app.handler.Run"
	log := a.log.With(slog.String("op", op))
	l, err := net.Listen("tcp", fmt.Sprintf("%d", a.port))
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
