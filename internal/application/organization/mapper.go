package organization

import "github.com/lwmacct/251117-go-ddd-template/internal/domain/organization"

// ============================================================================
// Organization Mappers
// ============================================================================

// ToOrgDTO 将组织实体转换为 DTO
func ToOrgDTO(org *organization.Organization) *OrgDTO {
	if org == nil {
		return nil
	}
	return &OrgDTO{
		ID:          org.ID,
		Name:        org.Name,
		DisplayName: org.DisplayName,
		Description: org.Description,
		Avatar:      org.Avatar,
		Status:      org.Status,
		CreatedAt:   org.CreatedAt,
		UpdatedAt:   org.UpdatedAt,
	}
}

// ToOrgDTOs 将组织实体列表转换为 DTO 列表
func ToOrgDTOs(orgs []*organization.Organization) []*OrgDTO {
	if len(orgs) == 0 {
		return []*OrgDTO{}
	}
	dtos := make([]*OrgDTO, 0, len(orgs))
	for _, org := range orgs {
		if dto := ToOrgDTO(org); dto != nil {
			dtos = append(dtos, dto)
		}
	}
	return dtos
}

// ============================================================================
// Team Mappers
// ============================================================================

// ToTeamDTO 将团队实体转换为 DTO
func ToTeamDTO(team *organization.Team) *TeamDTO {
	if team == nil {
		return nil
	}
	return &TeamDTO{
		ID:             team.ID,
		OrganizationID: team.OrganizationID,
		Name:           team.Name,
		DisplayName:    team.DisplayName,
		Description:    team.Description,
		CreatedAt:      team.CreatedAt,
		UpdatedAt:      team.UpdatedAt,
	}
}

// ToTeamDTOs 将团队实体列表转换为 DTO 列表
func ToTeamDTOs(teams []*organization.Team) []*TeamDTO {
	if len(teams) == 0 {
		return []*TeamDTO{}
	}
	dtos := make([]*TeamDTO, 0, len(teams))
	for _, team := range teams {
		if dto := ToTeamDTO(team); dto != nil {
			dtos = append(dtos, dto)
		}
	}
	return dtos
}

// ============================================================================
// Member Mappers
// ============================================================================

// ToMemberDTO 将成员实体转换为 DTO
func ToMemberDTO(member *organization.Member) *MemberDTO {
	if member == nil {
		return nil
	}
	return &MemberDTO{
		ID:             member.ID,
		OrganizationID: member.OrganizationID,
		UserID:         member.UserID,
		Role:           string(member.Role),
		JoinedAt:       member.JoinedAt,
	}
}

// ToMemberDTOs 将成员实体列表转换为 DTO 列表
func ToMemberDTOs(members []*organization.Member) []*MemberDTO {
	if len(members) == 0 {
		return []*MemberDTO{}
	}
	dtos := make([]*MemberDTO, 0, len(members))
	for _, member := range members {
		if dto := ToMemberDTO(member); dto != nil {
			dtos = append(dtos, dto)
		}
	}
	return dtos
}

// ============================================================================
// Team Member Mappers
// ============================================================================

// ToTeamMemberDTO 将团队成员实体转换为 DTO
func ToTeamMemberDTO(member *organization.TeamMember) *TeamMemberDTO {
	if member == nil {
		return nil
	}
	return &TeamMemberDTO{
		ID:       member.ID,
		TeamID:   member.TeamID,
		UserID:   member.UserID,
		Role:     string(member.Role),
		JoinedAt: member.JoinedAt,
	}
}

// ToTeamMemberDTOs 将团队成员实体列表转换为 DTO 列表
func ToTeamMemberDTOs(members []*organization.TeamMember) []*TeamMemberDTO {
	if len(members) == 0 {
		return []*TeamMemberDTO{}
	}
	dtos := make([]*TeamMemberDTO, 0, len(members))
	for _, member := range members {
		if dto := ToTeamMemberDTO(member); dto != nil {
			dtos = append(dtos, dto)
		}
	}
	return dtos
}
