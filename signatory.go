package tonic

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// A Signatory is used to sign and validate sessions.
type Signatory struct {
	Secret string
	TTL    time.Duration
	Method jwt.SigningMethod
}

// ErrInvalidTTL is returned if a TTL is too short, or negative.
var ErrInvalidTTL = errors.New("invalid TTL")

//nolint:gochecknoglobals // Needs to be global as it's a fallback.
var defaultSecret string

const expiryClaim = "exp"

// Sign a set of claims from the gin context.
func (s *Signatory) Sign(ctx *gin.Context, claims ...string) (string, error) {
	if s.TTL <= 0 {
		return "", fmt.Errorf("%w: %v", ErrInvalidTTL, s.TTL)
	}

	s.Initialise()

	payload := jwt.MapClaims{}

	for _, claim := range claims {
		payload[claim], _ = ctx.Get(claim)
	}

	payload[expiryClaim] = time.Now().Add(s.TTL).Unix()

	token := jwt.NewWithClaims(s.Method, payload)
	tokenString, err := token.SignedString([]byte(s.Secret))

	if err != nil {
		err = fmt.Errorf("failed to sign JWT: %w", err)
	}

	return tokenString, err
}

// Validate the given token, adding the claims to the context if it is valid.
// Returns true if the token is valid, false otherwise.
func (s *Signatory) Validate(ctx *gin.Context, tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(_ *jwt.Token) (any, error) {
		return []byte(s.Secret), nil
	})

	if err != nil {
		return false
	}

	// Not entirely sure how it's possible to get here without MapClaims, and
	// all attempts to stuff odd things into tokens have lead to either this
	// code not failing, or validation failing before we get here. Since we
	// don't want an unexpected panic due to "this shouldn't happen" (famous
	// last words) we'll treat non map claims as being empty and ignore any cast
	// failure.
	//
	//nolint:errcheck // see above.
	claims, _ := token.Claims.(jwt.MapClaims)

	for k, v := range claims {
		ctx.Set(k, v)
	}

	return true
}

// Initialise a Signatory, ensuring a secret is set. If no secret is set then
// the default secret is used. This is a random string generated on first use.
func (s *Signatory) Initialise() {
	if defaultSecret == "" {
		defaultSecret = GenerateSecret(rand.Reader)
	}

	if s.Secret == "" {
		s.Secret = defaultSecret
	}

	if s.Method == nil {
		s.Method = jwt.SigningMethodHS512
	}
}

// GenerateSecret will generate a new secret key. For most uses rand.Reader
// should be used to ensure a cryptographically secure secret. GenerateSecret
// will panic if it fails.
func GenerateSecret(reader io.Reader) string {
	const secretLength = 64

	secret := make([]byte, secretLength)
	_, err := io.ReadFull(reader, secret)

	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(secret)
}
