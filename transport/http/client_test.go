package http

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/Barty-Uruk/kfmstest/configs"
)

func TestClient_Request(t *testing.T) {
	type fields struct {
		JwtToken     string
		RequestToken string
		ObjectID     string
		Cert         string
		config       configs.Auth
	}
	type args struct {
		method string
		url    string
		body   []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Response
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Client{
				JwtToken:     tt.fields.JwtToken,
				RequestToken: tt.fields.RequestToken,
				ObjectID:     tt.fields.ObjectID,
				Cert:         tt.fields.Cert,
				config:       tt.fields.config,
			}
			got, err := a.Request(tt.args.method, tt.args.url, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Request() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.Request() = %v, want %v", got, tt.want)
			}
		})
	}
}
