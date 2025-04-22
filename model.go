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
	TenantId           string
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
	TenantId     string
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
	TenantId   string
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
	TenantId        string
}

type CategoryModel struct {
	Userid     int
	DataAccess int
}

var Categorymodel CategoryModel

// Parent Category List
func (categories CategoryModel) CategoryGroupList(offset int, limit int, filter Filter, createonly bool, DB *gorm.DB, tenantid string) (category []TblCategories, count int64, err error) {

	var categorycount int64

	query := DB.Table("tbl_categories").Where("is_deleted = 0 and parent_id=0 and  tenant_id = ?", tenantid).Order("tbl_categories.created_on desc")

	if createonly && categories.DataAccess == 1 {

		query = query.Where("tbl_categories.created_by = ? ", categories.Userid)
	}

	if filter.Keyword != "" {

		query = query.Where("LOWER(TRIM(category_name)) like LOWER(TRIM(?))", "%"+filter.Keyword+"%")
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
func (categories CategoryModel) UpdateCategory(category *TblCategories, DB *gorm.DB, tenantid string) error {

	if category.ParentId == 0 && category.ImagePath == "" {

		if err := DB.Table("tbl_categories").Where("id = ? and  tenant_id = ?", category.Id, tenantid).UpdateColumns(map[string]interface{}{"category_name": category.CategoryName, "category_slug": category.CategorySlug, "description": category.Description, "modified_by": category.ModifiedBy, "modified_on": category.ModifiedOn}).Error; err != nil {

			return err
		}
	} else {
		if err := DB.Table("tbl_categories").Where("id = ? and  tenant_id = ?", category.Id, tenantid).UpdateColumns(map[string]interface{}{"category_name": category.CategoryName, "parent_id": category.ParentId, "category_slug": category.CategorySlug, "description": category.Description, "image_path": category.ImagePath, "modified_by": category.ModifiedBy, "modified_on": category.ModifiedOn}).Error; err != nil {

			return err
		}
	}

	return nil
}

// Children Category List
func (cate CategoryModel) GetCategoryList(categ CategoriesListReq, flag int, DB *gorm.DB, tenantid string) (categorylist []TblCategories, count int64) {

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
			query.Raw(` `+res+` select count(count(distinct(cat_tree.id))) from cat_tree where is_deleted = 0 and LOWER(TRIM(category_name)) like LOWER(TRIM(?)) group by cat_tree.id and  tenant_id = ?`, "%"+categ.Keyword+"%", tenantid).Count(&categorycount)

			return categorylist, categorycount
		}

		query = query.Raw(` `+res+` select distinct(cat_tree.id),cat_tree.* from cat_tree where is_deleted = 0 `+selectGroupRemove+outerlevel+onlycategories+` and  tenant_id = ? and LOWER(TRIM(category_name)) like LOWER(TRIM(?)) limit(?) offset(?) `, tenantid, "%"+categ.Keyword+"%", categ.Limit, categ.Offset)

	} else if flag == 0 {

		query = query.Raw(``+res+` SELECT distinct(cat_tree.id),cat_tree.* FROM cat_tree where is_deleted = 0 `+selectGroupRemove+outerlevel+onlycategories+`  and  tenant_id = ? and id not in (?) order by id desc limit(?) offset(?) `, tenantid, categ.ParentId, categ.Limit, categ.Offset)

	} else if flag == 1 {

		query = query.Raw(``+res+` SELECT * FROM cat_tree where is_deleted = 0 and  tenant_id = ? order by id desc `, categ.ParentId, tenantid)
	}

	if categ.Limit != 0 {

		query.Find(&categorylist)

		return categorylist, categorycount

	} else {

		DB.Raw(` `+res+` SELECT count(*) FROM cat_tree where is_deleted = 0 and id not in (?) and  tenant_id = ?  group by cat_tree.id order by id desc`, categ.ParentId, categ.ParentId, tenantid).Count(&categorycount)

		return categorylist, categorycount
	}

}

/*getCategory Details*/
func (cate CategoryModel) GetCategoryById(categoryId int, DB *gorm.DB, tenantid string) (categorylist TblCategories, err error) {

	if err := DB.Table("tbl_categories").Where("is_deleted=0 and id= ? and  tenant_id = ?", categoryId, tenantid).First(&categorylist).Error; err != nil {

		return TblCategories{}, err
	}
	return categorylist, nil
}

