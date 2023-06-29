package jwt

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/domdavis/tonic"
	"github.com/domdavis/tonic/config"
	"github.com/gin-gonic/gin"
)

// ContentType of the returned token.
const ContentType = "application/jwt"

// Header used to hold the JWT token.
const Header = "Authorization"

// Prefix is applied/removed from the front of the JWT token.
const Prefix = "Bearer"

// Sign a response with an authentication token. The token will contain the
// authorisation claims taken from the context, and a TLL. Sign will send the
// response to the client so no further action is required once called.
// Failure to sign the token will result in a 500 error response.
//
// The token can be used with the Authenticate middleware handler.
func Sign(c *gin.Context, security config.Security, claims ...string) {
	signatory := &tonic.Signatory{TTL: security.SessionTTL, Secret: security.Secret}
	signatory.Initialise()

	if token, err := signatory.Sign(c, claims...); err != nil {
		//nolint:errcheck // Gin is handling this for us.
		_ = c.Error(fmt.Errorf("failed to drop authorisation cookie: %w", err))
		c.String(http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError))
	} else {
		c.Data(http.StatusOK, ContentType, []byte(token))
	}
}

// Set the JWT in the headers. A prefix will be added to the token if needed.
func Set(r *http.Request, token string) {
	if !strings.HasPrefix(token, Prefix) {
		token = fmt.Sprintf("%s %s", Prefix, token)
	}

	r.Header.Set(Header, token)
}

// Authenticate a request by checking for a JWT token in the Authorisation
// header. Authenticate uses sensible defaults if no configuration is set, and
// will use a generated secret if none is set. Failure to authenticate will
// abort the middleware chain and return http.StatusUnauthorized.
func Authenticate(security config.Security) gin.HandlerFunc {
	signatory := &tonic.Signatory{Secret: security.Secret}
	signatory.Initialise()

	return func(c *gin.Context) {
		token := c.GetHeader(Header)
		token = strings.TrimPrefix(token, Prefix)
		token = strings.TrimSpace(token)

		if !signatory.Validate(c, token) {
			c.Abort()
			c.String(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		} else {
			c.Next()
		}
	}
}
