package main

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/k-komarov/passbase-currency-api/internal/pkg/constants"
	"github.com/k-komarov/passbase-currency-api/internal/pkg/generated"
	"github.com/k-komarov/passbase-currency-api/internal/pkg/graph"
	"github.com/k-komarov/passbase-currency-api/internal/pkg/middlewares"
	"github.com/k-komarov/passbase-currency-api/internal/pkg/model"
	"github.com/k-komarov/passbase-currency-api/pkg/fixer_api_client"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"runtime/debug"
)

const defaultPort = "3000"
const fixerAccessKey = "3ce51fc854d40824c79825222c6ca953"

var projects = map[string]*model.Project{}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	//logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)

	fixerClient := fixer_api_client.NewClient("http://data.fixer.io/api", fixerAccessKey)

	gqlHandler := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	gqlHandler.AddTransport(transport.POST{})
	gqlHandler.Use(extension.Introspection{})
	gqlHandler.SetRecoverFunc(recoverFunc)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))

	h := middlewares.WithAuthorization(projects)(gqlHandler)
	h = middlewares.WithContextValues(map[interface{}]interface{}{
		constants.CTX_FIXER_CLIENT: fixerClient,
		constants.CTX_PROJECTS:     projects,
	})(h)
	http.Handle("/query", h)

	logrus.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	logrus.Fatal(http.ListenAndServe(":"+port, nil))
}

func recoverFunc(ctx context.Context, err interface{}) (userMessage error) {
	logrus.
		WithField("stack_trace", string(debug.Stack())).
		WithField("query", graphql.GetOperationContext(ctx).RawQuery).
		WithField("variables", graphql.GetOperationContext(ctx).Variables).
		Errorf("Panic: %+v", err)

	return errors.New(http.StatusText(http.StatusInternalServerError))
}
