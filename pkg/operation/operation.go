package operation

// Operation is a unified operation identifier.
//
// Format: {scope}:{type}:{identifier}
//
// Common scopes:
//   - public: Public operations (no authentication required)
//   - sys:    System administration operations
//   - self:   User self-service operations
//
// Examples:
//   - public:auth:login
//   - sys:users:create
//   - self:profile:update
type Operation string

// Scope returns the scope part of the operation.
//
//	Operation("sys:users:create").Scope() // "sys"
func (o Operation) Scope() string {
	return URN(o).Scope()
}

// Type returns the type part of the operation.
//
//	Operation("sys:users:create").Type() // "users"
func (o Operation) Type() string {
	return URN(o).Type()
}

// Identifier returns the identifier part of the operation.
//
//	Operation("sys:users:create").Identifier() // "create"
func (o Operation) Identifier() string {
	return URN(o).Identifier()
}

// String returns the operation as a string.
func (o Operation) String() string {
	return string(o)
}

// HTTPMethod represents an HTTP request method.
type HTTPMethod string

// HTTP method constants.
const (
	GET    HTTPMethod = "GET"
	POST   HTTPMethod = "POST"
	PUT    HTTPMethod = "PUT"
	DELETE HTTPMethod = "DELETE"
	PATCH  HTTPMethod = "PATCH"
)
