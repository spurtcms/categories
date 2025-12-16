package categories

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/spurtcms/categories/migration"
	"gorm.io/gorm"
)

func CategoriesSetup(config Config) *Categories {

	migration.AutoMigration(config.DB, config.DataBaseType)

	return &Categories{
		DB:               config.DB,
		AuthEnable:       config.AuthEnable,
		PermissionEnable: config.PermissionEnable,
		Auth:             config.Auth,
		Permissions:      config.Permissions,
	}

}

/*ListCategory*/
func (cate *Categories) ListCategory(offset int, limit int, filter Filter, parent_id int, tenantid string) (tblcat []TblCategories, categories []TblCategories, parentcategory TblCategories, categorycount int64, err error) {

	if Autherr := AuthandPermission(cate); Autherr != nil {

		return
	}
	//get particular category
	parentcategory, err1 := Categorymodel.GetCategoryById(parent_id, cate.DB, tenantid)

	if err1 != nil {

		fmt.Println(err)
	}

	//get overall count category

	var categorylist []TblCategories

	var categorylists []TblCategories

	var categorys []TblCategories

	if err1 != nil {
		fmt.Println(err)
	}
	_, count := Categorymodel.GetSubCategoryList(&categorylist, 0, 0, tenantid, filter, parent_id, 0, cate.DB)

	childcategorys, _ := Categorymodel.GetSubCategoryList(&categorys, offset, limit, tenantid, filter, parent_id, 1, cate.DB)

	childcategory, _ := Categorymodel.GetSubCategoryList(&categorylist, offset, limit, tenantid, filter, parent_id, 0, cate.DB)

	for _, val := range *childcategory {

		if !val.ModifiedOn.IsZero() {

			val.DateString = val.ModifiedOn.Format("02 Jan 2006 03:04 PM")

		} else {
			val.DateString = val.CreatedOn.Format("02 Jan 2006 03:04 PM")

		}

		categorylists = append(categorylists, val)

	}
	var AllCategorieswithSubCategories []Arrangecategories

	GetData, _ := Categorymodel.GetCategoryTree(parent_id, cate.DB, tenantid)

	var pid int

	for _, categories := range GetData {

		var addcat Arrangecategories

		addcat.Assingedcategoryid = categories.Id

		var individualid []CatgoriesOrd

		pid = categories.ParentId

	LOOP:
		for _, GetParent := range GetData {

			var indivi CatgoriesOrd

			if pid == GetParent.Id {

				pid = GetParent.ParentId

				indivi.Id = GetParent.Id

				indivi.Category = GetParent.CategoryName

				individualid = append(individualid, indivi)

				if pid != 0 {

					goto LOOP // this loop execute until parentid=0

				}
			}

		}

		var ReverseOrder Arrangecategories

		addcat.Categories = append(addcat.Categories, individualid...)

		var singlecat []CatgoriesOrd

		for i := len(addcat.Categories) - 1; i >= 0; i-- {

			var Sing CatgoriesOrd

			Sing.Id = addcat.Categories[i].Id

			Sing.Category = addcat.Categories[i].Category

			singlecat = append(singlecat, Sing)

		}

		var newcate CatgoriesOrd

		newcate.Id = categories.Id

		newcate.Category = categories.CategoryName

		addcat.Categories = append(addcat.Categories, newcate)

		singlecat = append(singlecat, newcate)

		ReverseOrder.Assingedcategoryid = addcat.Assingedcategoryid

		ReverseOrder.Categories = singlecat

		AllCategorieswithSubCategories = append(AllCategorieswithSubCategories, ReverseOrder)

	}

	var FinalCategoryList []Arrangecategories

	for _, val := range AllCategorieswithSubCategories {

		var infinalarray Arrangecategories

		for index, res := range val.Categories {

			if index < len(val.Categories)-1 {

				// var cate CatgoriesOrd

				cate := res

				infinalarray.Categories = append(infinalarray.Categories, cate)

			}

		}

		infinalarray.Assingedcategoryid = val.Assingedcategoryid

		FinalCategoryList = append(FinalCategoryList, infinalarray)
	}

	var FinalModalCategoryList []Arrangecategories

	for _, val := range AllCategorieswithSubCategories {

		var infinalarray Arrangecategories

		for index, res := range val.Categories {

			if index < len(val.Categories) {

				// var cate CatgoriesOrd

				cate := res

				infinalarray.Categories = append(infinalarray.Categories, cate)
			}
		}

		infinalarray.Assingedcategoryid = val.Assingedcategoryid

		FinalModalCategoryList = append(FinalModalCategoryList, infinalarray)
	}

	var FinalModalCategoriesList []TblCategories

	for index, val := range *childcategorys {

		// var finalcat TblCategory

		finalcat := val

		for cindex, val2 := range FinalModalCategoryList {

			if index+offset == cindex && val2.Assingedcategoryid == val.Id {

				for _, va3 := range val2.Categories {

					finalcat.Parent = append(finalcat.Parent, va3.Category)
				}
			}
		}
		FinalModalCategoriesList = append(FinalModalCategoriesList, finalcat)
	}

	var FinalCategoriesList []TblCategories

	for _, val := range categorylists {

		finalcat := val

		for _, val2 := range FinalCategoryList {

			if val2.Assingedcategoryid == val.Id {

				for _, va3 := range val2.Categories {

					finalcat.Parent = append(finalcat.Parent, va3.Category)
				}
			}

		}
		FinalCategoriesList = append(FinalCategoriesList, finalcat)
	}

	return FinalCategoriesList, FinalModalCategoriesList, parentcategory, count, nil
}

