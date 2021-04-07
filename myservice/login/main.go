package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

type Response events.APIGatewayV2HTTPResponse

func errorResponse(err error) Response {
	return Response{
		StatusCode: 500,
		Body:       err.Error(),
	}
}

var clientID string = os.Getenv("OIDC_CLIENT_ID")
var clientSecret string = os.Getenv("OIDC_CLIENT_SECRET")
var issuer string = os.Getenv("OIDC_ISSUER")
var provider *oidc.Provider
var config oauth2.Config

func init() {
	ctx := context.Background()
	var err error

	provider, err = oidc.NewProvider(ctx, issuer)
	if err != nil {
		log.Fatal(err)
	}

	config = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "https://6s2o3b5oql.execute-api.us-east-1.amazonaws.com/auth/cognito/callback",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (Response, error) {
	state, err := randString(16)
	if err != nil {
		return errorResponse(err), err
	}
	nonce, err := randString(16)
	if err != nil {
		return errorResponse(err), err
	}

	resp := Response{
		StatusCode: 301,
		Headers: map[string]string{
			"Location":      config.AuthCodeURL(state, oidc.Nonce(nonce)),
			"Pragma":        "no-cache",
			"Cache-Control": "no-cache",
		},
		Cookies: []string{
			fmt.Sprintf("state=%s", state),
			fmt.Sprintf("nonce=%s", nonce),
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
