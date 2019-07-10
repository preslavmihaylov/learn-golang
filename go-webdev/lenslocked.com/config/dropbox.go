package config

type DropboxConfig struct {
	ID          string `json:"id"`
	Secret      string `json:"secret"`
	AuthURL     string `json:"auth_url"`
	TokenURL    string `json:"token_url"`
	RedirectURL string `json:"redirect_url"`
}
