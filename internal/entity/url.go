package entity

import "time"

type Url struct {
    ID          int64     `gorm:"primaryKey;autoIncrement;column:id"`
    UserID      int64     `gorm:"column:user_id;not null"`
    ShortCode   string    `gorm:"column:short_code;not null;uniqueIndex"`
    OriginalUrl string    `gorm:"column:original_url;not null"`
    Hits        int64     `gorm:"column:hits;default:0"`
    CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt   time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (Url) TableName() string {
    return "urls"
}