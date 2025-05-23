package healthchecks

import (
	"crypto"
	"github.com/xenolf/lego/acme"

	"github.com/cloud-gov/cf-cdn-service-broker/config"
)

type User struct {
	Email        string
	Registration *acme.RegistrationResource
	key          crypto.PrivateKey
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetRegistration() *acme.RegistrationResource {
	return u.Registration
}

func (u *User) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func LetsEncrypt(settings config.Settings) error {
	user := &User{key: "cheese"}
	_, err := acme.NewClient("https://acme-v01.api.letsencrypt.org/directory", user, acme.RSA2048)
	return err
}
