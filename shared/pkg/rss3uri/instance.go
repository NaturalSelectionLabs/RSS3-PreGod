package rss3uri

import (
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
)

type Instance interface {
	fmt.Stringer

	GetPrefix() string
	GetIdentity() string
	GetSuffix() string
	UriString() string
}

var (
	_ Instance = PlatformInstance{}
	_ Instance = NetworkInstance{}
)

type PlatformInstance struct {
	Prefix   constants.PrefixName     `json:"prefix"`
	Identity string                   `json:"identity"`
	Platform constants.PlatformSymbol `json:"platform"`
}

func (p PlatformInstance) GetPrefix() string {
	return string(p.Prefix)
}

func (p PlatformInstance) GetIdentity() string {
	return p.Identity
}

func (p PlatformInstance) GetSuffix() string {
	return string(p.Platform)
}

func (p PlatformInstance) String() string {
	return fmt.Sprintf("%s:%s@%s", p.Prefix, p.Identity, p.Platform)
}

func (n PlatformInstance) UriString() string {
	return fmt.Sprintf("%s://%s", Scheme, n.String())
}

type NetworkInstance struct {
	Prefix   constants.PrefixName    `json:"prefix"`
	Identity string                  `json:"identity"`
	Network  constants.NetworkSymbol `json:"network"`
}

func (n NetworkInstance) GetPrefix() string {
	return string(n.Prefix)
}

func (n NetworkInstance) GetIdentity() string {
	return n.Identity
}

func (n NetworkInstance) GetSuffix() string {
	return string(n.Network)
}

func (n NetworkInstance) String() string {
	return fmt.Sprintf("%s:%s@%s", n.Prefix, n.Identity, n.Network)
}

func (n NetworkInstance) UriString() string {
	return fmt.Sprintf("%s://%s", Scheme, n.String())
}

func NewAccountInstance(identity string, platform constants.PlatformSymbol) Instance {
	r, err := NewInstance("account", identity, string(platform))
	if err != nil {
		logger.Errorf("Error when creating account instance: %s", err)

		return PlatformInstance{
			Prefix:   constants.PrefixNameAccount,
			Identity: identity,
			Platform: platform,
		}
	}

	return r
}

func NewNoteInstance(identity string, network constants.NetworkSymbol) Instance {
	r, err := NewInstance("note", identity, string(network))
	if err != nil {
		logger.Errorf("Error when creating note instance: %s", err)

		return NetworkInstance{
			Prefix:   constants.PrefixNameNote,
			Identity: identity,
			Network:  network,
		}
	}

	return r
}

func NewAssetInstance(identity string, network constants.NetworkSymbol) Instance {
	r, err := NewInstance("asset", identity, string(network))
	if err != nil {
		logger.Errorf("Error when creating asset instance: %s", err)

		return NetworkInstance{
			Prefix:   constants.PrefixNameAsset,
			Identity: identity,
			Network:  network,
		}
	}

	return r
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
