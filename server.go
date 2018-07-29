package fstorage

import (
	"fmt"
	"net/http"

	"github.com/dzen-it/fstorage/api"
	"github.com/dzen-it/fstorage/middlewares"
	"github.com/dzen-it/fstorage/storage"
	"github.com/dzen-it/fstorage/utils/log"

	"github.com/go-chi/chi"
)

const (
	defaultHTTPAddr            = ":8080"
	defaultRPS         float64 = 10             // 10 request per second per one IP
	defaultCS          int     = 3              // 3 concurent sessions
	defaultMaxFilesize int64   = 10 << 20 * 100 // 100Mb
)

var (
	ErrConnectToRedis = fmt.Errorf("Error connect to redis")
)

// Server is instance of http server
type Server struct {
	router      chi.Router
	middlewares []func(next http.Handler) http.Handler
	filestorage storage.Storage
	MaxFilesize int64
	RPS         float64
	CS          int // concurent sessions
}

func NewServer(filestorage storage.Storage) *Server {
	srv := new(Server)

	srv.filestorage = filestorage
	srv.RPS = defaultRPS
	srv.CS = defaultCS
	srv.MaxFilesize = defaultMaxFilesize

	return srv
}

// Start runs the server
func (s *Server) Start(listenAddr string) (err error) {
	handlerAPI := api.NewAPI(s.filestorage, s.MaxFilesize)
	t := middlewares.NewThrottler(s.RPS, s.CS)

	s.router = chi.NewRouter()
	s.router.Use(middlewares.MiddlewareThrottling(t))
	s.router.Use(s.middlewares...)

	s.router.With(api.MiddlewareValidateFilename, api.MiddlewareValidateHeadersControl).
		Put("/files/{filename}", handlerAPI.UploadFile())

	s.router.With(api.MiddlewareValidateHash).
		Get("/files/{hash}", handlerAPI.DownloadFile())

	s.router.With(api.MiddlewareValidateHash).
		Delete("/files/{hash}", handlerAPI.DeleteFile())

	if len(listenAddr) == 0 {
		listenAddr = defaultHTTPAddr
	}

	httpServer := &http.Server{
		Addr:    listenAddr,
		Handler: s.router,
	}

	log.Infow("Start File Storage listening", "address", httpServer.Addr)
	return httpServer.ListenAndServe() // Blocks!
}

func (s *Server) MountHandler(pattern string, handler http.Handler) {
	s.router.Mount(pattern, handler)
}

func (s *Server) AddMiddleware(m ...func(next http.Handler) http.Handler) {
	s.middlewares = append(s.middlewares, m...)
}
