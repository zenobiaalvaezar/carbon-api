package models

type TreeCategory struct {
	ID   int    `gorm:"primaryKey" json:"id"`
	Name string `gorm:"type:varchar(255);unique;not null" json:"name"`
}

//// TableName overrides the table name for TreeCategory
//func (TreeCategory) TableName() string {
//	return "tree_categories"
//}
