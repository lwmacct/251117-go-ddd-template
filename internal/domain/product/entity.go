package product

import "time"

// Status 产品状态
type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)

// Type 产品类型
type Type string

const (
	TypePersonal Type = "personal" // 个人产品
	TypeTeam     Type = "team"     // 团队产品（组织订阅）
)

// Product 产品实体。
//
// 产品是多租户系统中可供组织或用户订阅的服务单元。
type Product struct {
	ID          uint      `json:"id"`
	Code        string    `json:"code"`        // 产品代码，用于路由和前端识别，如 "teamTask", "crm"
	Name        string    `json:"name"`        // 产品名称
	Type        Type      `json:"type"`        // 产品类型：personal 或 team
	Description string    `json:"description"` // 产品描述
	Price       float64   `json:"price"`       // 基础价格（月付）
	Status      Status    `json:"status"`      // 产品状态
	LayoutRef   string    `json:"layout_ref"`  // 前端 Layout 组件引用，如 "TeamTaskLayout"
	MaxSeats    int       `json:"max_seats"`   // 团队产品最大席位数量，0 表示无限制
	TrialDays   int       `json:"trial_days"`  // 试用天数，0 表示无试用
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

// IsTeamProduct 报告是否为团队产品。
func (p *Product) IsTeamProduct() bool {
	return p.Type == TypeTeam
}

// HasTrial 报告是否支持试用。
func (p *Product) HasTrial() bool {
	return p.TrialDays > 0
}

// GetMaxSeats 报告最大席位数量，0 表示无限制。
func (p *Product) GetMaxSeats() int {
	if p.MaxSeats <= 0 {
		return 0 // 无限制
	}
	return p.MaxSeats
}
