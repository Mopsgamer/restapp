package internal

import (
	"restapp/internal/model"
	"restapp/internal/model_request"
	"strconv"
)

func (r Responder) MessageCreate() error {
	req := new(model_request.MessageCreate)
	if err := r.Bind().URI(req); err != nil {
		return r.SendString(MessageErrInvalidRequest)
	}
	if err := r.Bind().Form(req); err != nil {
		return r.SendString(MessageErrInvalidRequest)
	}

	user, _ := r.User()
	if user == nil {
		return nil
	}

	message := req.Message(user.Id)
	if r.DB.GroupMemberById(message.GroupId, message.AuthorId) == nil {
		return r.SendString(MessageErrNotGroupMember)
	}

	if !model.IsValidMessageContent(message.Content) {
		return r.SendString(MessageErrMessageContent + " Length: " + strconv.Itoa(len(message.Content)) + "/" + model.ContentMaxLengthString)
	}

	messageId := r.DB.MessageCreate(*message)
	if messageId == nil {
		return r.SendString(MessageFatalDatabaseQuery)
	}

	return r.SendString("")
}
