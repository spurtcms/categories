package categories

import (
	"time"

	"gorm.io/gorm"
)

type tblcategories struct {
	Id                 int
	CategoryName       string
	CategorySlug       string
	Description        string
	ImagePath          string
	CreatedOn          time.Time
	CreatedBy          int
	ModifiedOn         time.Time `gorm:"DEFAULT:NULL"`
	ModifiedBy         int       `gorm:"DEFAULT:NULL"`
	IsDeleted          int
	DeletedOn          time.Time `gorm:"DEFAULT:NULL"`
	DeletedBy          int       `gorm:"DEFAULT:NULL"`
	ParentId           int
	CreatedDate        string   `gorm:"-"`
	ModifiedDate       string   `gorm:"-"`
	DateString         string   `gorm:"-"`
	ParentCategoryName string   `gorm:"-"`
	Parent             []string `gorm:"-"`
	ParentWithChild    []Result `gorm:"-"`
}

type Filter struct {
	Keyword  string
	Category string
	Status   string
	FromDate string
	ToDate   string
}

type Result struct {
	CategoryName string
}

type Arrangecategories struct {
	Categories []CatgoriesOrd
}

type CatgoriesOrd struct {
	Id       int
	Category string
}

type CategoryCreate struct {
	Id           int
	CategoryName string
	CategorySlug string
	Description  string
	ImagePath    string
	ParentId     int
	CreatedBy    int
}

type CategoryModel struct{}

var C CategoryModel

// Parent Category List
func (categories CategoryModel) GetCategoryList(offset int, limit int, filter Filter, DB *gorm.DB) (category []tblcategories, count int64, err error) {

	var categorycount int64

	query := DB.Table("tbl_categories").Where("is_deleted = 0 and parent_id=0").Order("id desc")

	if filter.Keyword != "" {

		query = query.Where("LOWER(TRIM(category_name)) ILIKE LOWER(TRIM(?))", "%"+filter.Keyword+"%")
	}

	if limit != 0 {

		query.Limit(limit).Offset(offset).Find(&categories)

		return category, categorycount, nil

	}

	query.Find(&categories).Count(&categorycount)

	if query.Error != nil {

		return []tblcategories{}, 0, query.Error
	}

	return category, categorycount, nil

}

func (categories CategoryModel) CreateCategory(category tblcategories, DB *gorm.DB) error {

	if err := DB.Table("tbl_categories").Create(&category).Error; err != nil {

		return err
	}

	return nil
}
