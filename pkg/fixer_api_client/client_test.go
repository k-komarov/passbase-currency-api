package fixer_api_client

import (
	"context"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

type httpClientMock struct {
	HTTPClient
}

func (httpClient httpClientMock) Do(r *http.Request) (*http.Response, error) {
	if strings.HasSuffix(r.URL.Path, "latest") {
		if strings.Contains(r.URL.String(), "bad-json-url") {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader("this is not a JSON")),
			}, nil
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader("{\n\"success\": true,\n\"timestamp\": 1634659269,\n\"base\": \"EUR\",\n\"date\": \"2021-10-19\",\n\"rates\": {\n\"USD\": 1.163203\n}\n}")),
		}, nil
	}
	return nil, nil
}

func Test_client_GetLatestRate(t *testing.T) {
	type fields struct {
		BaseURL    string
		AccessKey  string
		HttpClient HTTPClient
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *LatestRateResponse
		wantErr bool
	}{
		{
			name: "Can't parse baseUrl", fields: struct {
				BaseURL    string
				AccessKey  string
				HttpClient HTTPClient
			}{
				BaseURL:    "bad url\f",
				AccessKey:  "any",
				HttpClient: &httpClientMock{},
			},
			args: struct {
				ctx context.Context
			}{
				ctx: context.TODO(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Can't decode json response",
			fields: struct {
				BaseURL    string
				AccessKey  string
				HttpClient HTTPClient
			}{
				BaseURL:    "bad-json-url",
				AccessKey:  "any",
				HttpClient: &httpClientMock{},
			},
			args: struct {
				ctx context.Context
			}{
				ctx: context.TODO(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(tt.fields.BaseURL, tt.fields.AccessKey, WithHttpClient(tt.fields.HttpClient))
			got, err := c.GetLatestEURToUSDRate(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLatestEURToUSDRate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLatestEURToUSDRate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
