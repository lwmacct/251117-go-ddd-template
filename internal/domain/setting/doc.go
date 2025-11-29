// Package setting 定义系统配置领域模型。
//
// 本包管理应用程序的动态配置项，支持运行时修改而无需重启服务。
// 适用于：
//   - 功能开关：动态启用/禁用特定功能
//   - 业务参数：如分页大小、超时时间等可调参数
//   - 系统设置：站点名称、Logo、联系方式等
//
// 配置分类 (Category)：
//   - general: 通用配置
//   - security: 安全相关配置
//   - notification: 通知配置
//   - backup: 备份配置
//
// 值类型 (ValueType)：
//   - string: 字符串
//   - number: 数值
//   - boolean: 布尔值
//   - json: JSON 对象
//
// 与配置文件的区别：
// 配置文件适合静态配置（数据库连接、端口等），本模块适合需要运行时调整的动态配置。
package setting
