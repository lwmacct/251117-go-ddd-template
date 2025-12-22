package menu

import "errors"

var (
	// ErrMenuNotFound 菜单不存在
	ErrMenuNotFound = errors.New("menu not found")

	// ErrParentMenuNotFound 父菜单不存在
	ErrParentMenuNotFound = errors.New("parent menu not found")

	// ErrCircularReference 循环引用（菜单不能是自己的子菜单）
	ErrCircularReference = errors.New("circular reference detected")

	// ErrInvalidMenuPath 无效的菜单路径
	ErrInvalidMenuPath = errors.New("invalid menu path")

	// ErrMenuHasChildren 菜单下有子菜单
	ErrMenuHasChildren = errors.New("menu has children")

	// ErrInvalidMenuOrder 无效的菜单排序
	ErrInvalidMenuOrder = errors.New("invalid menu order")

	// ErrDuplicateMenuPath 菜单路径重复
	ErrDuplicateMenuPath = errors.New("duplicate menu path")
)