func (cate CategoryModel) GetCategoryTree(categoryID int, DB *gorm.DB, tenantid string) ([]TblCategories, error) {
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
			WHERE id = ? and  tenant_id =?
			UNION ALL
			SELECT cat.id, cat.category_name,
			cat.category_slug,
			cat.parent_id,
			cat.created_on,
			cat.modified_on,
			cat.is_deleted
			FROM tbl_categories AS cat
			JOIN cat_tree ON cat.parent_id = cat_tree.id and  cat.tenant_id =?
		)
		SELECT *
		FROM cat_tree WHERE IS_DELETED = 0 order by id desc
	`, categoryID, tenantid, tenantid).Scan(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (cate CategoryModel) DeleteallCategoryById(category *TblCategories, categoryId []int, spacecatid int, tenantid string, DB *gorm.DB) error {

	// if err := DB.Table("tbl_spaces").Where("page_category_id = ? and  tenant_id = ?", spacecatid, tenantid).Updates(TblCategories{IsDeleted: category.IsDeleted, DeletedOn: category.DeletedOn, DeletedBy: category.DeletedBy}).Error; err != nil {

	// 	return err

	// }

	// if err := DB.Table("tbl_jobs").Where("categories_id = ? and  tenant_id = ?", spacecatid, tenantid).UpdateColumns(map[string]interface{}{"categories_id": 0}).Error; err != nil {

	// 	return err
	// }

	if err := DB.Table("tbl_categories").Where("id in(?) and  tenant_id = ?", categoryId, tenantid).Updates(TblCategories{IsDeleted: category.IsDeleted, DeletedOn: category.DeletedOn, DeletedBy: category.DeletedBy}).Error; err != nil {

		return err

	}

	return nil
}

func (cate CategoryModel) DeleteCategoryById(category *TblCategories, categoryId int, tenantid string, DB *gorm.DB) error {

	if err := DB.Table("tbl_categories").Where("id = ? and  tenant_id = ?", categoryId, tenantid).Updates(TblCategories{IsDeleted: category.IsDeleted, DeletedOn: category.DeletedOn, DeletedBy: category.DeletedBy}).Error; err != nil {

		return err

	}

	return nil
}

// Get Childern list
func (cate CategoryModel) GetCategoryDetails(id int, tenantid string, DB *gorm.DB) (category TblCategories, err error) {

	if err := DB.Table("tbl_categories").Where("id=? and  tenant_id = ?", id, tenantid).First(&category).Error; err != nil {

		return TblCategories{}, err
	}
	return category, nil

}

// Check category group name already exists
func (cate CategoryModel) CheckCategoryGroupName(category TblCategories, userid int, name string, DB *gorm.DB, tenantid string) error {

	if userid == 0 {

		if err := DB.Table("tbl_categories").Where("LOWER(TRIM(category_name))=LOWER(TRIM(?)) and is_deleted=0 and  tenant_id = ?", name, tenantid).First(&category).Error; err != nil {

			return err
		}
	} else {

		if err := DB.Table("tbl_categories").Where("LOWER(TRIM(category_name))=LOWER(TRIM(?)) and id not in (?) and is_deleted=0 and  tenant_id = ?", name, userid, tenantid).First(&category).Error; err != nil {

			return err
		}
	}

	return nil
}

func (cate CategoryModel) GetAllParentCategory(DB *gorm.DB, tenantid string) (categories []TblCategories, err error) {

	if err := DB.Table("tbl_categories").Where("parent_id=0 and is_deleted=0 and  tenant_id = ?", tenantid).Find(&categories).Error; err != nil {

		return []TblCategories{}, err
	}
	return categories, nil
}

// Check sub category name already exists
func (cate CategoryModel) CheckSubCategoryName(categoryid []int, currentid int, name string, DB *gorm.DB, tenantid string) (category TblCategories, err error) {

	if len(categoryid) == 0 {

		if err := DB.Table("tbl_categories").Where("LOWER(TRIM(category_name))=LOWER(TRIM(?)) and is_deleted=0 and  tenant_id = ? and parent_id=?", name, tenantid, currentid).First(&category).Error; err != nil {

			return TblCategories{}, err
		}
	} else {

		if err := DB.Table("tbl_categories").Where("LOWER(TRIM(category_name))=LOWER(TRIM(?)) and id in (?) and id not in (?) and is_deleted=0 and  tenant_id = ?", name, categoryid, currentid, tenantid).First(&category).Error; err != nil {

			return TblCategories{}, err
		}
	}

	return category, nil
}

// update imagepath
func (cate CategoryModel) UpdateImagePath(Imagepath string, tenantid string, DB *gorm.DB) error {

	if err := DB.Table("tbl_categories").Where("image_path=? and  tenant_id = ?", Imagepath, tenantid).UpdateColumns(map[string]interface{}{
		"image_path": ""}).Error; err != nil {

		return err
	}

	return nil

}

// Children Category List
func (cate CategoryModel) GetSubCategoryList(categories *[]TblCategories, offset int, limit int, tenantid string, filter Filter, parent_id int, flag int, DB *gorm.DB) (categorylist *[]TblCategories, count int64) {

	var categorycount int64

	res := `WITH RECURSIVE cat_tree AS (
		SELECT id, category_name, category_slug,image_path, parent_id,created_on,modified_on,is_deleted,tenant_id
		FROM tbl_categories
		WHERE id = ? and  tenant_id = ?
		UNION ALL
		SELECT cat.id, cat.category_name, cat.category_slug, cat.image_path ,cat.parent_id,cat.created_on,cat.modified_on,
		cat.is_deleted,cat.tenant_id
		FROM tbl_categories AS cat
		JOIN cat_tree ON cat.parent_id = cat_tree.id and cat.tenant_id = ?)`

	query := DB

	if filter.Keyword != "" {

		if limit == 0 {
			query.Raw(` `+res+` select count(*) from cat_tree where is_deleted = 0 and parent_id != 0 and LOWER(TRIM(category_name)) LIKE LOWER(TRIM(?)) group by cat_tree.id `, parent_id, tenantid, tenantid, "%"+filter.Keyword+"%").Count(&categorycount)

			return categories, categorycount
		}
		query = query.Raw(` `+res+` select * from cat_tree where is_deleted = 0 and parent_id != 0 and LOWER(TRIM(category_name)) LIKE LOWER(TRIM(?)) limit ? offset ?  `, parent_id, tenantid, tenantid, "%"+filter.Keyword+"%", limit, offset)
	} else if flag == 0 {
		query = query.Raw(``+res+` SELECT * FROM cat_tree where is_deleted = 0 and id not in (?) order by id desc limit ? offset ? `, parent_id, tenantid, tenantid, parent_id, limit, offset)
	} else if flag == 1 {
		query = query.Raw(``+res+` SELECT * FROM cat_tree where is_deleted = 0 order by id desc`, parent_id, tenantid, tenantid)
	}
	if limit != 0 {

		query.Find(&categories)

		return categories, categorycount

	} else {

		DB.Raw(` `+res+` SELECT count(*) FROM cat_tree where is_deleted = 0 and id not in (?)  group by cat_tree.id order by id desc`, parent_id, tenantid, tenantid, parent_id).Count(&categorycount)

		return categories, categorycount
	}

}

func (cate CategoryModel) DeleteCategoryByIds(category *TblCategories, categoryId []int, tenantid string, DB *gorm.DB) error {

	if err := DB.Table("tbl_categories").Where("id in (?) and  tenant_id = ?", categoryId, tenantid).Updates(TblCategories{IsDeleted: category.IsDeleted, DeletedOn: category.DeletedOn, DeletedBy: category.DeletedBy}).Error; err != nil {

		return err

	}

	return nil
}

// multiselect get entry category Id function
func (cate CategoryModel) MultiGetEntryCategoryids(entryCategory *TblChannelEntrie, channelId []int, tenantid string, DB *gorm.DB) (entries []string, categoryRowIds []string, entryRowIds []string, err error) {

	var (
		entryId       []string
		categoryRowId []string
		entryrowId    []string
	)

	result := DB.Raw(`with recursive categories as
					(
						select id, category_name, category_slug, parent_id, is_deleted
						from tbl_categories 
						where id in (?) and  tenant_id = ?
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
func (cate CategoryModel) DeleteallCategoryByIds(category *TblCategories, categoryId []int, spacecatid []int, tenantid string, DB *gorm.DB) error {

	// if err := DB.Table("tbl_spaces").Where("page_category_id in (?) and  tenant_id = ?", spacecatid, tenantid).Updates(TblCategories{IsDeleted: category.IsDeleted, DeletedOn: category.DeletedOn, DeletedBy: category.DeletedBy}).Error; err != nil {

	// 	fmt.Println(err)
	// 	// return err

	// }

	// if err := DB.Table("tbl_jobs").Where("categories_id in (?) and  tenant_id = ?", spacecatid, tenantid).UpdateColumns(map[string]interface{}{"categories_id": 0}).Error; err != nil {

	// 	return err
	// }

	if err := DB.Table("tbl_categories").Where("id in(?) and  tenant_id = ?", categoryId, tenantid).Updates(TblCategories{IsDeleted: category.IsDeleted, DeletedOn: category.DeletedOn, DeletedBy: category.DeletedBy}).Error; err != nil {

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

func (cate CategoryModel) DeleteChannelCategoryids(channelCategory *TblChannelCategorie, channelId [][]int, rowId []string, categoryId int, DB *gorm.DB, tenantid string) error {

	for i := 0; i < len(channelId); i++ {
		for j := 0; j < len(channelId[i]); j++ {
			if categoryId == channelId[i][j] {
				result := DB.Debug().Where("id = ? and  tenant_id = ?", rowId[i], tenantid).Delete(&TblChannelCategorie{})
				if result.Error != nil {
					return result.Error
				}
				break
			}
		}

	}

	return nil
}

func (cate CategoryModel) GetEntryCategoryids(entryCategory *TblChannelEntrie, channelId int, DB *gorm.DB, tenantid string) (entries []string, categoryRowIds []string, entryRowIds []string, err error) {

	var (
		entryId       []string
		categoryRowId []string
		entryrowId    []string
	)

	result := DB.Raw(`with recursive categories as
					(
						select id, category_name, category_slug, parent_id, is_deleted
						from tbl_categories 
						where id = ? and  tenant_id = ?
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

func (cate CategoryModel) DeleteEntryCategoryids(channelCategory *TblChannelEntrie, entryId string, rowId int, DB *gorm.DB, tenantid string) error {

	result := DB.Debug().Table("tbl_channel_entries").Where("id = ? and  tenant_id = ?", rowId, tenantid).UpdateColumn("categories_id", entryId)

	if result.Error != nil {

		return result.Error
	}

	return nil
}

func (Cat CategoryModel) GetHierarchicalCategoriesMappedInEntries(categoryIds []string, categories *[]TblCategories, db *gorm.DB, dbType string) (err error) {

	var condition string

	switch {

	case dbType == "postgres":

		condition = `C.ID::TEXT`

	case dbType == "mysql":

		condition = `CAST(C.ID AS CHAR)`
	}

	if err := db.Debug().Raw("WITH RECURSIVE CATHIERARCHY AS ( SELECT C.ID,C.DESCRIPTION,C.IMAGE_PATH,C.PARENT_ID,C.CATEGORY_SLUG,C.CATEGORY_NAME,C.CREATED_ON,C.CREATED_BY,C.MODIFIED_ON,C.MODIFIED_BY,C.IS_DELETED,C.DELETED_ON FROM tbl_categories AS C WHERE C.IS_DELETED = 0 AND "+condition+" IN (?) UNION SELECT TC.ID,TC.DESCRIPTION,TC.IMAGE_PATH,TC.PARENT_ID,TC.CATEGORY_SLUG,TC.CATEGORY_NAME,TC.CREATED_ON,TC.CREATED_BY,TC.MODIFIED_ON,TC.MODIFIED_BY,TC.IS_DELETED,TC.DELETED_ON FROM tbl_categories AS TC INNER JOIN CATHIERARCHY AS CH ON CH.PARENT_ID = TC.ID WHERE TC.IS_DELETED = 0) SELECT * FROM CATHIERARCHY AS CH ORDER BY CH.PARENT_ID", categoryIds).Find(&categories).Error; err != nil {

		return err
	}

	return nil
}

func (cat CategoryModel) FlexibleCategoryList(limit, offset, categoryGrpId, hierarchylevel int, tenantId string, excludeGroup, excludeParent, checkEntriesPresence, exactLevelOnly bool, channelSlug, categoryGrpSlug string, db *gorm.DB) (categories []TblCategories, count int64, err error) {

	var (
		hierarchyString, fromHierarchyString, selectHierarchyString                  string
		exactLevel, limitString, offsetString, createOnlyString, chanBasedCategories string
		categoryString, selectParentRemove, removeGroup, EntryMappedCategories       string
		convTenantId                                                                 = tenantId
	)

	if cat.DataAccess == 1 {

		createOnlyString = ` and ct.created_by = ` + strconv.Itoa(cat.Userid)
	}

	if excludeGroup {

		removeGroup = ` and ct.parent_id != 0`
	}

	if channelSlug != "" {

		var joinCondition string

		if db.Config.Dialector.Name() == "mysql" {

			joinCondition = ` find_in_set(ct.id,tcc.category_id) > 0`

		} else if db.Config.Dialector.Name() == "postgres" {

			joinCondition = ` ct.id = any(string_to_array(tcc.category_id,',')::Integer[])`
		}

		chanBasedCategories = `inner join tbl_channel_categories as tcc on ` + joinCondition + `inner join tbl_channels as tc on tc.id = tcc.channel_id and tc.slug_name ='` + channelSlug + `'`
	}

	if checkEntriesPresence {

		var joinCondition string

		if db.Config.Dialector.Name() == "mysql" {

			joinCondition = ` find_in_set(ct.id,ce.categories_id) > 0`

		} else if db.Config.Dialector.Name() == "postgres" {

			joinCondition = ` ct.id = any(string_to_array(ce.categories_id,',')::Integer[])`
		}

		EntryMappedCategories = ` inner join tbl_channel_entries as ce on ` + joinCondition + ` and ce.is_deleted = 0 and ce.status = 1 and ce.tenant_id = '` + convTenantId
	}

	if categoryGrpId != 0 {

		categoryString = ` id = ` + strconv.Itoa(categoryGrpId) + ` and `

		if excludeParent {

			selectParentRemove = ` And ct.id != ` + strconv.Itoa(categoryGrpId)
		}

	} else if categoryGrpSlug != "" {

		categoryString = ` id = (select id from tbl_categories where is_deleted = 0 and category_slug = '` + categoryGrpSlug + `' and tenant_id = '` + convTenantId + `') and `

		if excludeParent {

			selectParentRemove = ` And ct.id != (select id from tbl_categories where is_deleted = 0 and category_slug = '` + categoryGrpSlug + `' and tenant_id = '` + convTenantId + `')`
		}
	}

	if hierarchylevel != 0 {

		hierarchyString = ` and cat_tree.level < ` + strconv.Itoa(hierarchylevel)

		fromHierarchyString = `,cat_tree.level + 1`

		selectHierarchyString = `,1 as level`

	}

	if exactLevelOnly {

		exactLevel = ` and level = ` + strconv.Itoa(hierarchylevel)
	}

	if limit > 0 {

		limitString = `limit ` + strconv.Itoa(limit)
	}

	if offset > -1 {

		offsetString = ` offset ` + strconv.Itoa(offset)
	}

	res := `with recursive cat_tree AS (
	select id, category_name, category_slug, description,image_path, parent_id, created_on,modified_on, modified_by,is_deleted,tenant_id` + selectHierarchyString + ` from tbl_categories where ` + categoryString + ` tenant_id = '` + convTenantId + `' union select tbl_categories.id, tbl_categories.category_name, tbl_categories.category_slug,tbl_categories.description, tbl_categories.image_path, tbl_categories.parent_id, tbl_categories.created_on,tbl_categories.modified_on,tbl_categories.modified_by,tbl_categories.is_deleted,tbl_categories.tenant_id` + fromHierarchyString + ` from tbl_categories join cat_tree on tbl_categories.parent_id = cat_tree.id where tbl_categories.tenant_id = '` + convTenantId + `' ` + hierarchyString + ` )`

	if err := db.Debug().Raw(` ` + res + `select ct.* from cat_tree as ct ` + chanBasedCategories + EntryMappedCategories + ` where ct.is_deleted = 0  and ct.tenant_id ='` + convTenantId + `' ` + selectParentRemove + ` ` + exactLevel + ` ` + removeGroup + ` ` + createOnlyString + ` order by ct.id desc ` + limitString + offsetString).Find(&categories).Error; err != nil {

		return []TblCategories{}, 0, err
	}

	if err := db.Debug().Raw(` ` + res + `select count(distinct(ct.id)) from cat_tree as ct ` + chanBasedCategories + EntryMappedCategories + ` where ct.is_deleted = 0 and ct.tenant_id ='` + convTenantId + `' ` + selectParentRemove + ` ` + exactLevel + ` ` + removeGroup + ` ` + createOnlyString + ` group by ct.id order by ct.id desc`).Count(&count).Error; err != nil {

		return []TblCategories{}, 0, err
	}

	return categories, count, nil
}
