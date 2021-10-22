package model

import (
	"context"
	"github.com/k-komarov/passbase-currency-api/internal/pkg/constants"
)

func ProjectFromContext(ctx context.Context) *Project {
	if p, ok := ctx.Value(constants.CTX_PROJECT).(*Project); ok {
		return p
	}
	return nil
}
