package v1

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"


	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"

	"github.com/siangyeh8818/gin.project.template/pkg/env"
	"github.com/siangyeh8818/md5-hash-server/internal/app"
	"github.com/siangyeh8818/md5-hash-server/internal/service/tag_service"

	"github.com/siangyeh8818/gin.project.template/pkg/export"
	"github.com/siangyeh8818/gin.project.template/pkg/logging"
	"github.com/siangyeh8818/gin.project.template/pkg/setting"
	"github.com/siangyeh8818/gin.project.template/pkg/util"
)


func UploadUrlfile(c *gin.Context) {
	//var (
	//	appG = app.Gin{C: c}
	//)
	
	//form, _ := c.MultipartForm()
	//file, _, _ := c.Request.FormFile("file")
	//files := form.File["files"]
	//log.Println(len(files))
	//for i, file := range files {
	//	log.Println(i)
	//	log.Println(file.Filename)
//
//		c.SaveUploadedFile(file, "tmp/"+file.Filename)
//	}

	file, _ := c.FormFile("file")
	//if err != nil {
	//	c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
	//	return
	//}
	c.SaveUploadedFile(file, "/tmp/"+file.Filename)
	hashresult := md5.Sum([]byte(Readfile("/tmp/"+file.Filename)))
	//hashresult := md5.Sum([]byte(Readfile("/Users/a0000/Documents/siang-github/md5-hash-server/"+files[0].Filename)))
	log.Printf("%x",hashresult) 

	c.String(http.StatusOK, hex.EncodeToString(hashresult[:]), nil)
	//c.String(http.StatusOK, string(hashresult[:]), len(files))
	//appG.Response(http.StatusOK, env.ERROR_EXIST_TAG, hashresult)
	return

}

func GetTags(c *gin.Context) {
	appG := app.Gin{C: c}
	name := c.Query("name")
	state := -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	tagService := tag_service.Tag{
		Name:     name,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}
	tags, err := tagService.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, env.ERROR_GET_TAGS_FAIL, nil)
		return
	}

	count, err := tagService.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, env.ERROR_COUNT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, env.SUCCESS, map[string]interface{}{
		"lists": tags,
		"total": count,
	})
}

type AddTagForm struct {
	Name      string `form:"name" valid:"Required;MaxSize(100)"`
	CreatedBy string `form:"created_by" valid:"Required;MaxSize(100)"`
	State     int    `form:"state" valid:"Range(0,1)"`
}




func AddTag(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form AddTagForm
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != env.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	tagService := tag_service.Tag{
		Name:      form.Name,
		CreatedBy: form.CreatedBy,
		State:     form.State,
	}
	exists, err := tagService.ExistByName()
	if err != nil {
		appG.Response(http.StatusInternalServerError, env.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if exists {
		appG.Response(http.StatusOK, env.ERROR_EXIST_TAG, nil)
		return
	}

	err = tagService.Add()
	if err != nil {
		appG.Response(http.StatusInternalServerError, env.ERROR_ADD_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, env.SUCCESS, nil)
}

type EditTagForm struct {
	ID         int    `form:"id" valid:"Required;Min(1)"`
	Name       string `form:"name" valid:"Required;MaxSize(100)"`
	ModifiedBy string `form:"modified_by" valid:"Required;MaxSize(100)"`
	State      int    `form:"state" valid:"Range(0,1)"`
}

// @Summary Update article tag
// @Produce  json
// @Param id path int true "ID"
// @Param name body string true "Name"
// @Param state body int false "State"
// @Param modified_by body string true "ModifiedBy"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags/{id} [put]
func EditTag(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form = EditTagForm{ID: com.StrTo(c.Param("id")).MustInt()}
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != env.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	tagService := tag_service.Tag{
		ID:         form.ID,
		Name:       form.Name,
		ModifiedBy: form.ModifiedBy,
		State:      form.State,
	}

	exists, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, env.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, env.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = tagService.Edit()
	if err != nil {
		appG.Response(http.StatusInternalServerError, env.ERROR_EDIT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, env.SUCCESS, nil)
}

// @Summary Delete article tag
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags/{id} [delete]
func DeleteTag(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, env.INVALID_PARAMS, nil)
	}

	tagService := tag_service.Tag{ID: id}
	exists, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, env.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, env.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	if err := tagService.Delete(); err != nil {
		appG.Response(http.StatusInternalServerError, env.ERROR_DELETE_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, env.SUCCESS, nil)
}

// @Summary Export article tag
// @Produce  json
// @Param name body string false "Name"
// @Param state body int false "State"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags/export [post]
func ExportTag(c *gin.Context) {
	appG := app.Gin{C: c}
	name := c.PostForm("name")
	state := -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	tagService := tag_service.Tag{
		Name:  name,
		State: state,
	}

	filename, err := tagService.Export()
	if err != nil {
		appG.Response(http.StatusInternalServerError, env.ERROR_EXPORT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, env.SUCCESS, map[string]string{
		"export_url":      export.GetExcelFullUrl(filename),
		"export_save_url": export.GetExcelPath() + filename,
	})
}

// @Summary Import article tag
// @Produce  json
// @Param file body file true "Excel File"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags/import [post]
func ImportTag(c *gin.Context) {
	appG := app.Gin{C: c}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, env.ERROR, nil)
		return
	}

	tagService := tag_service.Tag{}
	err = tagService.Import(file)
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, env.ERROR_IMPORT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, env.SUCCESS, nil)
}

func Readfile(path string)  (string){
    f, err := ioutil.ReadFile(path)
    if err != nil {
        log.Println("read file ail", err)
    }
    return string(f)
}