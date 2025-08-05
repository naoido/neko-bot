package discordutil

import "github.com/bwmarrin/discordgo"

// HasAdminPermission : 管理者権限を持っている場合はtrue
func HasAdminPermission(member *discordgo.Member) bool {
	return member.Permissions&discordgo.PermissionAdministrator != 0
}
