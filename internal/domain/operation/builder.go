package operation

import op "github.com/lwmacct/251117-go-ddd-template/pkg/operation"

// ============================================================================
// Definer (操作定义器)
// ============================================================================

// Definer 操作定义构建器，提供链式 DSL 定义操作元数据。
type Definer struct {
	op   Operation
	meta operationMeta
}

// D 创建操作定义器。
func D(o Operation) *Definer {
	return &Definer{op: o}
}

// Use 应用预定义基类配置。
func (d *Definer) Use(base *Base) *Definer {
	d.meta.Path = base.pathPrefix
	d.meta.Group = base.group
	d.meta.AuditCategory = base.category
	return d
}

// Get 设置 GET 方法和路径后缀。
func (d *Definer) Get(path string) *Definer {
	return d.route(HttpGET, path)
}

// Post 设置 POST 方法和路径后缀。
func (d *Definer) Post(path string) *Definer {
	return d.route(HttpPOST, path)
}

// Put 设置 PUT 方法和路径后缀。
func (d *Definer) Put(path string) *Definer {
	return d.route(HttpPUT, path)
}

// Delete 设置 DELETE 方法和路径后缀。
func (d *Definer) Delete(path string) *Definer {
	return d.route(HttpDELETE, path)
}

// Patch 设置 PATCH 方法和路径后缀。
func (d *Definer) Patch(path string) *Definer {
	return d.route(HttpPATCH, path)
}

// Audit 设置审计配置。
func (d *Definer) Audit(action string, auditOp AuditOperation) *Definer {
	d.meta.AuditAction = action
	d.meta.AuditOperation = auditOp
	return d
}

// AuditCreate 设置创建审计（自动推导 action）。
func (d *Definer) AuditCreate() *Definer {
	return d.Audit(string(d.meta.AuditCategory)+".create", AuditOpCreate)
}

// AuditUpdate 设置更新审计（自动推导 action）。
func (d *Definer) AuditUpdate() *Definer {
	return d.Audit(string(d.meta.AuditCategory)+".update", AuditOpUpdate)
}

// AuditDelete 设置删除审计（自动推导 action）。
func (d *Definer) AuditDelete() *Definer {
	return d.Audit(string(d.meta.AuditCategory)+".delete", AuditOpDelete)
}

// AuditAccess 设置访问审计（自动推导 action）。
func (d *Definer) AuditAccess() *Definer {
	return d.Audit(string(d.meta.AuditCategory)+".access", AuditOpAccess)
}

// AuditAuth 设置认证审计（自动推导 action）。
func (d *Definer) AuditAuth() *Definer {
	return d.Audit(string(d.meta.AuditCategory)+".authenticate", AuditOpAuthenticate)
}

// I18n 设置中英文标签。
func (d *Definer) I18n(label, desc string) *Definer {
	d.meta.Label = label
	d.meta.Description = desc
	return d
}

// Group 覆盖 Swagger 分组。
func (d *Definer) Group(g string) *Definer {
	d.meta.Group = g
	return d
}

// Build 返回 Operation 和 Meta 对。
func (d *Definer) Build() (Operation, operationMeta) {
	return d.op, d.meta
}

func (d *Definer) route(method HTTPMethod, path string) *Definer {
	d.meta.Method = method
	d.meta.Path += path
	return d
}

// ============================================================================
// Base (预定义基类)
// ============================================================================

// Base 操作基类配置，封装路径前缀、分组和审计分类。
type Base struct {
	pathPrefix string
	group      string
	category   AuditCategory
}

// B 创建基类构建器。
func B() *Base {
	return &Base{}
}

// Path 设置路径前缀。
func (b *Base) Path(prefix string) *Base {
	b.pathPrefix = prefix
	return b
}

// Group 设置 Swagger 分组。
func (b *Base) Group(g string) *Base {
	b.group = g
	return b
}

// Cat 设置审计分类。
func (b *Base) Cat(c AuditCategory) *Base {
	b.category = c
	return b
}

// ============================================================================
// Registry Builder
// ============================================================================

// Build 从 Definer 列表构建注册表。
func Build(definers ...*Definer) op.Registry[operationMeta] {
	reg := make(op.Registry[operationMeta], len(definers))
	for _, d := range definers {
		o, meta := d.Build()
		reg[o] = meta
	}
	return reg
}
