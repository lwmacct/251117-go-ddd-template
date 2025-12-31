package product

import "time"

// Status 产品状态
type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)

// Product 产品实体。
//
// 产品是多租户系统中可供组织或用户订阅的服务单元。
type Product struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// IsActive 报告产品是否处于激活状态。
func (p *Product) IsActive() bool {
	return p.Status == StatusActive
}

// Activate 激活产品。
func (p *Product) Activate() {
	p.Status = StatusActive
}

// Deactivate 停用产品。
func (p *Product) Deactivate() {
	p.Status = StatusInactive
}
