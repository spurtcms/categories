package categories

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
func MigrateTables(db *gorm.DB) {

	db.AutoMigrate(
		&TblCategories{},
	)
}
