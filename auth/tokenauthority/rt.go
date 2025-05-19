package tokenauthority

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

func (t *TokenAuthority) IssueRT() (rt string, expiresAt time.Time, err error) {

	rtBytes := make([]byte, t.rtSizeInBytes)

	_, err = rand.Read(rtBytes)
	if err != nil {
		return "", time.Time{}, errors.New("failed to generate random bytes for RT" + err.Error())
	}

	rt = hex.EncodeToString(rtBytes)
	expiresAt = time.Now().UTC().Add(t.rtLifetime)
	return rt, expiresAt, nil
}
