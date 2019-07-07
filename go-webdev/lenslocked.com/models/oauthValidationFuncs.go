package models

type oauthValidationFunc func(oauth *OAuth) error

func runOAuthValidationFuncs(oauth *OAuth, fns ...oauthValidationFunc) error {
	for _, fn := range fns {
		err := fn(oauth)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ov *oauthValidator) requireNotNil(oauth *OAuth) error {
	if oauth == nil {
		return ErrNilModel
	}

	return nil
}

func (ov *oauthValidator) requireUserID(oauth *OAuth) error {
	if oauth.UserID == 0 {
		return ErrUserIDRequired
	}

	return nil
}

func (ov *oauthValidator) requireService(oauth *OAuth) error {
	if oauth.Service == "" {
		return ErrOAuthServiceRequired
	}

	return nil
}

func (ov *oauthValidator) requireToken(oauth *OAuth) error {
	if oauth.Token.AccessToken == "" {
		return ErrOAuthTokenRequired
	}

	return nil
}
