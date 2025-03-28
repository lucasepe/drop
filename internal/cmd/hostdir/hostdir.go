package hostdir

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lucasepe/drop/internal/handlers/fileserver"
	"github.com/lucasepe/drop/internal/middleware"
	"github.com/lucasepe/drop/internal/tools"
	"github.com/lucasepe/x/getopt"
)

func Do(args []string, opts []getopt.OptArg) (err error) {
	dir, err := os.Getwd()
	if err != nil {
		dir = "."
	}
	if len(args) > 0 {
		dir = args[0]
	}

	addr := tools.Str(opts, []string{"-a"}, "127.0.0.1:8080")
	cert := tools.Str(opts, []string{"-c"}, "")
	key := tools.Str(opts, []string{"-k"}, "")

	middlewares := []func(http.Handler) http.Handler{
		middleware.AllowedMethods(),
		middleware.Extra(),
	}

	headers, err := loadConfig(os.DirFS(dir), ".headers")
	if err == nil && headers != nil {
		middlewares = append(middlewares, middleware.Headers(headers))
	}

	middlewares = append(middlewares, middleware.Logger())

	users, err := loadConfig(os.DirFS(dir), ".users")
	if err == nil && users != nil {
		middlewares = append(middlewares, middleware.BasicAuth(users))
	}

	server := &http.Server{
		Addr:         addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 50 * time.Second,
		IdleTimeout:  30 * time.Second,
		Handler:      middleware.Chain(fileserver.New(http.Dir(dir)), middlewares...),
	}

	ctx, stop := signal.NotifyContext(context.Background(), []os.Signal{
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	}...)
	defer stop()

	go func() {
		if len(cert) > 0 && len(key) > 0 {
			log.Fatal(server.ListenAndServeTLS(cert, key))
			return
		}

		log.Fatal(server.ListenAndServe())
	}()

	log.Printf("server is ready to handle requests @ %s\n", server.Addr)
	<-ctx.Done()

	stop()
	log.Printf("\rserver is shutting down gracefully, press Ctrl+C again to force\n")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	return server.Shutdown(ctx)
}
