package operation

// Registry is a generic operation registry that supports any metadata type.
//
// Usage:
//
//	type myMeta struct {
//	    Method string
//	    Path   string
//	}
//
//	var registry = operation.Registry[myMeta]{}
//	registry.Register("sys:users:create", myMeta{Method: "POST", Path: "/users"})
//
//	if meta, ok := registry.Get("sys:users:create"); ok {
//	    fmt.Println(meta.Method) // "POST"
//	}
type Registry[M any] map[Operation]M

// Get retrieves the metadata for an operation.
// Returns the metadata and true if found, zero value and false otherwise.
func (r Registry[M]) Get(op Operation) (M, bool) {
	m, ok := r[op]
	return m, ok
}

// Register adds or updates operation metadata.
func (r Registry[M]) Register(op Operation, meta M) {
	r[op] = meta
}

// All returns all registered operations and their metadata.
func (r Registry[M]) All() map[Operation]M {
	return r
}

// Operations returns all registered operation keys.
func (r Registry[M]) Operations() []Operation {
	ops := make([]Operation, 0, len(r))
	for op := range r {
		ops = append(ops, op)
	}
	return ops
}

// Len returns the number of registered operations.
func (r Registry[M]) Len() int {
	return len(r)
}
