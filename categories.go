package categories

import (
	"fmt"
	"strings"
	"time"
)

func CategoriesSetup(config Config) *Categories {

	MigrateTables(config.DB)

	return &Categories{
		DB:               config.DB,
		AuthEnable:       config.AuthEnable,
		PermissionEnable: config.PermissionEnable,
		Auth:             config.Auth,
		Permissions:      config.Permissions,
	}

}

/*ListCategory*/
func (cate *Categories) ListCategory(catlist CategoriesListReq) (tblcat []tblcategories, categories []tblcategories, parentcategory tblcategories, categorycount int64, err error) {

	if Autherr := AuthandPermission(cate); Autherr != nil {

		return
	}
	//get particular category
	parentcategory, err1 := Categorymodel.GetCategoryById(catlist.ParentId, cate.DB)

	if err1 != nil {

		fmt.Println(err)
	}

	//get overall count category
	_, count := Categorymodel.GetCategoryList(catlist, 0, cate.DB)

	childcategorys, _ := Categorymodel.GetCategoryList(catlist, 1, cate.DB)

	childcategory, _ := Categorymodel.GetCategoryList(catlist, 0, cate.DB)

	var AllCategorieswithSubCategories []Arrangecategories

	var id int

	if catlist.ParentId != 0 {

		id = catlist.ParentId

	} else if catlist.CategoryGroupId != 0 {

		id = catlist.CategoryGroupId
	}

	GetData, _ := Categorymodel.GetCategoryTree(id, cate.DB)

	var pid int

	for _, categories := range GetData {

		var addcat Arrangecategories

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
		FinalModalCategoryList = append(FinalModalCategoryList, infinalarray)
	}

	var FinalModalCategoriesList []tblcategories

	for index, val := range childcategorys {

		// var finalcat TblCategory

		finalcat := val

		for cindex, val2 := range FinalModalCategoryList {

			if index == cindex {

				for _, va3 := range val2.Categories {

					finalcat.Parent = append(finalcat.Parent, va3.Category)
				}
			}
		}
		FinalModalCategoriesList = append(FinalModalCategoriesList, finalcat)
	}
	var FinalCategoriesList []tblcategories

	for index, val := range childcategory {

		// var finalcat TblCategory

		finalcat := val

		for cindex, val2 := range FinalCategoryList {

			if index+catlist.Offset == cindex {

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

	var category tblcategories

	category.CategoryName = req.CategoryName

	category.CategorySlug = strings.ToLower(req.CategoryName)

	category.Description = req.Description

	category.ImagePath = req.ImagePath

	category.CreatedBy = req.CreatedBy

	category.ParentId = req.ParentId

	category.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	err := Categorymodel.CreateCategory(&category, cate.DB)

	if err != nil {

		return err
	}

	return nil

}

/*Update Sub category*/
func (cate *Categories) UpdateSubCategory(req CategoryCreate) error {

	if Autherr := AuthandPermission(cate); Autherr != nil {

		return Autherr
	}

	if req.Id <= 0 || req.CategoryName == "" {

		return ErrorCategoryName
	}

	var category tblcategories

	category.CategoryName = req.CategoryName

	category.CategorySlug = strings.ToLower(req.CategoryName)

	category.Description = req.Description

	category.ImagePath = req.ImagePath

	category.ParentId = req.ParentId

	category.CreatedBy = req.CreatedBy

	category.Id = req.Id

	category.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	category.ModifiedBy = req.ModifiedBy

	err := Categorymodel.UpdateCategory(&category, cate.DB)

	if err != nil {

		return err
	}

	return nil

}

/*Delete Sub Category*/
func (cate *Categories) DeleteSubCategory(categoryid int, modifiedby int) error {

	if Autherr := AuthandPermission(cate); Autherr != nil {

		return Autherr
	}

	var category tblcategories

	category.DeletedBy = modifiedby

	category.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	category.IsDeleted = 1

	err := Categorymodel.DeleteCategoryById(&category, categoryid, cate.DB)

	if err != nil {

		return err
	}

	return nil

}

// Get Sub Category List
func (cate *Categories) GetSubCategoryDetails(categoryid int) (categorys tblcategories, err error) {

	if Autherr := AuthandPermission(cate); Autherr != nil {

		return tblcategories{}, Autherr
	}

	category, err := Categorymodel.GetCategoryDetails(categoryid, cate.DB)

	if err != nil {

		return tblcategories{}, err
	}

	return category, nil

}

/*Remove entries cover image if media image delete*/
func (cate *Categories) UpdateImagePath(ImagePath string) error {

	if Autherr := AuthandPermission(cate); Autherr != nil {

		return Autherr
	}

	err := Categorymodel.UpdateImagePath(ImagePath, cate.DB)

	if err != nil {

		return err
	}

	return nil

}
