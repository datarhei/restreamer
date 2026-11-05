package handler

import (
	"net/http"
	"net/http/pprof"

	"github.com/labstack/echo/v4"
)

// The ProfilingHandler type provides a function to register the porfiling endpoints
type ProfilingHandler struct{}

// NewProfiling returns a new Profiling type
func NewProfiling() *ProfilingHandler {
	return &ProfilingHandler{}
}

// Register registers the different golang profiling endpoinds with a router
// @Summary Retrieve profiling data from the application
// @Description Retrieve profiling data from the application
// @ID profiling
// @Produce text/html
// @Success 200 {string} string
// @Failure 404 {string} string
// @Router /profiling [get]
func (p *ProfilingHandler) Register(r *echo.Group) {
	r.GET("/", p.handler(pprof.Index))
	r.GET("/cmdline", p.handler(pprof.Cmdline))
	r.GET("/profile", p.handler(pprof.Profile))
	r.POST("/symbol", p.handler(pprof.Symbol))
	r.GET("/symbol", p.handler(pprof.Symbol))
	r.GET("/trace", p.handler(pprof.Trace))
	r.GET("/allocs", p.handler(pprof.Handler("allocs").ServeHTTP))
	r.GET("/block", p.handler(pprof.Handler("block").ServeHTTP))
	r.GET("/goroutine", p.handler(pprof.Handler("goroutine").ServeHTTP))
	r.GET("/heap", p.handler(pprof.Handler("heap").ServeHTTP))
	r.GET("/mutex", p.handler(pprof.Handler("mutex").ServeHTTP))
	r.GET("/threadcreate", p.handler(pprof.Handler("threadcreate").ServeHTTP))
}

func (p *ProfilingHandler) handler(h http.HandlerFunc) echo.HandlerFunc {
	handler := http.HandlerFunc(h)

	return func(c echo.Context) error {
		handler.ServeHTTP(c.Response(), c.Request())

		return nil
	}
}
