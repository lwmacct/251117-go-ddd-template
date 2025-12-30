// Package urn provides a Unified Resource Name type for access control systems.
//
// # Overview
//
// URN is a three-part identifier format: {scope}:{type}:{identifier}
//
//	{scope}:{type}:{identifier}
//	   │      │        │
//	   │      │        └─ Action (create) or ID (123/*/variable)
//	   │      └─ Module/resource type
//	   └─ Scope (can be hierarchical with . separator)
//
// # Separator Rules
//
//   - `:` separates main parts (scope, type, identifier)
//   - `.` separates scope hierarchy (e.g., sys.admin, org.acme.team.dev)
//
// # Examples
//
// Operations:
//
//	sys:users:create      // System-level user creation
//	self:profile:update   // User updates own profile
//	public:auth:login     // Public login operation
//
// Resources:
//
//	sys:user:123          // System user 123
//	self:user:@me         // Current user (requires variable resolution)
//	org.acme:user:*       // All users in org.acme
//
// # Wildcard Matching
//
//	*:*:*           → Matches all (super admin)
//	sys:*:*         → All operations in sys scope
//	sys.*:*:*       → sys and all sub-scopes (sys.admin, sys.readonly)
//	org.acme:users:* → All user operations in org.acme
//
// # Variable Resolution
//
// URN supports runtime variable substitution using [Resolver]:
//
//	r := urn.NewResolver(map[string]string{
//	    "@me":  "123",
//	    "@org": "acme",
//	})
//	r.ResolveString("self:user:@me")   // "self:user:123"
//	r.ResolveString("org.@org:*:*")    // "org.acme:*:*"
//
// # Thread Safety
//
// [URN] is an immutable string type. All methods are safe for concurrent use.
// [Resolver] is also safe for concurrent use after construction.
package urn
