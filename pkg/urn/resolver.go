package urn

import (
	"maps"
	"strings"
)

// Resolver performs variable substitution in URNs.
//
// Variables are arbitrary strings that get replaced with their values.
// Common convention uses @ prefix (e.g., @me, @org) but any string works.
type Resolver struct {
	vars map[string]string
}

// NewResolver creates a new Resolver with the given variable mappings.
//
// Example:
//
//	r := NewResolver(map[string]string{
//	    "@me":  "123",
//	    "@org": "acme",
//	})
//	r.ResolveString("self:user:@me")  // "self:user:123"
//	r.ResolveString("org.@org:*:*")   // "org.acme:*:*"
func NewResolver(vars map[string]string) *Resolver {
	// Copy the map to prevent external modification
	m := make(map[string]string, len(vars))
	maps.Copy(m, vars)
	return &Resolver{vars: m}
}

// Resolve replaces all variables in the URN with their values.
func (r *Resolver) Resolve(u URN) URN {
	if r == nil || len(r.vars) == 0 {
		return u
	}
	return URN(r.ResolveString(string(u)))
}

// ResolveString replaces all variables in the string with their values.
func (r *Resolver) ResolveString(s string) string {
	if r == nil || len(r.vars) == 0 {
		return s
	}
	for k, v := range r.vars {
		s = strings.ReplaceAll(s, k, v)
	}
	return s
}

// ContainsVar reports whether the URN contains any of the resolver's variables.
func (r *Resolver) ContainsVar(u URN) bool {
	if r == nil || len(r.vars) == 0 {
		return false
	}
	s := string(u)
	for k := range r.vars {
		if strings.Contains(s, k) {
			return true
		}
	}
	return false
}

// Vars returns a copy of the variable mappings.
func (r *Resolver) Vars() map[string]string {
	if r == nil {
		return nil
	}
	m := make(map[string]string, len(r.vars))
	maps.Copy(m, r.vars)
	return m
}
