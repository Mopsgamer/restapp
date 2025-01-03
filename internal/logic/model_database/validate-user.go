package model_database

import (
	"regexp"
)

const (
	RegexpUserPassword string = "^[a-zA-Z0-9, .~\\-+%$^&*_!?()[\\]{}`]{8,255}$"
	RegexpUserNick     string = `^.{1,255}$`
	RegexpUserName     string = `^[a-zA-Z0-9._]{1,255}$`
	// Source: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/email#basic_validation
	RegexpUserEmail string = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	RegexpUserPhone string = `^\+?[1-9]\d{1,14}$`
)

func IsValidUserPassword(str string) bool {
	return IsValidString(str, RegexpUserPassword, 255)
}

func IsValidUserNick(str string) bool {
	return IsValidString(str, RegexpUserNick, 255)
}

func IsValidUserName(str string) bool {
	return IsValidString(str, RegexpUserName, 255)
}

func IsValidUserEmail(str string) bool {
	return IsValidString(str, RegexpUserEmail, 255)
}

func IsValidUserPhone(str *string) bool {
	if str == nil {
		return true // allow no phone
	}
	newstr := regexp.MustCompile(`\s`).ReplaceAllString(*str, "")
	return IsValidString(newstr, RegexpUserPhone, 255)
}
