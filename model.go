package categories

import (
	"fmt"
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
	TenantId           int
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
	Categories         []CatgoriesOrd
	Assingedcategoryid int
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
	TenantId     int
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
	TenantId   int
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
	TenantId        int
}

type CategoryModel struct {
	Userid     int
	DataAccess int
}

var Categorymodel CategoryModel

// Parent Category List
func (categories CategoryModel) CategoryGroupList(offset int, limit int, filter Filter, createonly bool, DB *gorm.DB, tenantid int) (category []TblCategories, count int64, err error) {

	var categorycount int64

	query := DB.Table("tbl_categories").Where("is_deleted = 0 and parent_id=0 and (tenant_id is NULL or tenant_id = ?)", tenantid).Order("id desc")

	if createonly && categories.DataAccess == 1 {

		query = query.Where("tbl_categories.created_by = ? ", categories.Userid)
	}

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
func (categories CategoryModel) UpdateCategory(category *TblCategories, DB *gorm.DB, tenantid int) error {

	if category.ParentId == 0 && category.ImagePath == "" {

		if err := DB.Table("tbl_categories").Where("id = ? and (tenant_id is NULL or tenant_id = ?)", category.Id, tenantid).UpdateColumns(map[string]interface{}{"category_name": category.CategoryName, "category_slug": category.CategorySlug, "description": category.Description, "modified_by": category.ModifiedBy, "modified_on": category.ModifiedOn}).Error; err != nil {

			return err
		}
	} else {
		if err := DB.Table("tbl_categories").Where("id = ? and (tenant_id is NULL or tenant_id = ?)", category.Id, tenantid).UpdateColumns(map[string]interface{}{"category_name": category.CategoryName, "parent_id": category.ParentId, "category_slug": category.CategorySlug, "description": category.Description, "image_path": category.ImagePath, "modified_by": category.ModifiedBy, "modified_on": category.ModifiedOn}).Error; err != nil {

			return err
		}
	}

	return nil
}

// Children Category List
func (cate CategoryModel) GetCategoryList(categ CategoriesListReq, flag int, DB *gorm.DB, tenantid int) (categorylist []TblCategories, count int64) {

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
			query.Raw(` `+res+` select count(count(distinct(cat_tree.id))) from cat_tree where is_deleted = 0 and LOWER(TRIM(category_name)) ILIKE LOWER(TRIM(?)) group by cat_tree.id and (tenant_id is NULL or tenant_id = ?)`, "%"+categ.Keyword+"%", tenantid).Count(&categorycount)

			return categorylist, categorycount
		}

		query = query.Raw(` `+res+` select distinct(cat_tree.id),cat_tree.* from cat_tree where is_deleted = 0 `+selectGroupRemove+outerlevel+onlycategories+` and (tenant_id is NULL or tenant_id = ?) and LOWER(TRIM(category_name)) ILIKE LOWER(TRIM(?)) limit(?) offset(?) `, tenantid, "%"+categ.Keyword+"%", categ.Limit, categ.Offset)

	} else if flag == 0 {

		query = query.Raw(``+res+` SELECT distinct(cat_tree.id),cat_tree.* FROM cat_tree where is_deleted = 0 `+selectGroupRemove+outerlevel+onlycategories+`  and (tenant_id is NULL or tenant_id = ?) and id not in (?) order by id desc limit(?) offset(?) `, categ.ParentId, tenantid, categ.Limit, categ.Offset)

	} else if flag == 1 {

		query = query.Raw(``+res+` SELECT * FROM cat_tree where is_deleted = 0 and (tenant_id is NULL or tenant_id = ?) order by id desc `, categ.ParentId, tenantid)
	}

	if categ.Limit != 0 {

		query.Find(&categorylist)

		return categorylist, categorycount

	} else {

		DB.Raw(` `+res+` SELECT count(*) FROM cat_tree where is_deleted = 0 and id not in (?) and (tenant_id is NULL or tenant_id = ?)  group by cat_tree.id order by id desc`, categ.ParentId, categ.ParentId, tenantid).Count(&categorycount)

		return categorylist, categorycount
	}

}

/*getCategory Details*/
func (cate CategoryModel) GetCategoryById(categoryId int, DB *gorm.DB, tenantid int) (categorylist TblCategories, err error) {

	if err := DB.Table("tbl_categories").Where("is_deleted=0 and id= ? and (tenant_id is NULL or tenant_id = ?)", categoryId, tenantid).First(&categorylist).Error; err != nil {

		return TblCategories{}, err
	}
	return categorylist, nil
}

