package internal

import (
	"errors"
	"restapp/internal/model"
)

// Create new DB record.
func (db Database) UserCreate(user model.User) error {
	query :=
		`INSERT INTO app_users (
			nickname,
			username,
			email,
			phone,
			password,
			avatar,
			created_at,
			last_seen
		)
    	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := db.Sql.Exec(query,
		user.Nick,
		user.Name,
		user.Email,
		user.Phone,
		user.Password,
		user.Avatar,
		user.CreatedAt,
		user.LastSeen,
	)
	return err
}

// Change the existing DB record.
func (db Database) UserUpdate(user model.User) error {
	query :=
		`UPDATE app_users
    	SET
		nickname = ?,
		username = ?,
		email = ?,
		phone = ?,
		password = ?,
		avatar = ?,
		created_at = ?,
		last_seen = ?

        WHERE id = ?`
	_, err := db.Sql.Exec(query,
		user.Nick,
		user.Name,
		user.Email,
		user.Phone,
		user.Password,
		user.Avatar,
		user.CreatedAt,
		user.LastSeen,
		user.Id,
	)
	return err
}

// Delete the existing DB record.
func (db Database) UserDelete(userId uint) error {
	query := `DELETE FROM app_users WHERE id = ?`
	_, err := db.Sql.Exec(query, userId)
	return err
}

// Get the user by his email.
func (db Database) UserByEmail(email string) (*model.User, error) {
	user := new(model.User)
	query := `SELECT * FROM app_users WHERE email = ?`
	err := db.Sql.Get(user, query, email)
	if err != nil {
		user = nil
	}
	return user, err
}

// Get the user by his identificator.
func (db Database) UserById(userId uint) (*model.User, error) {
	user := new(model.User)
	query := `SELECT * FROM app_users WHERE id = ?`
	err := db.Sql.Get(user, query, userId)
	if err != nil {
		user = nil
	}
	return user, err
}

// Get the user by his username.
func (db Database) UserByUsername(username string) (*model.User, error) {
	user := new(model.User)
	query := `SELECT * FROM app_users WHERE username = ?`
	err := db.Sql.Get(user, query, username)
	if err != nil {
		user = nil
	}
	return user, err
}

func (db Database) UserOwnGroups(userId uint) (*[]model.Group, error) {
	groupList := new([]model.Group)
	query := `SELECT * FROM app_group_members WHERE user_id = ?
	RIGHT JOIN app_groups ON app_groups.id = app_group_members.group_id
	GROUP BY app_group_members.group_id
	`
	err := db.Sql.Get(groupList, query, userId)
	if err != nil {
		groupList = nil
	}
	return groupList, err
}

func (db Database) UserJoinGroup(newMember model.Member, group model.Group, groupMembers []model.Member) error {
	if group.Mode == model.GroupModeDm && len(groupMembers) <= 2 {
		return errors.New("can not join to the DM. already full")
	}
	query :=
		`INSERT INTO app_group_members (
			group_id,
			user_id,
			is_owner,
			is_banned,
			membername
		)
    	VALUES (?, ?, ?, ?, ?)`
	_, err := db.Sql.Exec(query,
		newMember.GroupId,
		newMember.UserId,
		newMember.IsOwner,
		newMember.IsBanned,
		newMember.Nick,
	)

	return err
}

func (db Database) UserLeaveGroup(userId uint, groupId uint) error {
	query := `DELETE FROM app_group_members WHERE group_id = ? AND user_id = ?`
	_, err := db.Sql.Exec(query, groupId, userId)
	return err
}

func (db Database) MessageCreate(message model.Message) error {
	query :=
		`INSERT INTO app_group_messages (
			id,
			group_id,
			author_id,
			content,
			created_at
		)
    	VALUES (?, ?, ?, ?, ?)`
	_, err := db.Sql.Exec(query,
		message.Id,
		message.GroupId,
		message.AuthorId,
		message.Content,
		message.CreatedAt,
	)

	return err
}
