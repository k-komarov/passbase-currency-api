package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type withCtxTestHandler struct {
	ctx context.Context
}

func (n *withCtxTestHandler) ServeHTTP(_ http.ResponseWriter, r *http.Request) {
	n.ctx = r.Context()
}

func TestWithContextValues(t *testing.T) {
	type args struct {
		keyValues map[interface{}]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[interface{}]interface{}
	}{
		{
			name: "Values present in context",
			args: struct{ keyValues map[interface{}]interface{} }{keyValues: map[interface{}]interface{}{"key1": "value1", "key2": "value2"}},
			want: map[interface{}]interface{}{"key1": "value1", "key2": "value2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("", "", nil)
			handler := &withCtxTestHandler{}
			WithContextValues(tt.args.keyValues)(handler).ServeHTTP(httptest.NewRecorder(), req)
			for key, want := range tt.want {
				if got := handler.ctx.Value(key); got != want {
					t.Errorf("TestWithContextValues() got = %v, want %v", got, want)
				}
			}
		})
	}
}
