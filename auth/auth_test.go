package auth

import (
	"reflect"
	"testing"

	"github.com/camptocamp/terraboard/config"
)

func TestSetup_simple(t *testing.T) {
	expected := "/log/me/out"

	c := config.Config{}
	c.Authentication.LogoutURL = expected

	Setup(&c)

	if logoutURL != expected {
		t.Fatalf("Expected %s, got %s", expected, logoutURL)
	}
}

func TestUserInfo(t *testing.T) {
	expected := User{
		Name:      "foo",
		LogoutURL: "/log/me/out",
		AvatarURL: "http://www.gravatar.com/avatar/b48def645758b95537d4424c84d1a9ff",
	}

	u := UserInfo("foo", "foo@example.com")

	if !reflect.DeepEqual(u, expected) {
		t.Fatalf("Expected %v, got %v", expected, u)
	}
}
