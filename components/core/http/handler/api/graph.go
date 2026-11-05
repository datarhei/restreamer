package api

import (
	"net/http"

	"github.com/datarhei/core/v16/http/graph/graph"
	"github.com/datarhei/core/v16/http/graph/resolver"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
)

type GraphHandler struct {
	resolver          resolver.Resolver
	path              string
	queryHandler      *handler.Server
	playgroundHandler http.HandlerFunc
}

// NewRestream return a new Restream type. You have to provide a valid Restreamer instance.
func NewGraph(resolver resolver.Resolver, path string) *GraphHandler {
	g := &GraphHandler{
		resolver: resolver,
		path:     path,
	}

	g.queryHandler = handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &g.resolver}))
	g.playgroundHandler = playground.Handler("GraphQL", path)

	return g
}

// Query the GraphAPI
// @Summary Query the GraphAPI
// @Description Query the GraphAPI
// @ID graph-query
// @Accept json
// @Produce json
// @Param query body api.GraphQuery true "GraphQL Query"
// @Success 200 {object} api.GraphResponse
// @Failure 400 {object} api.GraphResponse
// @Security ApiKeyAuth
// @Router /api/graph/query [post]
func (g *GraphHandler) Query(c echo.Context) error {
	g.queryHandler.ServeHTTP(c.Response(), c.Request())

	return nil
}

// Load GraphQL playground
// @Summary Load GraphQL playground
// @Description Load GraphQL playground
// @ID graph-playground
// @Produce html
// @Success 200
// @Security ApiKeyAuth
// @Router /api/graph [get]
func (g *GraphHandler) Playground(c echo.Context) error {
	g.playgroundHandler.ServeHTTP(c.Response(), c.Request())

	return nil
}
