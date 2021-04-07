package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

func errorResponse(err error) events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 500,
		Body:       err.Error(),
	}
}

var clientID string = os.Getenv("OIDC_CLIENT_ID")
var clientSecret string = os.Getenv("OIDC_CLIENT_SECRET")
var issuer string = os.Getenv("OIDC_ISSUER")
var provider *oidc.Provider
var oidcConfig *oidc.Config
var verifier *oidc.IDTokenVerifier
var config oauth2.Config

func init() {
	ctx := context.Background()
	var err error

	provider, err = oidc.NewProvider(ctx, issuer)
	if err != nil {
		log.Fatal(err)
	}

	oidcConfig = &oidc.Config{
		ClientID: clientID,
	}

	verifier = provider.Verifier(oidcConfig)

	config = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "https://6s2o3b5oql.execute-api.us-east-1.amazonaws.com/auth/cognito/callback",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	code, ok := request.QueryStringParameters["code"]
	if !ok {
		err := errors.New("No code in query string")
		return errorResponse(err), err
	}

	oauth2Token, err := config.Exchange(ctx, code)
	if err != nil {
		err = errors.New("Failed to exchange token: " + err.Error())
		return errorResponse(err), err
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		err = errors.New("No id_token field in oauth2 token.")
		return errorResponse(err), err
	}

	_, err = verifier.Verify(ctx, rawIDToken)

	if err != nil {
		err = errors.New("Failed to verify ID Token: " + err.Error())
		return errorResponse(err), err
	}

	resp := events.APIGatewayV2HTTPResponse{
		StatusCode: 301,
		Body:       rawIDToken,
		Headers: map[string]string{
			"Location":      "/",
			"Pragma":        "no-cache",
			"Cache-Control": "no-cache",
		},
		Cookies: []string{
			fmt.Sprintf("idToken=%s; Path=/", rawIDToken),
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
