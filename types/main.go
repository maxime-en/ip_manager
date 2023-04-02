package types

import (
	"regexp"
)

// Contains all prefixes of a service
type Service struct {
	Prefixes      []*Prefix
	PrefixesByKey map[string]*Prefix
}

// Prefix generator
func NewService() *Service {
	return &Service{
		Prefixes:      make([]*Prefix, 0),
		PrefixesByKey: make(map[string]*Prefix),
	}
}

// Check if key only contains allowed chars and is not longer than 50 chars
func isValidKey(key string) bool {
	// 50 chars max
	if len(key) > 50 {
		return false
	}

	// Allowed chars: [a-z0-9.-]
	match, _ := regexp.MatchString("^[a-z0-9.-]+$", key)
	return match
}