func (cate CategoryModel) GetCategoryTree(categoryID int, DB *gorm.DB, tenantid int) ([]TblCategories, error) {
	var categories []TblCategories
	err := DB.Debug().Raw(`
		WITH RECURSIVE cat_tree AS (
			SELECT id, 	category_name,
			category_slug,
			parent_id,
			created_on,
			modified_on,
			is_deleted
			FROM tbl_categories
			WHERE id = ? and (tenant_id is NULL or tenant_id =?)
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
	`, categoryID, tenantid).Scan(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (cate CategoryModel) DeleteallCategoryById(category *TblCategories, categoryId []int, spacecatid, tenantid int, DB *gorm.DB) error {

	// if err := DB.Table("tbl_spaces").Where("page_category_id = ? and (tenant_id is NULL or tenant_id = ?)", spacecatid, tenantid).Updates(TblCategories{IsDeleted: category.IsDeleted, DeletedOn: category.DeletedOn, DeletedBy: category.DeletedBy}).Error; err != nil {

	// 	return err

	// }

	// if err := DB.Table("tbl_jobs").Where("categories_id = ? and (tenant_id is NULL or tenant_id = ?)", spacecatid, tenantid).UpdateColumns(map[string]interface{}{"categories_id": 0}).Error; err != nil {

	// 	return err
	// }

	if err := DB.Table("tbl_categories").Where("id in(?) and (tenant_id is NULL or tenant_id = ?)", categoryId, tenantid).Updates(TblCategories{IsDeleted: category.IsDeleted, DeletedOn: category.DeletedOn, DeletedBy: category.DeletedBy}).Error; err != nil {

		return err

	}

	return nil
}

func (cate CategoryModel) DeleteCategoryById(category *TblCategories, categoryId, tenantid int, DB *gorm.DB) error {

	if err := DB.Table("tbl_categories").Where("id = ? and (tenant_id is NULL or tenant_id = ?)", categoryId).Updates(TblCategories{IsDeleted: category.IsDeleted, DeletedOn: category.DeletedOn, DeletedBy: category.DeletedBy}).Error; err != nil {

		return err

	}

	return nil
}

// Get Childern list
func (cate CategoryModel) GetCategoryDetails(id, tenantid int, DB *gorm.DB) (category TblCategories, err error) {

	if err := DB.Table("tbl_categories").Where("id=? and (tenant_id is NULL or tenant_id = ?)", id, tenantid).First(&category).Error; err != nil {

		return TblCategories{}, err
	}
	return category, nil

}

// Check category group name already exists
func (cate CategoryModel) CheckCategoryGroupName(category TblCategories, userid int, name string, DB *gorm.DB, tenantid int) error {

	if userid == 0 {

		if err := DB.Table("tbl_categories").Where("LOWER(TRIM(category_name))=LOWER(TRIM(?)) and is_deleted=0 and (tenant_id is NULL or tenant_id = ?)", name, tenantid).First(&category).Error; err != nil {

			return err
		}
	} else {

		if err := DB.Table("tbl_categories").Where("LOWER(TRIM(category_name))=LOWER(TRIM(?)) and id not in (?) and is_deleted=0 and (tenant_id is NULL or tenant_id = ?)", name, userid, tenantid).First(&category).Error; err != nil {

			return err
		}
	}

	return nil
}

func (cate CategoryModel) GetAllParentCategory(DB *gorm.DB, tenantid int) (categories []TblCategories, err error) {

	if err := DB.Table("tbl_categories").Where("parent_id=0 and is_deleted=0 and (tenant_id is NULL or tenant_id = ?)", tenantid).Find(&categories).Error; err != nil {

		return []TblCategories{}, err
	}
	return categories, nil
}

// Check sub category name already exists
func (cate CategoryModel) CheckSubCategoryName(categoryid []int, currentid int, name string, DB *gorm.DB, tenantid int) (category TblCategories, err error) {

	if len(categoryid) == 0 {

		if err := DB.Table("tbl_categories").Where("LOWER(TRIM(category_name))=LOWER(TRIM(?)) and is_deleted=0 and (tenant_id is NULL or tenant_id = ?)", name, tenantid).First(&category).Error; err != nil {

			return TblCategories{}, err
		}
	} else {

		if err := DB.Table("tbl_categories").Where("LOWER(TRIM(category_name))=LOWER(TRIM(?)) and id in (?) and id not in (?) and is_deleted=0 and (tenant_id is NULL or tenant_id = ?)", name, categoryid, currentid, tenantid).First(&category).Error; err != nil {

			return TblCategories{}, err
		}
	}

	return category, nil
}

// update imagepath
func (cate CategoryModel) UpdateImagePath(Imagepath string, tenantid int, DB *gorm.DB) error {

	if err := DB.Table("tbl_categories").Where("image_path=? and (tenant_id is NULL or tenant_id = ?)", Imagepath).UpdateColumns(map[string]interface{}{
		"image_path": ""}).Error; err != nil {

		return err
	}

	return nil

}

// Children Category List
func (cate CategoryModel) GetSubCategoryList(categories *[]TblCategories, offset int, limit, tenantid int, filter Filter, parent_id int, flag int, DB *gorm.DB) (categorylist *[]TblCategories, count int64) {

	var categorycount int64

	res := `WITH RECURSIVE cat_tree AS (
		SELECT id, category_name, category_slug,image_path, parent_id,created_on,modified_on,is_deleted
		FROM tbl_categories
		WHERE id = ? and (tenant_id is NULL or tenant_id=?)
		UNION ALL
		SELECT cat.id, cat.category_name, cat.category_slug, cat.image_path ,cat.parent_id,cat.created_on,cat.modified_on,
		cat.is_deleted
		FROM tbl_categories AS cat
		JOIN cat_tree ON cat.parent_id = cat_tree.id )`

	query := DB

	if filter.Keyword != "" {

		if limit == 0 {
			query.Raw(` `+res+` select count(*) from cat_tree where is_deleted = 0 and parent_id != 0 and LOWER(TRIM(category_name)) ILIKE LOWER(TRIM(?)) group by cat_tree.id and (tenant_id is NULL or tenant_id = ?) `, parent_id, "%"+filter.Keyword+"%", tenantid).Count(&categorycount)

			return categories, categorycount
		}
		query = query.Raw(` `+res+` select * from cat_tree where is_deleted = 0 and parent_id != 0 and LOWER(TRIM(category_name)) ILIKE LOWER(TRIM(?)) limit ? offset ?  and (tenant_id is NULL or tenant_id = ?) `, parent_id, "%"+filter.Keyword+"%", limit, offset, tenantid)
	} else if flag == 0 {
		query = query.Raw(``+res+` SELECT * FROM cat_tree where is_deleted = 0 and id not in (?) order by id desc limit ? offset ? and (tenant_id is NULL or tenant_id = ?)  `, parent_id, parent_id, limit, offset, tenantid)
	} else if flag == 1 {
		query = query.Raw(``+res+` SELECT * FROM cat_tree where is_deleted = 0 order by id desc  and (tenant_id is NULL or tenant_id = ?) `, parent_id, tenantid)
	}
	if limit != 0 {

		query.Find(&categories)

		return categories, categorycount

	} else {

		DB.Raw(` `+res+` SELECT count(*) FROM cat_tree where is_deleted = 0 and id not in (?)  group by cat_tree.id order by id desc and (tenant_id is NULL or tenant_id = ?) , `, parent_id, parent_id, tenantid).Count(&categorycount)

		return categories, categorycount
	}

}

func (cate CategoryModel) DeleteCategoryByIds(category *TblCategories, categoryId []int, tenantid int, DB *gorm.DB) error {

	if err := DB.Table("tbl_categories").Where("id in (?) and (tenant_id is NULL or tenant_id = ?)", categoryId, tenantid).Updates(TblCategories{IsDeleted: category.IsDeleted, DeletedOn: category.DeletedOn, DeletedBy: category.DeletedBy}).Error; err != nil {

		return err

	}

	return nil
}

// multiselect get entry category Id function
func (cate CategoryModel) MultiGetEntryCategoryids(entryCategory *TblChannelEntrie, channelId []int, tenantid int, DB *gorm.DB) (entries []string, categoryRowIds []string, entryRowIds []string, err error) {

	var (
		entryId       []string
		categoryRowId []string
		entryrowId    []string
	)

	result := DB.Raw(`with recursive categories as
					(
						select id, category_name, category_slug, parent_id, is_deleted
						from tbl_categories 
						where id in (?) and (tenant_id is NULL or tenant_id = ?)
						union all
						select TC.id,TC.category_name,TC.category_slug,TC.parent_id,TC.is_deleted
						from categories C 
						join tbl_categories TC on C.id = TC.parent_id
					)
					select id from categories`, channelId, tenantid).Scan(&categoryRowId)
	if result.Error != nil {
		return
	}

	result = DB.Table("tbl_channel_entries").Pluck("categories_id", &entryId)
	if result.Error != nil {
		return
	}

	result = DB.Table("tbl_channel_entries").Pluck("id", &entryrowId)
	if result.Error != nil {
		return
	}

	return entryId, categoryRowId, entryrowId, nil
}

// delete all categories
func (cate CategoryModel) DeleteallCategoryByIds(category *TblCategories, categoryId []int, spacecatid []int, tenantid int, DB *gorm.DB) error {

	if err := DB.Table("tbl_spaces").Where("page_category_id in (?) and (tenant_id is NULL or tenant_id = ?)", spacecatid, tenantid).Updates(TblCategories{IsDeleted: category.IsDeleted, DeletedOn: category.DeletedOn, DeletedBy: category.DeletedBy}).Error; err != nil {

		fmt.Println(err)
		// return err

	}

	if err := DB.Table("tbl_jobs").Where("categories_id in (?) and (tenant_id is NULL or tenant_id = ?)", spacecatid, tenantid).UpdateColumns(map[string]interface{}{"categories_id": 0}).Error; err != nil {

		return err
	}

	if err := DB.Table("tbl_categories").Where("id in(?) and (tenant_id is NULL or tenant_id = ?)", categoryId, tenantid).Updates(TblCategories{IsDeleted: category.IsDeleted, DeletedOn: category.DeletedOn, DeletedBy: category.DeletedBy}).Error; err != nil {

		return err

	}

	return nil
}

func (cate CategoryModel) GetChannelCategoryids(channelcategory *TblChannelCategorie, DB *gorm.DB) (categories []string, rowIds []string, err error) {

	// var categoryId string
	var (
		categoryId []string
		rowId      []string
	)

	result := DB.Table("tbl_channel_categories").Pluck("category_id", &categoryId)
	if result.Error != nil {
		return
	}

	result = DB.Table("tbl_channel_categories").Pluck("id", &rowId)
	if result.Error != nil {
		return
	}

	return categoryId, rowId, nil

}

func (cate CategoryModel) DeleteChannelCategoryids(channelCategory *TblChannelCategorie, channelId [][]int, rowId []string, categoryId int, DB *gorm.DB, tenantid int) error {

	for i := 0; i < len(channelId); i++ {
		for j := 0; j < len(channelId[i]); j++ {
			if categoryId == channelId[i][j] {
				result := DB.Debug().Where("id = ? and (tenant_id is NULL or tenant_id = ?)", rowId[i]).Delete(&TblChannelCategorie{})
				if result.Error != nil {
					return result.Error
				}
				break
			}
		}

	}

	return nil
}

func (cate CategoryModel) GetEntryCategoryids(entryCategory *TblChannelEntrie, channelId int, DB *gorm.DB, tenantid int) (entries []string, categoryRowIds []string, entryRowIds []string, err error) {

	var (
		entryId       []string
		categoryRowId []string
		entryrowId    []string
	)

	result := DB.Raw(`with recursive categories as
					(
						select id, category_name, category_slug, parent_id, is_deleted
						from tbl_categories 
						where id = ? and (tenant_id is NULL or tenant_id = ?)
						union all
						select TC.id,TC.category_name,TC.category_slug,TC.parent_id,TC.is_deleted
						from categories C 
						join tbl_categories TC on C.id = TC.parent_id
					)
					select id from categories`, channelId, tenantid).Scan(&categoryRowId)
	if result.Error != nil {
		return
	}

	result = DB.Debug().Table("tbl_channel_entries").Pluck("categories_id", &entryId)
	if result.Error != nil {
		return
	}

	result = DB.Debug().Table("tbl_channel_entries").Pluck("id", &entryrowId)
	if result.Error != nil {
		return
	}

	return entryId, categoryRowId, entryrowId, nil
}

// delete entry category ids

func (cate CategoryModel) DeleteEntryCategoryids(channelCategory *TblChannelEntrie, entryId string, rowId int, DB *gorm.DB, tenantid int) error {

	result := DB.Debug().Table("tbl_channel_entries").Where("id = ? and (tenant_id is NULL or tenant_id = ?)", rowId, tenantid).UpdateColumn("categories_id", entryId)

	if result.Error != nil {

		return result.Error
	}

	return nil
}

func (Cat CategoryModel) GetHierarchicalCategoriesMappedInEntries(categoryIds []string, categories *[]TblCategories, db *gorm.DB) (err error) {

	if err := db.Debug().Raw("WITH RECURSIVE CATHIERARCHY AS ( SELECT C.ID,C.DESCRIPTION,C.IMAGE_PATH,C.PARENT_ID,C.CATEGORY_SLUG,C.CATEGORY_NAME,C.CREATED_ON,C.CREATED_BY,C.MODIFIED_ON,C.MODIFIED_BY,C.IS_DELETED,C.DELETED_ON FROM TBL_CATEGORIES AS C WHERE C.IS_DELETED = 0 AND C.ID::TEXT IN (?) UNION SELECT TC.ID,TC.DESCRIPTION,TC.IMAGE_PATH,TC.PARENT_ID,TC.CATEGORY_SLUG,TC.CATEGORY_NAME,TC.CREATED_ON,TC.CREATED_BY,TC.MODIFIED_ON,TC.MODIFIED_BY,TC.IS_DELETED,TC.DELETED_ON FROM TBL_CATEGORIES AS TC INNER JOIN CATHIERARCHY AS CH ON CH.PARENT_ID = TC.ID WHERE TC.IS_DELETED = 0) SELECT * FROM CATHIERARCHY AS CH ORDER BY CH.PARENT_ID", categoryIds).Find(&categories).Error; err != nil {

		return err
	}

	return nil
}