package model

//定义数据库模型（User、Course、Video、Order），使用 GORM 标签指定字段属性

import "time"

// 1. User 用户模型（没问题，无需改）
type User struct {
	ID       int64  `gorm:"primaryKey;autoIncrement"`
	Username string `gorm:"type:varchar(50);not null;unique"` // 唯一索引防重复注册
	Password string `gorm:"type:varchar(255);not null"`       // 加密存储
	Role     string `gorm:"type:varchar(20);not null"`        // admin/student
}

// 2. Course 课程模型（没问题，无需改）
type Course struct {
	ID          int64   `gorm:"primaryKey;autoIncrement"`
	Title       string  `gorm:"type:varchar(100);not null"`    // 课程标题（非空）
	Price       float64 `gorm:"type:decimal(10,2);not null"`   // 价格（decimal类型适合金额）
	Intro       string  `json:"intro"`                         // 课程简介（可选）
	CoverOssKey string  `gorm:"not null" json:"cover_oss_key"` // OSS封面存储路径（非空）
	CreatedBy   int64   `gorm:"not null" json:"created_by"`    // 创建者ID（关联User的ID）
}

// 3. Video 视频模型（关键修改：删掉 DeletedAt 字段）
type Video struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseID  int64     `gorm:"not null" json:"course_id"` // 关联Course的ID（非空）
	OssKey    string    `gorm:"not null" json:"oss_key"`   // OSS视频存储路径（非空）
	CreatedAt time.Time `json:"created_at"`
}

// 4. Order 订单模型
type Order struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"not null" json:"user_id"`    // 关联User的ID（非空）
	CourseID  int64     `gorm:"not null" json:"course_id"`  // 关联Course的ID（非空）
	Status    string    `gorm:"default:paid" json:"status"` // 订单状态（默认已支付）
	CreatedAt time.Time `json:"created_at"`
}
