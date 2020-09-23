package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"fullRssAPI/logs"
	"fullRssAPI/model/db"
	"fullRssAPI/server"
	uuid "github.com/satori/go.uuid"
)

func Register(data db.SQLDataStorage) func(c *gin.Context) {
	return func(c *gin.Context) {
		var feedInfo db.FeedInfo
		c.BindJSON(&feedInfo)

		title, err := server.GetTitle(&feedInfo)
		if err != nil {
			logs.Error("GetTitle: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "BadRequest"})
			return
		}
		feedInfo.Title = title

		err = data.PutFeedInfo(&feedInfo)
		if err != nil {
			logs.Error("PutFeedInfo: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "BadRequest"})
			return
		}

		err = server.ServeFeed(&feedInfo)
		if err != nil {
			logs.Error("ServeFeed: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "BadRequest"})
			err = data.DeleteFeedInfo(feedInfo.ID)
			if err != nil {
				logs.Error("DeleteFeedInfo: %v", err)
				return
			}
			return
		}
		
		c.JSON(http.StatusOK, feedInfo)
	}
}

func GetFeedInfo(data db.SQLDataStorage) func(c *gin.Context) {
	return func(c *gin.Context) {
		feedInfoList, err := data.GetFeedInfo()
		if err != nil {
			logs.Error("GetFeedInfo: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "BadRequest"})
			return
		}
		c.JSON(http.StatusOK, feedInfoList)
	}
}

func Delete(data db.SQLDataStorage) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Query("id")
		err := data.DeleteFeedInfo(uuid.FromStringOrNil(id))
		if err != nil {
			logs.Error("DeleteFeedInfo: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "BadRequest"})
			return
		}
		if err := os.Remove(fmt.Sprintf("./rss/%s.xml", id)); err != nil {
			logs.Error("DeleteFile: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "BadRequest"})
			return
		}
	}
}

func GetRSS(c *gin.Context) {
	id := c.Param("id")
	c.Header("Content-Type", "application/xml")
	bytes, err := ioutil.ReadFile(fmt.Sprintf("./rss/%s.xml", id))
	if err != nil {
		logs.Error("GetRSS: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "BadRequest"})
		return
	}
	c.Writer.Write(bytes)
}

func Refresh(data db.SQLDataStorage) func(c *gin.Context) {
	return func(c *gin.Context) {
		feedInfoList, err := data.GetFeedInfo()
		if err != nil {
			logs.Error("GetFeedInfo: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "BadRequest"})
			return
		}
		for _, feedInfo := range feedInfoList {
			server.ServeFeed(feedInfo)
		}
	}
}