package http

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	log "github.com/Krew-Guru/kas/pkg/logger"
)

type Client struct {
	Cert   string ``
	logger *log.Logger
	client *http.Client
	config Config
}

type Config struct {
	Type        string
	Method      string
	Host        string
	Path        string
	ApiKey      string
	ContentType string
}

type Credentials struct {
	JwtToken     string ``
	RequestToken string ``
	ObjectID     string ``
	UserID       string ``
}

var netTransport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: 5 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 5 * time.Second,
}

func NewClient(logger *log.Logger) *Client {
	return &Client{
		logger: logger,
		client: &http.Client{
			Transport: netTransport,
		},
	}
}

func (a *Client) Request(ctx context.Context, logger *log.Logger, method, url string, body []byte) (context.Context, *http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return ctx, nil, err
	}
	jwtToken := getStringFromContext(ctx, logger, "jwtToken")
	if len(jwtToken) > 0 {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
	}
	objectID := getStringFromContext(ctx, logger, "X-Auth-Uas-Objectid")
	if len(objectID) > 0 {
		req.Header.Set("X-Auth-Uas-Objectid", objectID)
	}
	userid := getStringFromContext(ctx, logger, "X-Auth-Uas-Userid")
	if len(userid) > 0 {
		req.Header.Set("X-Auth-Uas-Userid", userid)
	}
	groupid := getStringFromContext(ctx, logger, "X-Auth-Uas-Groupid")
	if len(groupid) > 0 {
		req.Header.Set("X-Auth-Uas-Groupid", groupid)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := a.client.Do(req)
	if err != nil {
		return ctx, resp, err
	}
	//parse jwt header
	if len(resp.Header.Get("Authorization")) > 0 {
		credentials := strings.Split(resp.Header.Get("Authorization"), " ")
		if len(credentials) == 2 {
			ctx = context.WithValue(ctx, "jwtToken", credentials[1])
		}
	}
	return ctx, resp, err
}

func getStringFromContext(ctx context.Context, logger *log.Logger, name string) string {
	str, ok := ctx.Value(name).(string)
	if !ok {
		return ""
	}
	return str
}
