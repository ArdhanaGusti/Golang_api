package routes

import (
	"strconv"
	"time"

	"github.com/ArdhanaGusti/Golang_api/config"
	"github.com/ArdhanaGusti/Golang_api/models"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

func Home(c *gin.Context) {
	items := []models.Article{}
	// config.DB.Find(&items)
	config.DB.Preload("User").Find(&items)
	c.JSON(200, gin.H{
		"Message": "Berhasil akses home",
		"Data":    items,
	})
}

func GetArticle(c *gin.Context) {
	slug := c.Param("slug")
	var item models.Article
	if err := config.DB.First(&item, "slug = ?", slug).Error; err != nil {
		c.JSON(404, gin.H{"status": "error"})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"data": item,
	})
}

func GetArticleTag(c *gin.Context) {
	tag := c.Param("tag")
	items := []models.Article{}
	config.DB.Where("tag LIKE ?", "%"+tag+"%").Find(&items)
	c.JSON(200, gin.H{
		"data": items,
	})
}

func PostArticle(c *gin.Context) {

	slug := slug.Make(c.PostForm("title"))
	var oldItem models.Article
	if err := config.DB.First(&oldItem, "slug = ?", slug).Error; err == nil {
		slug = slug + strconv.FormatInt(time.Now().Unix(), 10)
	}

	item := models.Article{
		Title:  c.PostForm("title"),
		Desc:   c.PostForm("desc"),
		Tag:    c.PostForm("tag"),
		Slug:   slug,
		UserID: uint(c.MustGet("jwt_user_id").(float64)),
	}

	config.DB.Create(&item)

	c.JSON(200, gin.H{
		"status": "berhasil",
		"data":   item,
	})
}

func UpdateArticle(c *gin.Context) {
	slug := c.Param("slug")
	var item models.Article
	if err := config.DB.First(&item, "slug = ?", slug).Error; err != nil {
		c.JSON(404, gin.H{"status": "error"})
		c.Abort()
		return
	}
	if uint(c.MustGet("jwt_user_id").(float64)) != item.UserID {
		c.JSON(403, gin.H{"status": "error", "msg": "Data is forbidden"})
		c.Abort()
		return
	}
	config.DB.Model(&item).Where("slug = ?", slug).Updates(models.Article{Title: c.PostForm("title"), Desc: c.PostForm("desc")})
	c.JSON(200, gin.H{
		"data":    item,
		"message": "Berhasil di update",
	})
}

func DeleteArticle(c *gin.Context) {
	slug := c.Param("slug")
	var item models.Article
	config.DB.Where("slug = ?", slug).Delete(&item)
	c.JSON(200, gin.H{
		"msg": "Berhasil di hapus",
	})
}
