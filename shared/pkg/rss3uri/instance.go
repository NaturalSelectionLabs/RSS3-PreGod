package rss3uri

import (
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
)

type Instance interface {
	fmt.Stringer

	GetPrefix() string
	GetIdentity() string
	GetSuffix() string
}

var (
	_ Instance = &PlatformInstance{}
	_ Instance = &NetworkInstance{}
)

type PlatformInstance struct {
	Prefix   constants.PrefixName     `json:"prefix"`
	Identity string                   `json:"identity"`
	Platform constants.PlatformSymbol `json:"platform"`
}

func (p *PlatformInstance) GetPrefix() string {
	return string(p.Prefix)
}

func (p *PlatformInstance) GetIdentity() string {
	return p.Identity
}

func (p *PlatformInstance) GetSuffix() string {
	return string(p.Platform)
}

func (p *PlatformInstance) String() string {
	return fmt.Sprintf("%s:%s@%s", p.Prefix, p.Identity, p.Platform)
}

type NetworkInstance struct {
	Prefix   constants.PrefixName    `json:"prefix"`
	Identity string                  `json:"identity"`
	Network  constants.NetworkSymbol `json:"network"`
}

func (n *NetworkInstance) GetPrefix() string {
	return string(n.Prefix)
}

func (n *NetworkInstance) GetIdentity() string {
	return n.Identity
}

func (n *NetworkInstance) GetSuffix() string {
	return string(n.Network)
}

func (n *NetworkInstance) String() string {
	return fmt.Sprintf("%s:%s@%s", n.Prefix, n.Identity, n.Network)
}

func NewInstance(prefix, identity, platform string) (Instance, error) {
	if !constants.IsValidPrefix(prefix) {
		return nil, ErrInvalidPrefix
	}

	if identity == "" {
		return nil, ErrInvalidIdentity
	}

	switch prefix := constants.PrefixName(prefix); prefix {
	case constants.PrefixNameAccount:
		if !constants.IsValidPlatformSymbol(platform) {
			return nil, ErrInvalidPlatform
		}

		return &PlatformInstance{
			Prefix:   prefix,
			Identity: identity,
			Platform: constants.PlatformSymbol(platform),
		}, nil
	default:
		if !constants.IsValidNetworkName(platform) {
			return nil, ErrInvalidNetwork
		}

		return &NetworkInstance{
			Prefix:   prefix,
			Identity: identity,
			Network:  constants.NetworkSymbol(platform),
		}, nil
	}
}

func ParseInstance(rawInstance string) (Instance, error) {
	uri, err := Parse(fmt.Sprintf("%s://%s", Scheme, rawInstance))
	if err != nil {
		return nil, err
	}

	return uri.Instance, nil
}
