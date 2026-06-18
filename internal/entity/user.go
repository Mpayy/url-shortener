package entity

import "time"

type User struct {
    ID        int64     `gorm:"primaryKey;autoIncrement;column:id"`
    Username  string    `gorm:"column:username;not null;uniqueIndex"`
    Email     string    `gorm:"column:email;not null;uniqueIndex"`
    Password  string    `gorm:"column:password;not null"`
    Token     string    `gorm:"column:token"`
    CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (User) TableName() string {
    return "users"
}