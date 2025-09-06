package router

import (
	"L0/controller"
	sarKaf "L0/internal"
	"L0/models/compositeView/storage"
	"L0/models/order"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/golang-lru/v2/expirable"
)

func StartRouter(ctx context.Context, repoS storage.Repository, cache *expirable.LRU[string, order.Order], kafkaCont *sarKaf.SarKafCont) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	//gin.SetMode(gin.DebugMode)

	router := gin.New()
	router.StaticFS("/static", http.Dir("./view"))
	router.Use(gin.Logger(), gin.Recovery())

	cont := controller.NewController(repoS, cache, kafkaCont)

	router.GET("/order", cont.MainPage)
	router.GET("/", cont.RedirectOnMainPage)

	router.GET("/order/:order_uid", cont.GetOrderById)
	router.POST("/order/fake", cont.PostOrders)

	return router
}
