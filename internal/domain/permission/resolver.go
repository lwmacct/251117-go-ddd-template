package permission

import (
	"maps"
	"strings"
)

// Resolver 执行 URN 中的变量替换。
//
// 变量是任意字符串，会被替换为对应的值。
// 常用约定使用 @ 前缀（如 @me, @org），但任何字符串都可以。
type Resolver struct {
	vars map[string]string
}

// NewResolver 创建变量解析器。
//
// 示例：
//
//	r := NewResolver(map[string]string{
//	    "@me":  "123",
//	    "@org": "acme",
//	})
//	r.ResolveString("self:user:@me")  // "self:user:123"
//	r.ResolveString("org.@org:*:*")   // "org.acme:*:*"
func NewResolver(vars map[string]string) *Resolver {
	// 复制 map 防止外部修改
	m := make(map[string]string, len(vars))
	maps.Copy(m, vars)
	return &Resolver{vars: m}
}

// Resolve 替换 Operation 中的所有变量。
func (r *Resolver) Resolve(o Operation) Operation {
	if r == nil || len(r.vars) == 0 {
		return o
	}
	return Operation(r.ResolveString(string(o)))
}

// ResolveResource 替换 Resource 中的所有变量。
func (r *Resolver) ResolveResource(res Resource) Resource {
	if r == nil || len(r.vars) == 0 {
		return res
	}
	return Resource(r.ResolveString(string(res)))
}

// ResolveString 替换字符串中的所有变量。
func (r *Resolver) ResolveString(s string) string {
	if r == nil || len(r.vars) == 0 {
		return s
	}
	for k, v := range r.vars {
		s = strings.ReplaceAll(s, k, v)
	}
	return s
}

// ContainsVar 报告字符串是否包含解析器中的任何变量。
func (r *Resolver) ContainsVar(s string) bool {
	if r == nil || len(r.vars) == 0 {
		return false
	}
	for k := range r.vars {
		if strings.Contains(s, k) {
			return true
		}
	}
	return false
}

// Vars 返回变量映射的副本。
func (r *Resolver) Vars() map[string]string {
	if r == nil {
		return nil
	}
	m := make(map[string]string, len(r.vars))
	maps.Copy(m, r.vars)
	return m
}
