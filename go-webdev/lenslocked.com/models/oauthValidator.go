package models

type oauthValidator struct {
	OAuthDB
}

func (ov *oauthValidator) Create(oauth *OAuth) error {
	err := runOAuthValidationFuncs(oauth,
		ov.requireNotNil, ov.requireToken, ov.requireService, ov.requireUserID)
	if err != nil {
		return err
	}

	return ov.OAuthDB.Create(oauth)
}

func (ov *oauthValidator) Delete(id uint) error {
	if id <= 0 {
		return ErrIDInvalid
	}

	return ov.OAuthDB.Delete(id)
}