/*Add Category*/
func (cate *Categories) AddCategory(req CategoryCreate) error {

	if Autherr := AuthandPermission(cate); Autherr != nil {

		return Autherr
	}

	if req.CategoryName == "" {

		return ErrorCategoryName
	}

	var (
		category     TblCategories
		categorySlug string
	)

	category.CategoryName = req.CategoryName

	if req.CategorySlug == "" {
		categorySlug = strings.ToLower(strings.ReplaceAll(req.CategoryName, " ", "-"))
	} else {
		categorySlug = strings.ToLower(strings.ReplaceAll(req.CategorySlug, " ", "-"))
	}

	categorySlug = regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(categorySlug, "-")

	categorySlug = regexp.MustCompile(`-+`).ReplaceAllString(categorySlug, "-")

	categorySlug = strings.Trim(categorySlug, "-")

	category.CategorySlug = categorySlug

	category.Description = req.Description

	category.ImagePath = req.ImagePath

	category.CreatedBy = req.CreatedBy

	category.ParentId = req.ParentId

	category.TenantId = req.TenantId

	category.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	category.SeoTitle = req.SeoTitle

	category.SeoDescription = req.SeoDescription

	category.SeoKeyword = req.SeoKeyword

	err := Categorymodel.CreateCategory(&category, cate.DB)

	if err != nil {

		return err
	}

	return nil

}

/*Update Sub category*/
func (cate *Categories) UpdateSubCategory(req CategoryCreate, tenantid string) error {

	if Autherr := AuthandPermission(cate); Autherr != nil {

		return Autherr
	}

	if req.Id <= 0 || req.CategoryName == "" {

		return ErrorCategoryName
	}

	var (
		category     TblCategories
		categorySlug string
	)
	category.CategoryName = req.CategoryName
	if req.CategorySlug == "" {
		categorySlug = strings.ToLower(strings.ReplaceAll(req.CategoryName, " ", "-"))
	} else {
		categorySlug = strings.ToLower(strings.ReplaceAll(req.CategorySlug, " ", "-"))
	}
	categorySlug = regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(categorySlug, "-")
	categorySlug = regexp.MustCompile(`-+`).ReplaceAllString(categorySlug, "-")
	categorySlug = strings.Trim(categorySlug, "-")
	category.CategorySlug = categorySlug
	category.Description = req.Description
	category.ImagePath = req.ImagePath
	category.ParentId = req.ParentId
	category.CreatedBy = req.CreatedBy
	category.Id = req.Id
	category.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))
	category.ModifiedBy = req.ModifiedBy
	category.SeoTitle = req.SeoTitle
	category.SeoDescription = req.SeoDescription
	category.SeoKeyword = req.SeoKeyword

	err := Categorymodel.UpdateCategory(&category, cate.DB, tenantid)
	if err != nil {
		return err
	}

	return nil

}

