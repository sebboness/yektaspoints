package jwt

import (
	"encoding/json"
	"fmt"

	"github.com/MicahParks/keyfunc/v3"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/sebboness/yektaspoints/util/log"
)

var logger = log.Get()

type JwtParser interface {
	// GetJwtClaims parses the given token and returns a map of claims
	GetJwtClaims(tokenValue string) (map[string]interface{}, error)
}

type AwsJwtParser struct {
	JwtParser
	KeyFunc jwtv5.Keyfunc
}

// NewAwsJwtParser Creates a new JWT Parser
func NewAwsJwtParser(region, userPoolId string) (JwtParser, error) {
	jwtSetUrl := fmt.Sprintf("https://cognito-idp.%v.amazonaws.com/%v/.well-known/jwks.json", region, userPoolId)

	logger.Infof("getting jwt set from %v", jwtSetUrl)

	k, err := keyfunc.NewDefault([]string{jwtSetUrl}) // Context is used to end the refresh goroutine.
	if err != nil {
		return nil, fmt.Errorf("failed to create a keyfunc.Keyfunc from the server's url: %w", err)
	}

	return &AwsJwtParser{
		KeyFunc: k.Keyfunc,
	}, nil
}

func (p *AwsJwtParser) GetJwtClaims(tokenValue string) (map[string]interface{}, error) {
	var claims map[string]interface{}

	token, err := jwtv5.Parse(tokenValue, p.KeyFunc)

	if err != nil {
		return claims, fmt.Errorf("error parsing jwt token: %v", err.Error())
	} else {
		tokenJson, _ := json.Marshal(token.Claims)

		if err := json.Unmarshal(tokenJson, &claims); err != nil {
			return claims, fmt.Errorf("error unmarshalling token claims: %v", err.Error())
		}
	}

	return claims, nil
}
