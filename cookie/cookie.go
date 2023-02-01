package cookie

import (
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/idomdavis/tonic"
	"bitbucket.org/idomdavis/tonic/config"
	"github.com/gin-gonic/gin"
)

// Name of the dropped cookie.
const Name = "GinAndTonicAuth"

// Drop a cookie with an authentication token. The token will contain the
// authorisation claims taken from the context, and a TLL. The cookie will also
// be set with a TTL, but this is not used for authentication since it can be
// tampered with. It is simply used as a mechanism to allow the browser to tidy
// up expired cookies.
//
// The dropped cookie can be used with the Authenticate middleware handler.
func Drop(c *gin.Context, security config.Security, claims ...string) error {
	signatory := &tonic.Signatory{TTL: security.SessionTTL, Secret: security.Secret}
	signatory.Initialise()

	maxAge := int(security.SessionTTL.Round(time.Second).Seconds())

	if security.SessionTTL <= 0 || maxAge <= 0 {
		return fmt.Errorf("%w: %v", tonic.ErrInvalidTTL, security.SessionTTL)
	}

	token, err := signatory.Sign(c, claims...)

	if err != nil {
		return fmt.Errorf("failed to drop authorisation cookie: %w", err)
	}

	c.SetCookie(Name, token, maxAge, "/", "", true, true)

	return nil
}

// Authenticate a request by checking for a JWT token in a cookie. Authenticate
// uses sensible defaults if no configuration is set, and will use a generated
// secret if none is set. Failure to authenticate will abort the middleware
// chain and either redirect the request to the given URL, or return
// http.StatusUnauthorized if the redirect is blank.
func Authenticate(security config.Security, redirect string) gin.HandlerFunc {
	signatory := &tonic.Signatory{Secret: security.Secret}
	signatory.Initialise()

	return func(c *gin.Context) {
		var authorised bool

		if token, err := c.Cookie(Name); err == nil {
			authorised = signatory.Validate(c, token)
		}

		switch {
		case !authorised && redirect != "":
			c.Abort()
			c.Redirect(http.StatusTemporaryRedirect, redirect)
		case !authorised && redirect == "":
			c.Abort()
			c.String(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		default:
			c.Next()
		}
	}
}
