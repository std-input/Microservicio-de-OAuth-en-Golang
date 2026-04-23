package models

import "time"

type User struct {
	ID          string `json:"id" gorm:"primaryKey"`
	Email       string `json:"email" gorm:"uniqueIndex"`
	Name        string `json:"given_name"`
	FamilyName  string `json:"family_name"`
	Picture     string `json:"picture"`
	Description string `json:"description"`
	Role        string `json:"role"`
}

type SafeUser struct {
	ID          string `json:"id"`
	Name        string `json:"given_name"`
	FamilyName  string `json:"family_name"`
	Picture     string `json:"picture"`
	Description string `json:"description"`
	Role        string `json:"role"`
}

type RefreshToken struct {
	User       User   `gorm:"foreignKey:UserID;references:ID"`
	UserID     string `gorm:"primaryKey"`
	Expiration time.Time
	Token      string
}
