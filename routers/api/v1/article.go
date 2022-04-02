package v1

import (
	"gindemo/models"
	"gindemo/pkg/app"
	"gindemo/pkg/e"
	"gindemo/pkg/logging"
	"gindemo/pkg/setting"
	"gindemo/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)

//GetArticle
// @Summary 查询单篇文章
// @Produce  json
// @Param id query string true "id"
// @success 200 {object} gin.H
// @Router /api/v1/tag/:id [get]
func GetArticle(c *gin.Context) {

	appG := app.Gin{c}

	var data interface{}
	articleId := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	code := e.INVALID_PARAMS
	valid.Required(articleId, "id").Message("id必须传")
	valid.Min(articleId, 1, "id").Message("Id 必须大于0")
	if !valid.HasErrors() {
		if models.ExistTagByID(articleId) {
			data = models.GetActile(articleId)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}

	//c.JSON(http.StatusOK, gin.H{
	//	"code":    code,
	//	"msg":     e.GetMsg(code),
	//	"article": data,
	//})
	appG.Response(code, e.SUCCESS, data)

}

//GetArticles
// @Summary 查询多篇文章
// @Produce  json
// @Param state query string true "state"
// @Param tag_id query string true "tag_id"
// @success 200 {object} gin.H
// @Router /api/v1/tags [get]
func GetArticles(c *gin.Context) {
	data := make(map[string]interface{})
	maps := make(map[string]interface{})
	valid := validation.Validation{}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}
	var tagId = -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		maps["tag_id"] = tagId
		valid.Min(tagId, 0, "tag_id").Message("标签ID必须大于0")
	}
	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		code = e.SUCCESS
		data["lists"] = models.GetActiles(util.GetPage(c), setting.AppSetting.PageSize, maps)
		data["total"] = models.GetArticlesTotal(maps)
	} else {
		for _, err := range valid.Errors {
			logging.Info("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})

}

//AddArticle
// @Summary 新增文章
// @Produce  json
// @Param state query string true "state"
// @Param tag_id query string true "tag_id"
// @Param title query string true "title"
// @Param desc query string true "desc"
// @Param content query string true "content"
// @Param created_by query string true "created_by"
// @success 200 {object} gin.H
// @Router /api/v1/tags [post]
func AddArticle(c *gin.Context) {

	tagId := com.StrTo(c.Query("tag_id")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	createdBy := c.Query("created_by")
	state := com.StrTo(c.Query("state")).MustInt()
	coverImageUrl := c.Query("cover_image_url")
	valid := validation.Validation{}

	valid.Min(tagId, 1, "tag_id").Message("标签Id必须大于0")
	valid.Required(title, "title").Message("标题必须填写")
	valid.Required(desc, "desc").Message("描述必须填写")
	valid.Required(content, "content").Message("内容不能为空")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.Required(createdBy, "cover_image_url").Message("文章封面不能为空")
	valid.Range(state, 0, 1, "state").Message("状态只允许是 0 或 1")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistTagByID(tagId) {
			condition := make(map[string]interface{})
			condition["tag_id"] = tagId
			condition["title"] = title
			condition["desc"] = desc
			condition["content"] = content
			condition["created_by"] = createdBy
			condition["state"] = state
			condition["cover_image_url"] = coverImageUrl
			models.AddArticle(condition)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info("err key:%s,err value:%s", err.Key, err.Value)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]interface{}),
	})
}

//EditArticle
// @Summary 修改文章
// @Produce  json
// @Param id query string true "id"
// @Param tag_id query string true "tag_id"
// @Param title query string true "title"
// @Param desc query string true "desc"
// @Param content query string true "content"
// @Param modified_by query string true "modified_by"
// @success 200 {object} gin.H
// @Router /api/v1/tags/:id [put]
func EditArticle(c *gin.Context) {
	valid := validation.Validation{}

	id := com.StrTo(c.Param("id")).MustInt()
	tagId := com.StrTo(c.Query("tag_id")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	modifiedBy := c.Query("modified_by")
	coverImageUrl := c.Query("cover_image_url")

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	valid.Min(id, 1, "id").Message("ID必须大于0")
	valid.MaxSize(title, 100, "title").Message("标题最长为100字符")
	valid.MaxSize(desc, 255, "desc").Message("简述最长为255字符")
	valid.MaxSize(content, 65535, "content").Message("内容最长为65535字符")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.Required(coverImageUrl, "cover_image_url").Message("文章封面不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			if models.ExistTagByID(tagId) {
				data := make(map[string]interface{})
				if tagId > 0 {
					data["tag_id"] = tagId
				}
				if title != "" {
					data["title"] = title
				}
				if desc != "" {
					data["desc"] = desc
				}
				if content != "" {
					data["content"] = content
				}

				data["modified_by"] = modifiedBy

				models.EditArticle(id, data)
				code = e.SUCCESS
			} else {
				code = e.ERROR_NOT_EXIST_TAG
			}
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

//DeleteArticle
// @Summary 删除文章
// @Produce  json
// @Param id query string true "id"
// @success 200 {object} gin.H
// @Router /api/v1/tags/:id [delete]
func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Query("id")).MustInt()
	valid := validation.Validation{}
	valid.Required(id, "id").Message("ID 必须填写")
	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			models.DeleteArticle(id)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}