/*Delete Sub Category*/
func (cate *Categories) DeleteSubCategory(categoryid int, modifiedby int, tenantid string) error {

	if Autherr := AuthandPermission(cate); Autherr != nil {
		return Autherr
	}

	if err := cate.DeleteChannelsubCategories(categoryid, tenantid); err != nil {

		fmt.Println(err)
	}

	if err := cate.DeleteEntriessubCategories(categoryid, tenantid); err != nil {

		fmt.Println(err)
	}

	var category TblCategories
	category.DeletedBy = modifiedby
	category.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))
	category.IsDeleted = 1
	err := Categorymodel.DeleteCategoryById(&category, categoryid, tenantid, cate.DB)

	if err != nil {
		return err
	}

	return nil

}

// Get Sub Category List
func (cate *Categories) GetSubCategoryDetailsBySlug(categoryslug string, tenantid string) (categorys TblCategories, err error) {

	if Autherr := AuthandPermission(cate); Autherr != nil {
		return TblCategories{}, Autherr
	}

	category, err := Categorymodel.GetCategoryDetailsBySlug(categoryslug, tenantid, cate.DB)

	if err != nil {
		return TblCategories{}, err
	}

	return category, nil

}

// Get Sub Category List
func (cate *Categories) GetSubCategoryDetails(categoryid int, tenantid string) (categorys TblCategories, err error) {

	if Autherr := AuthandPermission(cate); Autherr != nil {
		return TblCategories{}, Autherr
	}

	category, err := Categorymodel.GetCategoryDetails(categoryid, tenantid, cate.DB)

	if err != nil {
		return TblCategories{}, err
	}

	return category, nil

}

/*Remove entries cover image if media image delete*/
func (cate *Categories) UpdateImagePath(ImagePath string, tenantid string) error {

	if Autherr := AuthandPermission(cate); Autherr != nil {
		return Autherr
	}

	err := Categorymodel.UpdateImagePath(ImagePath, tenantid, cate.DB)
	if err != nil {
		return err
	}

	return nil

}

// multiSelect channel delete category ids
func (cate CategoryModel) MultiDeleteChannelCategoryids(channelCategory *TblChannelCategorie, channelIds [][]int, rowId []string, categoryIds []int, DB *gorm.DB) error {

	for _, categoryId := range categoryIds {
		for i, channelId := range channelIds {
			for _, id := range channelId {
				if id == categoryId {
					result := DB.Debug().Where("id = ?", rowId[i]).Delete(&TblChannelCategorie{})
					if result.Error != nil {
						return result.Error
					}
					break
				}
			}
		}
	}

	return nil
}

func (cat *Categories) CategoryList(limit, offset, categoryGrpId, hierarchyLevel int, tenantId string, checkEntriesPresence, excludeGroup, excludeParent, exactLevelOnly bool, channelSlug, categoryGrpSlug string) (CategoryList []TblCategories, Count int, err error) {

	categories, count, err := Categorymodel.FlexibleCategoryList(limit, offset, categoryGrpId, hierarchyLevel, tenantId, excludeGroup, excludeParent, checkEntriesPresence, exactLevelOnly, channelSlug, categoryGrpSlug, cat.DB)

	if err != nil {

		return []TblCategories{}, 0, err
	}

	var FinalCategoriesList []TblCategories

	seenCategory := make(map[int]bool)

	for _, category := range categories {

		if !seenCategory[category.Id] {

			FinalCategoriesList = append(FinalCategoriesList, category)

			seenCategory[category.Id] = true
		}
	}

	return FinalCategoriesList, int(count), nil

}
