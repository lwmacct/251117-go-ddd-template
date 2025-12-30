package urn

import "strings"

// URN is a Unified Resource Name.
//
// Format: {scope}:{type}:{identifier}
//
// Examples:
//   - sys:users:create     Operation (system user creation)
//   - self:user:@me        Resource (current user)
//   - *:*:*                Wildcard (matches all)
type URN string

// Parts contains the parsed components of a URN.
type Parts struct {
	Scope      string   // Scope, e.g., "sys.admin" or "org.acme.team.dev"
	ScopeParts []string // Scope hierarchy, e.g., ["sys", "admin"]
	Type       string   // Type/module, e.g., "users" or "user"
	Identifier string   // Identifier, e.g., "create" or "123" or "*"
}

// Parse parses the URN into its components.
//
// Parsing rules:
//   - "*" → Parts{Scope: "*", Type: "*", Identifier: "*"}
//   - "scope:type:id" → Parts{Scope: scope, Type: type, Identifier: id}
//   - "scope:type" → Parts{Scope: scope, Type: type, Identifier: "*"}
//   - "scope" → Parts{Scope: scope, Type: "*", Identifier: "*"}
func (u URN) Parse() Parts {
	s := string(u)

	// Super wildcard
	if s == "*" {
		return Parts{
			Scope:      "*",
			ScopeParts: []string{"*"},
			Type:       "*",
			Identifier: "*",
		}
	}

	parts := strings.SplitN(s, ":", 3)

	scope := "*"
	typ := "*"
	identifier := "*"

	switch len(parts) {
	case 3:
		scope = parts[0]
		typ = parts[1]
		identifier = parts[2]
	case 2:
		scope = parts[0]
		typ = parts[1]
	case 1:
		scope = parts[0]
	}

	// Parse scope hierarchy
	var scopeParts []string
	if scope == "*" {
		scopeParts = []string{"*"}
	} else {
		scopeParts = strings.Split(scope, ".")
	}

	return Parts{
		Scope:      scope,
		ScopeParts: scopeParts,
		Type:       typ,
		Identifier: identifier,
	}
}

// Scope returns the scope component.
//
// Examples:
//
//	"sys:users:create" → "sys"
//	"sys.admin:config:update" → "sys.admin"
func (u URN) Scope() string {
	return u.Parse().Scope
}

// Type returns the type/module component.
//
// Example: "sys:users:create" → "users"
func (u URN) Type() string {
	return u.Parse().Type
}

// Identifier returns the identifier component.
//
// Examples:
//
//	"sys:users:create" → "create"
//	"sys:user:123" → "123"
func (u URN) Identifier() string {
	return u.Parse().Identifier
}

// IsWildcard reports whether the URN is a full wildcard.
func (u URN) IsWildcard() bool {
	return string(u) == "*" || string(u) == "*:*:*"
}

// String returns the string representation of the URN.
func (u URN) String() string {
	return string(u)
}

// New creates a URN from scope, type, and identifier.
func New(scope, typ, identifier string) URN {
	return URN(scope + ":" + typ + ":" + identifier)
}
