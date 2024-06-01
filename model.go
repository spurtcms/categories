package categories

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

type TblCategories struct {
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
	ModifiedBy   int
}

type CategoriesListReq struct {
	Limit           int
	Offset          int
	Keyword         string
	Category        string
	ParentId        int  //get particular parentid below all child categories
	CategoryGroupId int  //get particular group categories only
	HirerachyLevel  int  //how many level of children you want default zero get all
	OnlyCategories  bool //if you want categories only without group enable true
}

type TblChannelCategorie struct {
	Id         int
	ChannelId  int
	CategoryId string
	CreatedAt  int
	CreatedOn  time.Time
}

type TblChannelEntrie struct {
	Id              int
	Title           string
	Slug            string
	Description     string
	UserId          int
	ChannelId       int
	Status          int //0-draft 1-publish 2-unpublish
	CoverImage      string
	ThumbnailImage  string
	MetaTitle       string
	MetaDescription string
	Keyword         string
	CategoriesId    string
	RelatedArticles string
	Feature         int
	ViewCount       int
	CreateTime      time.Time
	PublishedTime   time.Time
	ImageAltTag     string
	Author          string
	SortOrder       int
	Excerpt         string
	ReadingTime     int
	Tags            string
	CreatedOn       time.Time
	CreatedBy       int
	ModifiedBy      int
	ModifiedOn      time.Time
	IsActive        int
	IsDeleted       int
	DeletedOn       time.Time
	DeletedBy       int
}

type CategoryModel struct{}

var Categorymodel CategoryModel

// Parent Category List
func (categories CategoryModel) CategoryGroupList(offset int, limit int, filter Filter, DB *gorm.DB) (category []TblCategories, count int64, err error) {

	var categorycount int64

	query := DB.Table("tbl_categories").Where("is_deleted = 0 and parent_id=0").Order("id desc")

	if filter.Keyword != "" {

		query = query.Where("LOWER(TRIM(category_name)) ILIKE LOWER(TRIM(?))", "%"+filter.Keyword+"%")
	}

	if limit != 0 {

		query.Limit(limit).Offset(offset).Find(&category)

		return category, categorycount, nil

	}

	query.Find(&category).Count(&categorycount)

	if query.Error != nil {

		return []TblCategories{}, 0, query.Error
	}

	return category, categorycount, nil

}

func (categories CategoryModel) CreateCategory(category *TblCategories, DB *gorm.DB) error {

	if err := DB.Table("tbl_categories").Create(&category).Error; err != nil {

		return err
	}

	return nil
}

// Update Children list
func (categories CategoryModel) UpdateCategory(category *TblCategories, DB *gorm.DB) error {

	if category.ParentId == 0 && category.ImagePath == "" {

		if err := DB.Table("tbl_categories").Where("id = ?", category.Id).UpdateColumns(map[string]interface{}{"category_name": category.CategoryName, "category_slug": category.CategorySlug, "description": category.Description, "modified_by": category.ModifiedBy, "modified_on": category.ModifiedOn}).Error; err != nil {

			return err
		}
	} else {
		if err := DB.Table("tbl_categories").Where("id = ?", category.Id).UpdateColumns(map[string]interface{}{"category_name": category.CategoryName, "parent_id": category.ParentId, "category_slug": category.CategorySlug, "description": category.Description, "image_path": category.ImagePath, "modified_by": category.ModifiedBy, "modified_on": category.ModifiedOn}).Error; err != nil {

			return err
		}
	}

	return nil
}

