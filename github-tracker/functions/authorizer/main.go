package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	jose "github.com/go-jose/go-jose/v3" //libreria que permite validar jwt en go de forma sencilla, asi se sabe si el token puede ser o no accesible
)

const GitHubUserAgentPrefix = "GitHub-Hookshot"

func handler(event events.APIGatewayV2CustomAuthorizerV2Request) (events.APIGatewayCustomAuthorizerResponse, error) {
	route := event.RouteArn

	path := event.RequestContext.HTTP.Path

	jsonData, err := json.Marshal(event)
	if err != nil {
		return denyRequest(fmt.Sprintf("error unmarshal event: %s", err.Error()))
	}

	fmt.Println(string(jsonData))

	if path == "/commit" {
		userAgent, ok := event.Headers["user-agent"]
		if ok {
			fmt.Println("user agent: ", userAgent)

			if strings.HasPrefix(userAgent, GitHubUserAgentPrefix) {
				return allowRequest(route)
			}

			return denyRequest("path /commit request not allowed")
		}
	}

	authToken, ok := event.Headers["authorization"] //header vienen en forma de mapa
	if !ok {
		return denyRequest("missing auth0 token")
	}

	tokenString := strings.TrimPrefix(authToken, "Bearer ") //extraemos el pedazo del inicio; bearer
	if tokenString == "" {
		return denyRequest("missing auth0 token")
	}

	return validateToken(tokenString, route) //verifica el token, devuelve la politica
}

func validateToken(authToken string, route string) (events.APIGatewayCustomAuthorizerResponse, error) {
	_, err := jose.ParseSigned(authToken)
	if err != nil {
		return denyRequest(fmt.Sprintf("invalid0 authh0 token: %s", err.Error()))
	}

	return allowRequest(route)
}

func allowRequest(route string) (events.APIGatewayCustomAuthorizerResponse, error) {
	return events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: "user",
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Effect:   "Allow",
					Action:   []string{"execute-api:Invoke"},
					Resource: []string{route},
				},
			},
		},
	}, nil
}

func denyRequest(reason string) (events.APIGatewayCustomAuthorizerResponse, error) {
	fmt.Println(reason)

	return events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: "anonymous",
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Effect:   "Deny",
					Action:   []string{"execute-api:Invoke"},
					Resource: []string{"*"},
				},
			},
		},
	}, fmt.Errorf(reason)
}

func main() {
	lambda.Start(handler) //llamamos la lambda
}
