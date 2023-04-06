package remy

import "github.com/wrapped-owls/goremy-di/remy/pkg/injopts"

func cacheOptsFromConfig(conf Config) (options injopts.CacheConfOption) {
	if conf.CanOverride {
		options |= injopts.CacheOptAllowOverride
	}

	if conf.DuckTypeElements {
		options |= injopts.CacheOptReturnAll
	}

	return
}
