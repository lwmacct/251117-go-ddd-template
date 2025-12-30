package urn

import "strings"

// Match reports whether the URN matches the given pattern.
//
// Matching rules:
//   - Each segment (scope, type, identifier) is matched independently
//   - * matches any value for that segment
//   - Scope segment: .* suffix matches sub-scopes (sys.* matches sys.admin)
//
// Examples:
//
//	URN("sys:users:create").Match("*:*:*")            // true - super wildcard
//	URN("sys:users:create").Match("sys:*:*")          // true - scope wildcard
//	URN("sys:users:create").Match("sys:users:*")      // true - type wildcard
//	URN("sys:users:create").Match("sys:users:create") // true - exact match
//	URN("sys.admin:config:update").Match("sys.*:*:*") // true - sub-scope wildcard
func (u URN) Match(pattern URN) bool {
	// Super wildcard
	if pattern == "*" || pattern == "*:*:*" {
		return true
	}

	// Exact match
	if pattern == u {
		return true
	}

	p := pattern.Parse()
	t := u.Parse()

	// Scope matching
	if !matchScope(p.Scope, t.Scope) {
		return false
	}

	// Type matching
	if p.Type != "*" && p.Type != t.Type {
		return false
	}

	// Identifier matching
	if p.Identifier != "*" && p.Identifier != t.Identifier {
		return false
	}

	return true
}

// matchScope matches scope with support for hierarchical wildcards.
//
// Rules:
//   - "*" matches any scope
//   - "sys" exactly matches "sys"
//   - "sys.*" matches "sys", "sys.admin", "sys.readonly", etc.
func matchScope(pattern, scope string) bool {
	// Wildcard matches all
	if pattern == "*" {
		return true
	}

	// Hierarchical wildcard: sys.* matches sys and all sub-scopes
	if prefix, found := strings.CutSuffix(pattern, ".*"); found {
		// Exact match of prefix or starts with prefix.
		return scope == prefix || strings.HasPrefix(scope, prefix+".")
	}

	// Exact match
	return pattern == scope
}

// Match is a convenience function that checks if target matches pattern.
func Match(pattern, target string) bool {
	return URN(target).Match(URN(pattern))
}

// MatchOperation is an alias for Match with operation semantics.
// Provided for clarity when matching operation URNs.
func MatchOperation(pattern, operation string) bool {
	return Match(pattern, operation)
}

// MatchResource is an alias for Match with resource semantics.
// Provided for clarity when matching resource URNs.
func MatchResource(pattern, resource string) bool {
	return Match(pattern, resource)
}
