package controller

import (
	sarKaf "L0/internal"
	"L0/models/compositeView/storage"
	"L0/models/order"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Controller struct {
	cache           *expirable.LRU[string, order.Order]
	repo            storage.Repository
	kafkaController *sarKaf.SarKafCont
}

func NewController(repo storage.Repository, cache *expirable.LRU[string, order.Order], kafkaController *sarKaf.SarKafCont) *Controller {
	return &Controller{cache: cache, repo: repo, kafkaController: kafkaController}
}

type Message struct {
	Message string `json:"message" example:"message"`
}

type HTTPError struct {
	Error   bool   `json:"error" example:"true"`
	Message string `json:"message" example:"forbidden"`
}

func NewHTTPError(ctx *gin.Context, b bool, err error) {
	ctx.JSON(http.StatusInternalServerError, HTTPError{
		Error:   b,
		Message: err.Error(),
	})
}

func (cont *Controller) GetOrderById(ctx *gin.Context) {
	var result order.Order

	order_uid := ctx.Param("order_uid")

	//Check cache
	result, hasInCache := cont.cache.Get(order_uid)
	if !hasInCache {
		idBytes, err := json.Marshal(order_uid)
		if err != nil {
			NewHTTPError(ctx, true, err)
		}
		msg := sarKaf.ControllerRequest(cont.kafkaController.Producer, cont.kafkaController.Consumer, idBytes, "getOrderById", "resGetOrderById")

		err = json.Unmarshal([]byte(msg), &result)
		if err != nil {
			NewHTTPError(ctx, true, err)
		}
		fmt.Printf("Записи %s нет в кэше, загружена из БД", order_uid)
		cont.cache.Add(result.Order_uid, result)

	} else {
		fmt.Printf("Запись %s есть в кэше", order_uid)
	}
	ctx.JSON(http.StatusOK, result)
}

func (cont *Controller) GetAllOrders(ctx *gin.Context) {
	var (
		result []order.Order
		err    error
	)

	result, err = cont.repo.FindAll(ctx, &pgxpool.Pool{})
	if err != nil {
		NewHTTPError(ctx, true, err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (cont *Controller) PostOrders(ctx *gin.Context) {
	var value []byte // Not needed (Заглушка)
	msg := sarKaf.ControllerRequest(cont.kafkaController.Producer, cont.kafkaController.Consumer, value, "postOrder", "resPostOrder")

	if msg == "Generated an order" {
		ctx.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": msg,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"error":   true,
			"message": msg,
		})
	}
}

func (cont *Controller) MainPage(ctx *gin.Context) {
	ctx.File(filepath.Join("./view", "index.html"))
}

func (cont *Controller) RedirectOnMainPage(ctx *gin.Context) {
	ctx.Redirect(http.StatusMovedPermanently, "/order")
}
