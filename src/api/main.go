package main

import (
  "net/http"
  "github.com/gin-gonic/gin"
)

func main() {
  router := gin.Default()

  router.GET("/hello", func(c *gin.Context) {
    c.String(http.StatusOK, "Hello World!!Wow!!")
  })

  router.Run()
}
