package operation

// Resource is a resource identifier.
//
// Format: {scope}:{type}:{id}
//
// Examples:
//   - sys:user:123       System user with ID 123
//   - self:user:@me      Current user (runtime substitution)
//   - org.acme:user:*    All users in org.acme organization
//   - *:*:*              All resources (wildcard)
//
// Special identifiers:
//   - * matches any value
//   - @me is substituted with current user ID at runtime
//   - @org is substituted with current organization ID at runtime
type Resource string

// Predefined resource constants.
const (
	// ResourceAll matches all resources.
	ResourceAll Resource = "*:*:*"
)

// String returns the resource as a string.
func (r Resource) String() string {
	return string(r)
}
