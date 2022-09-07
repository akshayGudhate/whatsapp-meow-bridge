package api

import (
	// internal packages
	http "net/http"
	// external packages
	gin "github.com/gin-gonic/gin"
	// local packages
	bridge "akshayGudhate/whatsapp-bridge/src/bridge"
)

// customer update form
type ST_Form_SendMessage struct {
	FromPhone   string `form:"fromPhone,omitempty"`
	ToPhone     string `form:"toPhone,omitempty"`
	MessageText string `form:"messageText,omitempty"`
}

//////////////////
//    routes    //
//////////////////

// send whatsapp message
func sendMessage(c *gin.Context) {
	// extract details
	var formData ST_Form_SendMessage
	err := c.ShouldBind(&formData)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// send whatsapp message
	responseString := *bridge.SendWhatsappMessage(
		&formData.FromPhone, &formData.ToPhone, &formData.MessageText,
	)
	if responseString != "" {
		// on error
		c.JSON(http.StatusOK, gin.H{
			"info": responseString,
			"data": nil,
		})
		return
	}

	// message sent
	c.JSON(http.StatusOK, gin.H{
		"info": "Message sent!",
		"data": nil,
	})
}
