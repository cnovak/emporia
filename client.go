package emporia

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	cognitosrp "github.com/alexrudd/cognito-srp/v4"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

const (
	apiBaseURL    = "https://api.emporiaenergy.com"
	cognitoRegion = "us-east-2"
	clientID      = "4qte47jbstod8apnfic0bunmrq"
	userPoolID    = "us-east-2_ghlOXVLi1"
)

// Client represents an Emporia API Client
type Client struct {
	httpClient   *http.Client
	AccessToken  string
	IDToken      string
	RefreshToken string
}

func NewClient(username string, password string) (*Client, error) {
	client := &Client{
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}

	if err := client.Authenticate(username, password); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) get(u *url.URL) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	return c.do(req)
}

func (c *Client) do(req *http.Request) (*http.Response, error) {

	req.Header.Set("authtoken", c.IDToken)

	var resp *http.Response

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Check for 401 Unauthorized
	if resp.StatusCode == http.StatusUnauthorized {
		// TODO: Refresh the token
		return nil, fmt.Errorf("unauthorized: %v", err)
	}

	if resp.StatusCode != 200 {
		// get error from message
		error := struct {
			Message string `json:"message"`
		}{}
		err := json.NewDecoder(resp.Body).Decode(&error)
		if err != nil {
			return nil, fmt.Errorf("cannot parse http error response: %v", err)
		}

		return nil, fmt.Errorf("http error: %v", error.Message)
	}

	return resp, nil
}

func (c *Client) Authenticate(username string, password string) error {
	// configure cognito srp
	csrp, _ := cognitosrp.NewCognitoSRP(username, password, userPoolID, clientID, nil)

	// configure cognito identity provider
	cfg, _ := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(cognitoRegion),
		config.WithCredentialsProvider(aws.AnonymousCredentials{}),
	)
	svc := cip.NewFromConfig(cfg)

	// initiate auth
	resp, err := svc.InitiateAuth(context.Background(), &cip.InitiateAuthInput{
		AuthFlow:       types.AuthFlowTypeUserSrpAuth,
		ClientId:       aws.String(csrp.GetClientId()),
		AuthParameters: csrp.GetAuthParams(),
	})
	if err != nil {
		return err
	}

	// respond to password verifier challenge
	if resp.ChallengeName == types.ChallengeNameTypePasswordVerifier {
		challengeResponses, _ := csrp.PasswordVerifierChallenge(resp.ChallengeParameters, time.Now())

		resp, err := svc.RespondToAuthChallenge(context.Background(), &cip.RespondToAuthChallengeInput{
			ChallengeName:      types.ChallengeNameTypePasswordVerifier,
			ChallengeResponses: challengeResponses,
			ClientId:           aws.String(csrp.GetClientId()),
		})
		if err != nil {
			panic(err)
		}

		c.AccessToken = *resp.AuthenticationResult.AccessToken
		c.IDToken = *resp.AuthenticationResult.IdToken
		c.RefreshToken = *resp.AuthenticationResult.RefreshToken
	}
	return nil
}
