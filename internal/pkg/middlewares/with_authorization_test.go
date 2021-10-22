package middlewares

import (
	"github.com/k-komarov/passbase-currency-api/internal/pkg/constants"
	"github.com/k-komarov/passbase-currency-api/internal/pkg/model"
	"net/http"
	"testing"
)

type withAuthorizationTestHandler struct {
	result bool
}

func (n *withAuthorizationTestHandler) ServeHTTP(_ http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(constants.CTX_PROJECT).(*model.Project)
	n.result = ok
}

func TestWithAuthorization(t *testing.T) {
	type args struct {
		projects map[string]*model.Project
		handler  http.Handler
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Without Authorization header",
			args: struct {
				projects map[string]*model.Project
				handler  http.Handler
			}{
				projects: nil,
				handler:  &withAuthorizationTestHandler{},
			},
			want: false,
		},
		{
			name: "With Authorization header",
			args: struct {
				projects map[string]*model.Project
				handler  http.Handler
			}{
				projects: map[string]*model.Project{
					"access-key": {},
				},
				handler: &withAuthorizationTestHandler{},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("", "", nil)
			if tt.want {
				req.Header.Set("Authorization", "Bearer access-key")
			}

			th := &withAuthorizationTestHandler{}
			WithAuthorization(tt.args.projects)(th).ServeHTTP(nil, req)

			if got := th.result; tt.want != got {
				t.Errorf("TestWithAuthorization() got = %v, want %v", got, tt.want)
			}
		})
	}
}
