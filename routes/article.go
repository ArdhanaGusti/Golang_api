package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ArdhanaGusti/Golang_api/config"
	"github.com/ArdhanaGusti/Golang_api/handler/failed"
	"github.com/ArdhanaGusti/Golang_api/handler/validation"
	"github.com/ArdhanaGusti/Golang_api/models"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

func Home(c *gin.Context) {
	items := []models.Article{}
	if err := config.DB.Preload("User").Find(&items).Error; err != nil {
		c.JSON(500, failed.FailedResponse{
			StatusCode: 500,
			Message:    err.Error(),
		})
		c.Abort()
		return
	}

	itemsJson, err := json.Marshal(items)
	if err != nil {
		c.JSON(500, failed.FailedResponse{
			StatusCode: 500,
			Message:    "Failed to set json because: " + err.Error(),
		})
		c.Abort()
	}

	if err := config.RDB.Set("articles", itemsJson, 0).Err(); err != nil {
		c.JSON(500, failed.FailedResponse{
			StatusCode: 500,
			Message:    "Failed to set redis because: " + err.Error(),
		})
		c.Abort()
		return
	}
	c.JSON(200, items)
}

func GetArticle(c *gin.Context) {
	slug := c.Param("slug")
	var item models.Article
	if err := config.DB.First(&item, "slug = ?", slug).Error; err != nil {
		c.JSON(500, failed.FailedResponse{
			StatusCode: 500,
			Message:    err.Error(),
		})
		c.Abort()
		return
	}

	c.JSON(200, item)
}

// func GetArticleTag(c *gin.Context) {
// 	tag := c.Param("tag")
// 	items := []models.Article{}
// 	if err := config.DB.Where("tag LIKE ?", "%"+tag+"%").Find(&items).Error; err != nil {
// 		c.JSON(500, failed.FailedResponse{
// 			StatusCode: 500,
// 			Message:    err.Error(),
// 		})
// 	}
// 	c.JSON(200, gin.H{
// 		"data": items,
// 	})
// }

func PostArticle(c *gin.Context) {
	var articlePayload validation.CreateArticlePayload

	if err := c.ShouldBind(&articlePayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	slug := slug.Make(articlePayload.Title)
	var oldItem models.Article
	if err := config.DB.First(&oldItem, "slug = ?", slug).Error; err == nil {
		slug = slug + strconv.FormatInt(time.Now().Unix(), 10)
	}

	item := models.Article{
		Title:  articlePayload.Title,
		Desc:   articlePayload.Desc,
		Tag:    articlePayload.Tag,
		Slug:   slug,
		UserID: uint(c.MustGet("jwt_user_id").(float64)),
	}

	if err := config.DB.Create(&item).Error; err != nil {
		c.JSON(500, failed.FailedResponse{
			StatusCode: 500,
			Message:    err.Error(),
		})
	}

	exist, _ := config.RDB.Exists("articles").Result()

	if exist > 0 {
		if err := config.RDB.Del("articles").Err(); err != nil {
			c.JSON(500, failed.FailedResponse{
				StatusCode: 500,
				Message:    "Failed to delete redis because: " + err.Error(),
			})
			c.Abort()
			return
		}
	}

	c.JSON(200, gin.H{
		"message": "Article " + articlePayload.Title + " Made Successfully",
	})
}

func UpdateArticle(c *gin.Context) {
	var articlePayload validation.CreateArticlePayload

	if err := c.ShouldBind(&articlePayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	slug := c.Param("slug")
	var item models.Article
	if err := config.DB.First(&item, "slug = ?", slug).Error; err != nil {
		c.JSON(404, gin.H{"status": "error"})
		c.Abort()
		return
	}

	if uint(c.MustGet("jwt_user_id").(float64)) != item.UserID {
		c.JSON(403, failed.FailedResponse{
			StatusCode: 403,
			Message:    "Data is forbidden",
		})
		c.Abort()
		return
	}

	updatedArticle := models.Article{
		Title: articlePayload.Title,
		Desc:  articlePayload.Desc,
		Tag:   articlePayload.Tag,
	}

	if err := config.DB.Model(&item).Where("slug = ?", slug).Updates(updatedArticle).Error; err != nil {
		c.JSON(500, failed.FailedResponse{
			StatusCode: 500,
			Message:    err.Error(),
		})
	}

	exist, _ := config.RDB.Exists("articles").Result()

	if exist > 0 {
		if err := config.RDB.Del("articles").Err(); err != nil {
			c.JSON(500, failed.FailedResponse{
				StatusCode: 500,
				Message:    "Failed to delete redis because: " + err.Error(),
			})
			c.Abort()
			return
		}
	}

	c.JSON(200, gin.H{
		"message": "Article " + updatedArticle.Title + " Updated Successfully",
	})
}

func DeleteArticle(c *gin.Context) {
	slug := c.Param("slug")
	var item models.Article
	if err := config.DB.Where("slug = ?", slug).First(&item).Error; err != nil {
		c.JSON(500, failed.FailedResponse{
			StatusCode: 500,
			Message:    err.Error(),
		})
	}

	var title = item.Title

	if err := config.DB.Where("slug = ?", slug).Delete(&item).Error; err != nil {
		c.JSON(500, failed.FailedResponse{
			StatusCode: 500,
			Message:    err.Error(),
		})
	}

	exist, _ := config.RDB.Exists("articles").Result()

	if exist > 0 {
		if err := config.RDB.Del("articles").Err(); err != nil {
			c.JSON(500, failed.FailedResponse{
				StatusCode: 500,
				Message:    "Failed to delete redis because: " + err.Error(),
			})
			c.Abort()
			return
		}
	}

	c.JSON(200, gin.H{
		"message": "Article " + title + " Deleted Successfully",
	})
}
