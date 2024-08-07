package mysql

import (
	"time"

	"gorm.io/gorm"
)

type TblCategory struct {
	Id           int       `gorm:"primaryKey;auto_increment"`
	CategoryName string    `gorm:"type:varchar(255)"`
	CategorySlug string    `gorm:"type:varchar(255)"`
	Description  string    `gorm:"type:varchar(255)"`
	ImagePath    string    `gorm:"type:varchar(255)"`
	ParentId     int       `gorm:"type:int"`
	CreatedOn    time.Time `gorm:"type:datetime"`
	CreatedBy    int       `gorm:"type:int"`
	ModifiedOn   time.Time `gorm:"type:datetime;DEFAULT:NULL"`
	ModifiedBy   int       `gorm:"DEFAULT:NULL;type:int"`
	IsDeleted    int       `gorm:"type:int"`
	DeletedOn    time.Time `gorm:"type:datetime;DEFAULT:NULL"`
	DeletedBy    int       `gorm:"type:int;DEFAULT:NULL;type:int"`
	TenantId     int       `gorm:"type:int;"`
}

type TblChannelCategory struct {
	Id         int       `gorm:"primaryKey;auto_increment"`
	ChannelId  int       `gorm:"type:int"`
	CategoryId string    `gorm:"type:varchar(255)"`
	CreatedAt  int       `gorm:"type:int"`
	CreatedOn  time.Time `gorm:"type:datetime"`
	TenantId   int       `gorm:"type:int;"`
}

// MigrateTable creates this package related tables in your database
func MigrationTables(db *gorm.DB) {

	if err := db.AutoMigrate(

		&TblCategory{},
	); err != nil {

		panic(err)
	}

	db.Exec(`INSERT INTO public.tbl_categories(id, category_name, category_slug, created_on, created_by,is_deleted, parent_id, description)	VALUES (1, 'Default Category', 'default_category', '2024-03-04 11:22:03', 1, 0, 0, 'Default_Category'),(2, 'Default1', 'default1', '2024-03-04 11:22:03', 1, 0, 1, 'Default1');`)

}