package model_database

import "github.com/jmoiron/sqlx/types"

type Role struct {
	Id    uint32  `db:"id"`
	Name  string  `db:"name"`
	Color *uint32 `db:"color"`

	ChatRead     types.BitBool `db:"perm_chat_read"`
	ChatWrite    types.BitBool `db:"perm_chat_write"`
	ChatDelete   types.BitBool `db:"perm_chat_delete"`
	Kick         types.BitBool `db:"perm_kick"`
	Ban          types.BitBool `db:"perm_ban"`
	ChangeGroup  types.BitBool `db:"perm_change_group"`
	ChangeMember types.BitBool `db:"perm_change_member"`
}

var RoleDefault = Role{
	Name:       "everyone",
	Color:      nil,
	ChatRead:   true,
	ChatWrite:  true,
	ChatDelete: true,
}

func mergeProp(prop types.BitBool, prop2 types.BitBool) types.BitBool {
	return prop || prop2
}

// Enabled rights have priority over disabled rights.
func (r *Role) Merge(roleList ...Role) {
	for _, role := range roleList {
		r.ChatRead = mergeProp(r.ChatRead, role.ChatRead)
		r.ChatWrite = mergeProp(r.ChatWrite, role.ChatWrite)
		r.ChatDelete = mergeProp(r.ChatDelete, role.ChatDelete)
		r.Kick = mergeProp(r.Kick, role.Kick)
		r.Ban = mergeProp(r.Ban, role.Ban)
		r.ChangeGroup = mergeProp(r.ChangeGroup, role.ChangeGroup)
		r.ChangeMember = mergeProp(r.ChangeMember, role.ChangeMember)
	}
}
