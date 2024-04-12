package categories

import (
	"strings"
	"time"
)

/*List Category Group*/
func (cat *Categories) CategoryGroupList(limit int, offset int, filter Filter) (Categorylist []tblcategories, categorycount int64, err error) {

	if AuthError := AuthandPermission(cat); AuthError != nil {

		return []tblcategories{}, 0, AuthError
	}

	_, Total_categories, _ := C.GetCategoryList(0, 0, filter, cat.DB)

	categorygrplist, _, cerr := C.GetCategoryList(offset, limit, filter, cat.DB)

	if cerr != nil {

		return []tblcategories{}, 0, cerr
	}

	return categorygrplist, Total_categories, nil

}

/*Add Category Group*/
func (cat *Categories) CreateCategoryGroup(req CategoryCreate) error {

	if req.CategoryName == "" {

		return ErrorCategoryName
	}

	var category tblcategories

	category.CategoryName = req.CategoryName

	category.CategorySlug = strings.ToLower(req.CategoryName)

	category.Description = req.Description

	category.CreatedBy = req.CreatedBy

	category.ParentId = 0

	category.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	err := C.CreateCategory(category, cat.DB)

	if err != nil {

		return err
	}

	return nil

}
