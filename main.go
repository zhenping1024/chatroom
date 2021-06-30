package main

import "github.com/gin-gonic/gin"

func main(){
	Router:=gin.Default()
	go h.run()
	Router.GET("/ws",myws)
	Router.Run(":9090")
}