// Children Category List
func (cate CategoryModel) GetCategoryList(categ CategoriesListReq, flag int, DB *gorm.DB) (categorylist []TblCategories, count int64) {

	var categorycount int64

	var id int

	if categ.ParentId != 0 {

		id = categ.ParentId

	} else if categ.CategoryGroupId != 0 {

		id = categ.CategoryGroupId
	}

	category_string := ""

	selectGroupRemove := ""

	if categ.CategoryGroupId != 0 {

		category_string = `WHERE id = ` + strconv.Itoa(id)

		selectGroupRemove = `AND id != ` + strconv.Itoa(id)
	}

	onlycategories := ""

	if categ.OnlyCategories {

		onlycategories = `and parent_id!=0`
	}

	hierarchy_string := ""

	fromhierarchy_string := ""

	selecthierarchy_string := ""

	outerlevel := ""

	if categ.HirerachyLevel != 0 {

		hierarchy_string = ` WHERE CAT_TREE.LEVEL < ` + strconv.Itoa(categ.HirerachyLevel)

		fromhierarchy_string = `,CAT_TREE.LEVEL + 1`

		selecthierarchy_string = `,0 AS LEVEL`

		outerlevel = ` and level = ` + strconv.Itoa(categ.HirerachyLevel)

	}

	res := `WITH RECURSIVE cat_tree AS (
		SELECT id, category_name, category_slug,image_path, parent_id,created_on,modified_on,is_deleted ` + selecthierarchy_string + `
		FROM tbl_categories ` + category_string + `
		UNION
		SELECT cat.id, cat.category_name, cat.category_slug, cat.image_path ,cat.parent_id,cat.created_on,cat.modified_on,
		cat.is_deleted ` + fromhierarchy_string + `
		FROM tbl_categories AS cat
		JOIN cat_tree ON cat.parent_id = cat_tree.id  ` + hierarchy_string + `)`

	query := DB

	if categ.Keyword != "" {

		if categ.Limit == 0 {
			query.Raw(` `+res+` select count(count(distinct(cat_tree.id))) from cat_tree where is_deleted = 0 and LOWER(TRIM(category_name)) ILIKE LOWER(TRIM(?)) group by cat_tree.id `, "%"+categ.Keyword+"%").Count(&categorycount)

			return categorylist, categorycount
		}

		query = query.Raw(` `+res+` select distinct(cat_tree.id),cat_tree.* from cat_tree where is_deleted = 0 `+selectGroupRemove+outerlevel+onlycategories+` and LOWER(TRIM(category_name)) ILIKE LOWER(TRIM(?)) limit(?) offset(?) `, "%"+categ.Keyword+"%", categ.Limit, categ.Offset)

	} else if flag == 0 {

		query = query.Raw(``+res+` SELECT distinct(cat_tree.id),cat_tree.* FROM cat_tree where is_deleted = 0 `+selectGroupRemove+outerlevel+onlycategories+`  and id not in (?) order by id desc limit(?) offset(?) `, categ.ParentId, categ.Limit, categ.Offset)

	} else if flag == 1 {

		query = query.Raw(``+res+` SELECT * FROM cat_tree where is_deleted = 0 order by id desc `, categ.ParentId)
	}

	if categ.Limit != 0 {

		query.Find(&categorylist)

		return categorylist, categorycount

	} else {

		DB.Raw(` `+res+` SELECT count(*) FROM cat_tree where is_deleted = 0 and id not in (?)  group by cat_tree.id order by id desc`, categ.ParentId, categ.ParentId).Count(&categorycount)

		return categorylist, categorycount
	}

}

/*getCategory Details*/
func (cate CategoryModel) GetCategoryById(categoryId int, DB *gorm.DB) (categorylist TblCategories, err error) {

	if err := DB.Table("tbl_categories").Where("is_deleted=0 and id=?", categoryId).First(&categorylist).Error; err != nil {

		return TblCategories{}, err
	}
	return categorylist, nil
}

