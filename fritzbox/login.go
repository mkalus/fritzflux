package fritzbox

import (
	"github.com/bpicode/fritzctl/config"
	"github.com/bpicode/fritzctl/fritz"
	"net/url"
)

// login for home auto
func LoginHomeAuto(uri string, user string, password string) (fritz.HomeAuto, error) {
	myUrl, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	// open fritz
	h := fritz.NewHomeAuto(
		fritz.SkipTLSVerify(),
		fritz.URL(myUrl),
		fritz.Credentials(user, password),
	)

	err = h.Login()
	if err != nil {
		return nil, err
	}

	return h, nil
}

// login for fritzbox
func LoginFritzbox(uri string, user string, password string) (*fritz.Client, error) {
	myUrl, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	client := fritz.NewClientFromConfig(&config.Config{
		Net: &config.Net{
			Protocol: myUrl.Scheme,
			Host:     myUrl.Host,
			Port:     "",
		},
		Login: &config.Login{
			LoginURL: "/login_sid.lua",
			Username: user,
			Password: password,
		},
		Pki: &config.Pki{
			SkipTLSVerify:   true,
			CertificateFile: "",
		},
	})

	err = client.Login()
	if err != nil {
		return nil, err
	}

	return client, nil
}
