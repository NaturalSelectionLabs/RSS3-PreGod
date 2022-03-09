package rss3uri

import (
	"errors"
	"fmt"
	"net/url"
)

const (
	Scheme = "rss3"
)

var (
	ErrInvalidScheme   = errors.New("invalid scheme")
	ErrInvalidPrefix   = errors.New("invalid prefix")
	ErrInvalidIdentity = errors.New("invalid identity")
	ErrInvalidPlatform = errors.New("invalid platform")
	ErrInvalidNetwork  = errors.New("invalid network")
)

var (
	_ fmt.Stringer = &URI{}
)

type URI struct {
	Instance Instance `json:"instance"`
}

func (u *URI) String() string {
	value := url.URL{
		Scheme: Scheme,
		User:   url.UserPassword(u.Instance.GetPrefix(), u.Instance.GetIdentity()),
		Path:   u.Instance.GetSuffix(),
	}

	return value.String()
}

func New(instance Instance) *URI {
	return &URI{
		Instance: instance,
	}
}

func Parse(rawURL string) (*URI, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	if u.Scheme != Scheme {
		return nil, ErrInvalidScheme
	}

	prefix := u.User.Username()
	identity, _ := u.User.Password()
	platform := u.Hostname()

	instance, err := NewInstance(prefix, identity, platform)
	if err != nil {
		return nil, err
	}

	return &URI{
		Instance: instance,
	}, nil
}