func (cate CategoryModel) GetCategoryTree(categoryID int, DB *gorm.DB) ([]TblCategories, error) {
	var categories []TblCategories
	err := DB.Raw(`
		WITH RECURSIVE cat_tree AS (
			SELECT id, 	category_name,
			category_slug,
			parent_id,
			created_on,
			modified_on,
			is_deleted
			FROM tbl_categories
			WHERE id = ?
			UNION ALL
			SELECT cat.id, cat.category_name,
			cat.category_slug,
			cat.parent_id,
			cat.created_on,
			cat.modified_on,
			cat.is_deleted
			FROM tbl_categories AS cat
			JOIN cat_tree ON cat.parent_id = cat_tree.id
		)
		SELECT *
		FROM cat_tree WHERE IS_DELETED = 0 order by id desc
	`, categoryID).Scan(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (cate CategoryModel) DeleteallCategoryById(category *TblCategories, categoryId []int, spacecatid int, DB *gorm.DB) error {

	if err := DB.Table("tbl_spaces").Where("page_category_id", spacecatid).Updates(TblCategories{IsDeleted: category.IsDeleted, DeletedOn: category.DeletedOn, DeletedBy: category.DeletedBy}).Error; err != nil {

		return err

	}

	if err := DB.Table("tbl_categories").Where("id in(?)", categoryId).Updates(TblCategories{IsDeleted: category.IsDeleted, DeletedOn: category.DeletedOn, DeletedBy: category.DeletedBy}).Error; err != nil {

		return err

	}

	return nil
}

func (cate CategoryModel) DeleteCategoryById(category *TblCategories, categoryId int, DB *gorm.DB) error {

	if err := DB.Table("tbl_categories").Where("id=?", categoryId).Updates(TblCategories{IsDeleted: category.IsDeleted, DeletedOn: category.DeletedOn, DeletedBy: category.DeletedBy}).Error; err != nil {

		return err

	}

	return nil
}

// Get Childern list
func (cate CategoryModel) GetCategoryDetails(id int, DB *gorm.DB) (category TblCategories, err error) {

	if err := DB.Table("tbl_categories").Where("id=?", id).First(&category).Error; err != nil {

		return TblCategories{}, err
	}
	return category, nil

}

// Check category group name already exists
func (cate CategoryModel) CheckCategoryGroupName(category TblCategories, userid int, name string, DB *gorm.DB) error {

	if userid == 0 {

		if err := DB.Table("tbl_categories").Where("LOWER(TRIM(category_name))=LOWER(TRIM(?)) and is_deleted=0", name).First(&category).Error; err != nil {

			return err
		}
	} else {

		if err := DB.Table("tbl_categories").Where("LOWER(TRIM(category_name))=LOWER(TRIM(?)) and id not in (?) and is_deleted=0", name, userid).First(&category).Error; err != nil {

			return err
		}
	}

	return nil
}

func (cate CategoryModel) GetAllParentCategory(DB *gorm.DB) (categories []TblCategories, err error) {

	if err := DB.Table("tbl_categories").Where("parent_id=0 and is_deleted=0").Find(&categories).Error; err != nil {

		return []TblCategories{}, err
	}
	return categories, nil
}

// Check sub category name already exists
func (cate CategoryModel) CheckSubCategoryName(categoryid []int, currentid int, name string, DB *gorm.DB) (category TblCategories, err error) {

	if len(categoryid) == 0 {

		if err := DB.Table("tbl_categories").Where("LOWER(TRIM(category_name))=LOWER(TRIM(?)) and is_deleted=0", name).First(&category).Error; err != nil {

			return TblCategories{}, err
		}
	} else {

		if err := DB.Table("tbl_categories").Where("LOWER(TRIM(category_name))=LOWER(TRIM(?)) and id in (?) and id not in (?) and is_deleted=0", name, categoryid, currentid).First(&category).Error; err != nil {

			return TblCategories{}, err
		}
	}

	return category, nil
}

// update imagepath
func (cate CategoryModel) UpdateImagePath(Imagepath string, DB *gorm.DB) error {

	if err := DB.Table("tbl_categories").Where("image_path=?", Imagepath).UpdateColumns(map[string]interface{}{
		"image_path": ""}).Error; err != nil {

		return err
	}

	return nil

}


// Children Category List
func (cate CategoryModel) GetSubCategoryList(categories *[]TblCategories, offset int, limit int, filter Filter, parent_id int, flag int, DB *gorm.DB) (categorylist *[]TblCategories, count int64) {

	var categorycount int64

	res := `WITH RECURSIVE cat_tree AS (
		SELECT id, category_name, category_slug,image_path, parent_id,created_on,modified_on,is_deleted
		FROM tbl_categories
		WHERE id = ?
		UNION ALL
		SELECT cat.id, cat.category_name, cat.category_slug, cat.image_path ,cat.parent_id,cat.created_on,cat.modified_on,
		cat.is_deleted
		FROM tbl_categories AS cat
		JOIN cat_tree ON cat.parent_id = cat_tree.id )`

	query := DB

	if filter.Keyword != "" {

		if limit == 0 {
			query.Raw(` `+res+` select count(*) from cat_tree where is_deleted = 0 and LOWER(TRIM(category_name)) ILIKE LOWER(TRIM(?)) group by cat_tree.id `, parent_id, "%"+filter.Keyword+"%").Count(&categorycount)

			return categories, categorycount
		}
		query = query.Raw(` `+res+` select * from cat_tree where is_deleted = 0 and LOWER(TRIM(category_name)) ILIKE LOWER(TRIM(?)) limit ? offset ? `, parent_id, "%"+filter.Keyword+"%", limit, offset)
	} else if flag == 0 {
		query = query.Raw(``+res+` SELECT * FROM cat_tree where is_deleted = 0 and id not in (?) order by id desc limit ? offset ? `, parent_id, parent_id, limit, offset)
	} else if flag == 1 {
		query = query.Raw(``+res+` SELECT * FROM cat_tree where is_deleted = 0 order by id desc `, parent_id)
	}
	if limit != 0 {

		query.Find(&categories)

		return categories, categorycount

	} else {

		DB.Raw(` `+res+` SELECT count(*) FROM cat_tree where is_deleted = 0 and id not in (?)  group by cat_tree.id order by id desc`, parent_id, parent_id).Count(&categorycount)

		return categories, categorycount
	}

}