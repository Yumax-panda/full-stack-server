package model

import (
	"github.com/google/uuid"
	_ "gorm.io/gorm"
)

// スキルの用途 e.g. フロントエンド, バックエンド, etc...
type Usage struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Skills []*Skill  `json:"skills" gorm:"many2many:skill_usage_relations;"`
}

// スキルのカテゴリ e.g. フレームワーク, ライブラリ, etc...
type Category struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Skills []*Skill  `json:"skills" gorm:"many2many:skill_category_relations;"`
}

// 技術 e.g. React, Vue, etc...
type Method struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	SkillID uuid.UUID `json:"skill_id"`
}

type Skill struct {
	ID         uuid.UUID   `json:"id"`
	UserID     uuid.UUID   `json:"user_id"`
	Level      uint8       `json:"level"`
	Usages     []*Usage    `json:"usages" gorm:"many2many:skill_usage_relations;"`
	Categories []*Category `json:"categories" gorm:"many2many:skill_category_relations;"`
	Method     Method      `json:"method"`
}

func GetSkillsByUserID(userID string) ([]Skill, error) {
	var skills []Skill
	err := db.Where("user_id = ?", userID).Find(&skills).Error
	return skills, err
}
