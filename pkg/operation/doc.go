// Package operation provides types for URN-based access control systems.
//
// This package provides a complete toolkit for operation-centric RBAC:
//   - [Operation] and [Resource] types for access control
//   - [URN] type for parsing and manipulation
//   - [Match], [MatchOperation], [MatchResource] for pattern matching
//   - [Resolver] for variable substitution
//   - [Registry] for operation metadata storage
//
// # URN Format
//
// All types follow Unified Resource Name format: {scope}:{type}:{identifier}
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
// # Operations
//
// Operations identify actions in the system:
//
//	public:auth:login       // Public scope, auth type, login action
//	sys:users:create        // System scope, users type, create action
//	self:profile:update     // Self scope, profile type, update action
//
// # Resources
//
// Resources identify entities being acted upon:
//
//	sys:user:123            // System user 123
//	self:user:@me           // Current user (requires variable resolution)
//	org.acme:user:*         // All users in org.acme
//
// # Wildcard Matching
//
// Pattern matching supports wildcards:
//
//	*:*:*           → Matches all (super admin)
//	sys:*:*         → All operations in sys scope
//	sys.*:*:*       → sys and all sub-scopes (sys.admin, sys.readonly)
//	org.acme:users:* → All user operations in org.acme
//
// Example:
//
//	MatchOperation("sys:*:*", "sys:users:create")  // true
//	MatchResource("*:*:*", "sys:user:123")         // true
//
// # Variable Resolution
//
// [Resolver] supports runtime variable substitution:
//
//	r := NewResolver(map[string]string{
//	    "@me":  "123",
//	    "@org": "acme",
//	})
//	r.ResolveString("self:user:@me")   // "self:user:123"
//	r.ResolveString("org.@org:*:*")    // "org.acme:*:*"
//
// # Registry Usage
//
// The [Registry] type allows projects to define custom metadata:
//
//	type myMeta struct {
//	    Method string
//	    Path   string
//	}
//
//	var registry = operation.Registry[myMeta]{}
//	registry.Register("sys:users:create", myMeta{Method: "POST", Path: "/users"})
//
// # Thread Safety
//
// [URN] and [Operation] are immutable string types. All methods are safe for
// concurrent use. [Resolver] is also safe for concurrent use after construction.
package operation
