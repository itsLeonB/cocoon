package entity

type User struct {
	BaseEntity
	Email    string
	Password string
	Profile  UserProfile
}
