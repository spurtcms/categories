package categories

// Check category name already exists
func (cate *Categories) CheckCategroyGroupName(id int, name string) (bool, error) {

	var category Tblcategories

	err := Categorymodel.CheckCategoryGroupName(category, id, name, cate.DB)

	if err != nil {

		return false, err

	}

	return true, nil
}

/*Get All cateogry with parents and subcategory*/
func (cate *Categories) AllCategoriesWithSubList() (arrangecategories []Arrangecategories, CategoryNames []string) {

	getallparentcat, _ := Categorymodel.GetAllParentCategory(cate.DB)

	var AllCategorieswithSubCategories []Arrangecategories

	for _, Group := range getallparentcat {

		GetData, _ := Categorymodel.GetCategoryTree(Group.Id, cate.DB)

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

						goto LOOP //this loop get looped until parentid=0

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

	}

	/*This for Channel category show also individual group*/
	var FinalCategoryList []Arrangecategories

	for _, val := range AllCategorieswithSubCategories {

		if len(val.Categories) > 1 {

			var infinalarray Arrangecategories

			infinalarray.Categories = append(infinalarray.Categories, val.Categories...)

			FinalCategoryList = append(FinalCategoryList, infinalarray)
		}

	}

	var Categorynames []string

	for _, val := range FinalCategoryList {

		var name string

		for index, cat := range val.Categories {

			if len(val.Categories)-1 == index {

				name += cat.Category
			} else {
				name += cat.Category + " / "
			}

		}

		Categorynames = append(Categorynames, name)

	}

	return FinalCategoryList, Categorynames

}

// Check Sub category name already exists
func (cate *Categories) CheckSubCategroyName(id []int, Currentcategoryid int, name string) (bool, error) {

	if Autherr := AuthandPermission(cate); Autherr != nil {

		return false, Autherr
	}

	_, err := Categorymodel.CheckSubCategoryName(id, Currentcategoryid, name, cate.DB)

	if err != nil {

		return false, err
	}

	return true, nil
}
