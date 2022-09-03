package api

import (
	// internal packages
	http "net/http"
	// external packages
	gzip "github.com/gin-contrib/gzip"
	gin "github.com/gin-gonic/gin"
	// local packages
	middlewares "akshayGudhate/whatsapp-bridge/src/middlewares"
)

////////////////////
//     router     //
////////////////////

func GetAPIRouter() http.Handler {
	// router instance
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// cors middleware
	router.Use(middlewares.CORSMiddleware())
	// default compression middleware
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// routes
	apiRouter := router.Group("/api")
	{
		apiRouter.POST("/whatsapp/send", sendMessage)
		apiRouter.GET("/connect/qr", getConnectionQRCode)
	}

	// return router
	return router
}
