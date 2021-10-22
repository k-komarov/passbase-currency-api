package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/k-komarov/passbase-currency-api/internal/pkg/constants"
	"github.com/k-komarov/passbase-currency-api/internal/pkg/generated"
	"github.com/k-komarov/passbase-currency-api/internal/pkg/model"
	"github.com/k-komarov/passbase-currency-api/pkg/fixer_api_client"
	"github.com/sirupsen/logrus"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func (r *mutationResolver) CreateProject(ctx context.Context, project model.ProjectInput) (*model.Project, error) {
	p := &model.Project{
		Name:      project.Name,
		AccessKey: uuid.NewString(),
	}

	projects := ctx.Value(constants.CTX_PROJECTS).(map[string]*model.Project)
	projects[p.AccessKey] = p
	return p, nil
}

func (r *queryResolver) Convert(ctx context.Context, from *model.Symbol, to *model.Symbol, amount float64) (*model.ConversionResult, error) {
	logger := logrus.WithContext(ctx)
	if p := model.ProjectFromContext(ctx); p == nil {
		logger.Warnf("Attempt to access with invalid key")
		return nil, gqlerror.Errorf("Access key is not valid or not defined")
	}
	client, ok := ctx.Value(constants.CTX_FIXER_CLIENT).(fixer_api_client.Client)
	if !ok {
		logger.Error("Fixer client is missing in context")
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}
	//only base=EUR is available in free plan
	resp, err := client.GetLatestEURToUSDRate(ctx)
	if err != nil {
		logger.Errorf("Error retrieving rates from fixer.io: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	s := model.SymbolUsd.String()
	rate, ok := resp.Rates[s]
	if !ok {
		logger.Errorf("Error retrieving USD rate from fixer.io: %v", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	//Default direction EUR -> USD
	var result = amount * rate

	if *from == *to {
		result = 1
		rate = 1
	} else if *from == model.SymbolUsd && *to == model.SymbolEur {
		result = amount / rate
		rate = 1 / rate
	}

	result, _ = strconv.ParseFloat(fmt.Sprintf("%.5f", result), 64)
	rate, _ = strconv.ParseFloat(fmt.Sprintf("%.5f", rate), 64)

	return &model.ConversionResult{
		Timestamp: resp.Timestamp.Time,
		Rate:      rate,
		Result:    result,
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
