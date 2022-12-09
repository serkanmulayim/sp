// Code generated by goa v3.10.2, DO NOT EDIT.
//
// portfolio HTTP server
//
// Command:
// $ goa gen sp/design

package server

import (
	"context"
	"net/http"
	portfolio "sp/gen/portfolio"

	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// Server lists the portfolio service endpoint HTTP handlers.
type Server struct {
	Mounts             []*MountPoint
	PortfolioSummary   http.Handler
	GenHTTPOpenapiJSON http.Handler
}

// MountPoint holds information about the mounted endpoints.
type MountPoint struct {
	// Method is the name of the service method served by the mounted HTTP handler.
	Method string
	// Verb is the HTTP method used to match requests to the mounted handler.
	Verb string
	// Pattern is the HTTP request path pattern used to match requests to the
	// mounted handler.
	Pattern string
}

// New instantiates HTTP handlers for all the portfolio service endpoints using
// the provided encoder and decoder. The handlers are mounted on the given mux
// using the HTTP verb and path defined in the design. errhandler is called
// whenever a response fails to be encoded. formatter is used to format errors
// returned by the service methods prior to encoding. Both errhandler and
// formatter are optional and can be nil.
func New(
	e *portfolio.Endpoints,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
	fileSystemGenHTTPOpenapiJSON http.FileSystem,
) *Server {
	if fileSystemGenHTTPOpenapiJSON == nil {
		fileSystemGenHTTPOpenapiJSON = http.Dir(".")
	}
	return &Server{
		Mounts: []*MountPoint{
			{"PortfolioSummary", "GET", "/sp/portfolio/{account}"},
			{"./gen/http/openapi.json", "GET", "/openapi.json"},
		},
		PortfolioSummary:   NewPortfolioSummaryHandler(e.PortfolioSummary, mux, decoder, encoder, errhandler, formatter),
		GenHTTPOpenapiJSON: http.FileServer(fileSystemGenHTTPOpenapiJSON),
	}
}

// Service returns the name of the service served.
func (s *Server) Service() string { return "portfolio" }

// Use wraps the server handlers with the given middleware.
func (s *Server) Use(m func(http.Handler) http.Handler) {
	s.PortfolioSummary = m(s.PortfolioSummary)
}

// MethodNames returns the methods served.
func (s *Server) MethodNames() []string { return portfolio.MethodNames[:] }

// Mount configures the mux to serve the portfolio endpoints.
func Mount(mux goahttp.Muxer, h *Server) {
	MountPortfolioSummaryHandler(mux, h.PortfolioSummary)
	MountGenHTTPOpenapiJSON(mux, goahttp.Replace("", "/./gen/http/openapi.json", h.GenHTTPOpenapiJSON))
}

// Mount configures the mux to serve the portfolio endpoints.
func (s *Server) Mount(mux goahttp.Muxer) {
	Mount(mux, s)
}

// MountPortfolioSummaryHandler configures the mux to serve the "portfolio"
// service "portfolioSummary" endpoint.
func MountPortfolioSummaryHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/sp/portfolio/{account}", f)
}

// NewPortfolioSummaryHandler creates a HTTP handler which loads the HTTP
// request and calls the "portfolio" service "portfolioSummary" endpoint.
func NewPortfolioSummaryHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodePortfolioSummaryRequest(mux, decoder)
		encodeResponse = EncodePortfolioSummaryResponse(encoder)
		encodeError    = EncodePortfolioSummaryError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "portfolioSummary")
		ctx = context.WithValue(ctx, goa.ServiceKey, "portfolio")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountGenHTTPOpenapiJSON configures the mux to serve GET request made to
// "/openapi.json".
func MountGenHTTPOpenapiJSON(mux goahttp.Muxer, h http.Handler) {
	mux.Handle("GET", "/openapi.json", h.ServeHTTP)
}
