package auth

import (
	"crypto/md5"
	"fmt"

	"github.com/camptocamp/terraboard/config"
)

var logoutURL string

// User is an authenticated user
type User struct {
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	LogoutURL string `json:"logout_url"`
}

// Setup sets up authentication
func Setup(c *config.Config) {
	logoutURL = c.Web.LogoutURL
}

// UserInfo returns a User given a name and email
func UserInfo(name, email string) (user User) {
	user = User{
		LogoutURL: logoutURL,
	}

	if email != "" {
		user.Name = name
		user.AvatarURL = fmt.Sprintf("http://www.gravatar.com/avatar/%x", md5.Sum([]byte(email)))
	}

	return
}
