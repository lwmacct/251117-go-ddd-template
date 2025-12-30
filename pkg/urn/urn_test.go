package urn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURN_Parse(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantScope string
		wantType  string
		wantID    string
		wantParts []string
	}{
		{
			name:      "full URN",
			input:     "sys:users:create",
			wantScope: "sys",
			wantType:  "users",
			wantID:    "create",
			wantParts: []string{"sys"},
		},
		{
			name:      "hierarchical scope",
			input:     "org.acme.team.dev:tasks:create",
			wantScope: "org.acme.team.dev",
			wantType:  "tasks",
			wantID:    "create",
			wantParts: []string{"org", "acme", "team", "dev"},
		},
		{
			name:      "full wildcard",
			input:     "*:*:*",
			wantScope: "*",
			wantType:  "*",
			wantID:    "*",
			wantParts: []string{"*"},
		},
		{
			name:      "single wildcard",
			input:     "*",
			wantScope: "*",
			wantType:  "*",
			wantID:    "*",
			wantParts: []string{"*"},
		},
		{
			name:      "two parts",
			input:     "sys:users",
			wantScope: "sys",
			wantType:  "users",
			wantID:    "*",
			wantParts: []string{"sys"},
		},
		{
			name:      "one part",
			input:     "sys",
			wantScope: "sys",
			wantType:  "*",
			wantID:    "*",
			wantParts: []string{"sys"},
		},
		{
			name:      "with variable",
			input:     "self:user:@me",
			wantScope: "self",
			wantType:  "user",
			wantID:    "@me",
			wantParts: []string{"self"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := URN(tt.input)
			parts := u.Parse()

			assert.Equal(t, tt.wantScope, parts.Scope, "Scope")
			assert.Equal(t, tt.wantType, parts.Type, "Type")
			assert.Equal(t, tt.wantID, parts.Identifier, "Identifier")
			assert.Equal(t, tt.wantParts, parts.ScopeParts, "ScopeParts")
		})
	}
}

func TestURN_Methods(t *testing.T) {
	u := URN("sys.admin:config:update")

	assert.Equal(t, "sys.admin", u.Scope())
	assert.Equal(t, "config", u.Type())
	assert.Equal(t, "update", u.Identifier())
	assert.Equal(t, "sys.admin:config:update", u.String())
	assert.False(t, u.IsWildcard())

	assert.True(t, URN("*").IsWildcard())
	assert.True(t, URN("*:*:*").IsWildcard())
}

func TestNew(t *testing.T) {
	u := New("sys", "users", "create")
	assert.Equal(t, URN("sys:users:create"), u)
}

func TestURN_Match(t *testing.T) {
	tests := []struct {
		name    string
		target  string
		pattern string
		want    bool
	}{
		// Super wildcard
		{"super wildcard *:*:*", "sys:users:create", "*:*:*", true},
		{"super wildcard *", "sys:users:create", "*", true},

		// Exact match
		{"exact match", "sys:users:create", "sys:users:create", true},
		{"exact mismatch", "sys:users:create", "sys:users:delete", false},

		// Scope wildcard
		{"scope wildcard", "sys:users:create", "*:users:create", true},
		{"scope mismatch", "sys:users:create", "self:users:create", false},

		// Type wildcard
		{"type wildcard", "sys:users:create", "sys:*:create", true},
		{"type mismatch", "sys:users:create", "sys:roles:create", false},

		// Identifier wildcard
		{"id wildcard", "sys:users:create", "sys:users:*", true},
		{"id mismatch", "sys:users:create", "sys:users:delete", false},

		// Hierarchical scope
		{"hierarchical exact", "sys.admin:config:update", "sys.admin:config:update", true},
		{"hierarchical wildcard matches self", "sys:users:create", "sys.*:*:*", true},
		{"hierarchical wildcard matches child", "sys.admin:config:update", "sys.*:*:*", true},
		{"hierarchical wildcard matches grandchild", "sys.admin.readonly:config:read", "sys.*:*:*", true},
		{"hierarchical no match different root", "org.acme:users:list", "sys.*:*:*", false},

		// Org hierarchy
		{"org exact", "org.acme:users:list", "org.acme:users:list", true},
		{"org wildcard", "org.acme.team.dev:tasks:create", "org.acme.*:*:*", true},
		{"org partial no match", "org.acme:users:list", "org.other:*:*", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := URN(tt.target).Match(URN(tt.pattern))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMatch_Function(t *testing.T) {
	assert.True(t, Match("*:*:*", "sys:users:create"))
	assert.True(t, Match("sys:*:*", "sys:users:create"))
	assert.False(t, Match("self:*:*", "sys:users:create"))
}

func TestResolver(t *testing.T) {
	r := NewResolver(map[string]string{
		"@me":  "123",
		"@org": "acme",
	})

	t.Run("Resolve URN", func(t *testing.T) {
		assert.Equal(t, URN("self:user:123"), r.Resolve("self:user:@me"))
		assert.Equal(t, URN("org.acme:user:*"), r.Resolve("org.@org:user:*"))
	})

	t.Run("ResolveString", func(t *testing.T) {
		assert.Equal(t, "self:user:123", r.ResolveString("self:user:@me"))
		assert.Equal(t, "org.acme:user:*", r.ResolveString("org.@org:user:*"))
		assert.Equal(t, "sys:users:create", r.ResolveString("sys:users:create"))
	})

	t.Run("ContainsVar", func(t *testing.T) {
		assert.True(t, r.ContainsVar("self:user:@me"))
		assert.True(t, r.ContainsVar("org.@org:*:*"))
		assert.False(t, r.ContainsVar("sys:users:create"))
	})

	t.Run("Vars", func(t *testing.T) {
		vars := r.Vars()
		assert.Equal(t, "123", vars["@me"])
		assert.Equal(t, "acme", vars["@org"])

		// Ensure it's a copy
		vars["@me"] = "456"
		assert.Equal(t, "123", r.Vars()["@me"])
	})
}

func TestResolver_NilSafe(t *testing.T) {
	var r *Resolver

	assert.Equal(t, URN("self:user:@me"), r.Resolve("self:user:@me"))
	assert.Equal(t, "self:user:@me", r.ResolveString("self:user:@me"))
	assert.False(t, r.ContainsVar("self:user:@me"))
	assert.Nil(t, r.Vars())
}

func TestResolver_EmptyVars(t *testing.T) {
	r := NewResolver(nil)

	assert.Equal(t, URN("self:user:@me"), r.Resolve("self:user:@me"))
	assert.Equal(t, "self:user:@me", r.ResolveString("self:user:@me"))
	assert.False(t, r.ContainsVar("self:user:@me"))
}

func TestResolver_CustomVars(t *testing.T) {
	// Demonstrate that any variable format works
	r := NewResolver(map[string]string{
		"${user}": "alice",
		"$tenant": "corp",
	})

	assert.Equal(t, "tenant.corp:user:alice", r.ResolveString("tenant.$tenant:user:${user}"))
}
