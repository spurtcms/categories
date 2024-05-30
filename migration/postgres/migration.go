package postgres

import (
	"time"

	"gorm.io/gorm"
)

type TblCategories struct {
	Id           int       `gorm:"primaryKey;auto_increment;type:serial"`
	CategoryName string    `gorm:"type:character varying"`
	CategorySlug string    `gorm:"type:character varying"`
	Description  string    `gorm:"type:character varying"`
	ImagePath    string    `gorm:"type:character varying"`
	ParentId     int       `gorm:"type:integer"`
	CreatedOn    time.Time `gorm:"type:timestamp without time zone"`
	CreatedBy    int       `gorm:"type:integer"`
	ModifiedOn   time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	ModifiedBy   int       `gorm:"DEFAULT:NULL"`
	IsDeleted    int       `gorm:"type:integer"`
	DeletedOn    time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	DeletedBy    int       `gorm:"DEFAULT:NULL;type:integer"`
}

// MigrateTable creates this package related tables in your database
func MigrationTables(db *gorm.DB) {

	if err := db.AutoMigrate(

		&TblCategories{},
	); err != nil {

		panic(err)
	}

	db.Exec(`INSERT INTO public.tbl_categories(id, category_name, category_slug, created_on, created_by,is_deleted, parent_id, description)	VALUES (1, 'Default Category', 'default_category', '2024-03-04 11:22:03', 1, 0, 0, 'Default_Category'),(2, 'Default1', 'default1', '2024-03-04 11:22:03', 1, 0, 1, 'Default1');`)

}
