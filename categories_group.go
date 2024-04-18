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

	_, Total_categories, _ := Categorymodel.CategoryGroupList(0, 0, filter, cat.DB)

	categorygrplist, _, cerr := Categorymodel.CategoryGroupList(offset, limit, filter, cat.DB)

	if cerr != nil {

		return []tblcategories{}, 0, cerr
	}

	return categorygrplist, Total_categories, nil

}

/*Add Category Group*/
func (cat *Categories) CreateCategoryGroup(req CategoryCreate) error {

	if AuthError := AuthandPermission(cat); AuthError != nil {

		return AuthError
	}

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

	err := Categorymodel.CreateCategory(&category, cat.DB)

	if err != nil {

		return err
	}

	return nil

}

/*UpdateCategoryGroup*/
func (cat *Categories) UpdateCategoryGroup(req CategoryCreate) error {

	if AuthError := AuthandPermission(cat); AuthError != nil {

		return AuthError
	}

	if req.Id <= 0 || req.CategoryName == "" {

		return ErrorCategoryName
	}
	var category tblcategories

	category.Id = req.Id

	category.CategoryName = req.CategoryName

	category.Description = req.Description

	category.CategorySlug = strings.ToLower(req.CategoryName)

	category.ModifiedBy = req.ModifiedBy

	category.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	err := Categorymodel.UpdateCategory(&category, cat.DB)

	if err != nil {

		return err
	}

	return nil

}

/*DeleteCategoryGroup*/
func (cat *Categories) DeleteCategoryGroup(Categoryid int, modifiedby int) error {

	if AuthError := AuthandPermission(cat); AuthError != nil {

		return AuthError
	}

	GetData, _ := Categorymodel.GetCategoryTree(Categoryid, cat.DB)

	var individualid []int

	for _, GetParent := range GetData {

		indivi := GetParent.Id

		individualid = append(individualid, indivi)
	}

	spacecategory := individualid[0]

	var category tblcategories

	category.DeletedBy = modifiedby

	category.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	category.IsDeleted = 1

	err := Categorymodel.DeleteallCategoryById(&category, individualid, spacecategory, cat.DB)

	if err != nil {

		return err
	}

	return nil

}
