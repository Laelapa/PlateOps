package tokenauthority

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

// IssueRT generates a new refresh token (RT) and its expiration time.
// It creates a random byte slice of size rtSizeInBytes, encodes it to a hex string,
// and calculates the expiration time by adding rtLifetime to the current time.
// Returns the generated refresh token, its expiration time, and an error if any occurred during generation.
// rtSizeInBytes and rtLifetime are properties of the TokenAuthority instance.
//
// Returns:
//   - string: The generated refresh token in hex format.
//   - time.Time: The expiration time of the refresh token in UTC.
//   - error: An error if the random byte generation fails, otherwise nil.
func (t *TokenAuthority) IssueRT() (string, time.Time, error) {

	rtBytes := make([]byte, t.rtSizeInBytes)

	_, err := rand.Read(rtBytes)
	if err != nil {
		return "", time.Time{}, errors.New("failed to generate random bytes for RT" + err.Error())
	}

	rt := hex.EncodeToString(rtBytes)
	expiresAt := time.Now().UTC().Add(t.rtLifetime)
	return rt, expiresAt, nil
}
