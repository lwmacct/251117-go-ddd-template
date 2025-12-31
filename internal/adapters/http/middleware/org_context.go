package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/org"
)

// OrgContext 组织上下文中间件。
// 从路由参数提取 org_id，验证当前用户是否为组织成员。
// 验证通过后注入以下值到 Gin Context:
//   - org_id: uint - 组织 ID
//   - org_role: string - 用户在组织中的角色 (owner/admin/member)
func OrgContext(memberQuery org.MemberQueryRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取当前用户 ID
		userID, exists := c.Get("user_id")
		if !exists {
			response.Unauthorized(c, "user not authenticated")
			c.Abort()
			return
		}
		uid, ok := userID.(uint)
		if !ok {
			response.InternalError(c, "invalid user ID format")
			c.Abort()
			return
		}

		// 2. 解析路由参数中的 org_id
		orgIDStr := c.Param("org_id")
		if orgIDStr == "" {
			response.BadRequest(c, "org_id is required")
			c.Abort()
			return
		}
		orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "invalid org_id")
			c.Abort()
			return
		}

		// 3. 验证用户是否是组织成员
		member, err := memberQuery.GetByOrgAndUser(c.Request.Context(), uint(orgID), uid)
		if err != nil {
			response.Forbidden(c, "not a member of this organization")
			c.Abort()
			return
		}

		// 4. 注入组织上下文
		c.Set("org_id", uint(orgID))
		c.Set("org_role", string(member.Role))

		c.Next()
	}
}

// TeamContext 团队上下文中间件。
// 必须在 OrgContext 之后使用。
// 从路由参数提取 team_id，验证团队属于组织，验证用户是否为团队成员。
// 验证通过后注入以下值到 Gin Context:
//   - team_id: uint - 团队 ID
//   - team_role: string - 用户在团队中的角色 (lead/member)
func TeamContext(
	teamQuery org.TeamQueryRepository,
	teamMemberQuery org.TeamMemberQueryRepository,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取组织上下文（必须先通过 OrgContext）
		orgID, exists := c.Get("org_id")
		if !exists {
			response.InternalError(c, "org context not found, OrgContext middleware required")
			c.Abort()
			return
		}
		oid, ok := orgID.(uint)
		if !ok {
			response.InternalError(c, "invalid org_id format in context")
			c.Abort()
			return
		}

		// 2. 获取当前用户 ID
		userID, exists := c.Get("user_id")
		if !exists {
			response.Unauthorized(c, "user not authenticated")
			c.Abort()
			return
		}
		uid, ok := userID.(uint)
		if !ok {
			response.InternalError(c, "invalid user ID format")
			c.Abort()
			return
		}

		// 3. 解析路由参数中的 team_id
		teamIDStr := c.Param("team_id")
		if teamIDStr == "" {
			response.BadRequest(c, "team_id is required")
			c.Abort()
			return
		}
		teamID, err := strconv.ParseUint(teamIDStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "invalid team_id")
			c.Abort()
			return
		}

		// 4. 验证团队属于组织
		belongsTo, err := teamQuery.BelongsToOrg(c.Request.Context(), uint(teamID), oid)
		if err != nil {
			response.InternalError(c, "failed to verify team organization")
			c.Abort()
			return
		}
		if !belongsTo {
			response.NotFound(c, "team not found in this organization")
			c.Abort()
			return
		}

		// 5. 验证用户是否是团队成员
		member, err := teamMemberQuery.GetByTeamAndUser(c.Request.Context(), uint(teamID), uid)
		if err != nil {
			response.Forbidden(c, "not a member of this team")
			c.Abort()
			return
		}

		// 6. 注入团队上下文
		c.Set("team_id", uint(teamID))
		c.Set("team_role", string(member.Role))

		c.Next()
	}
}

// OrgContextOptional 可选的组织上下文中间件。
// 与 OrgContext 类似，但不要求用户必须是组织成员。
// 仅验证组织存在性，适用于需要组织 ID 但不限制成员的场景。
func OrgContextOptional(orgQuery org.QueryRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 解析路由参数中的 org_id
		orgIDStr := c.Param("org_id")
		if orgIDStr == "" {
			response.BadRequest(c, "org_id is required")
			c.Abort()
			return
		}
		orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
		if err != nil {
			response.BadRequest(c, "invalid org_id")
			c.Abort()
			return
		}

		// 2. 验证组织存在
		_, err = orgQuery.GetByID(c.Request.Context(), uint(orgID))
		if err != nil {
			response.NotFound(c, "organization")
			c.Abort()
			return
		}

		// 3. 注入组织 ID（不注入角色，因为用户可能不是成员）
		c.Set("org_id", uint(orgID))

		c.Next()
	}
}
