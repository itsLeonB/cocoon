package entity

import "github.com/itsLeonB/go-crud"

type User struct {
	crud.BaseEntity
	Email    string
	Password string
	Profile  UserProfile
}
