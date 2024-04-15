package v1

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/egor-denisov/biggest-change/internal/usecase"
	sl "github.com/egor-denisov/biggest-change/pkg/logger"
	"github.com/gin-gonic/gin"
)

var _defaultCountOfBlocks uint = 100

type statsOfChangingRoutes struct {
	sc usecase.StatsOfChanging
	l  *slog.Logger
}

func newStatsOfChanging(handler *gin.RouterGroup, l *slog.Logger, sc usecase.StatsOfChanging) {
	r := &statsOfChangingRoutes{sc, l}

	h := handler.Group("/")
	{
		h.GET("/get_biggest_change", r.getBiggestChange)
	}
}

type getBiggestChangeRequest struct {
	CountOfBlocks uint `form:"count_of_blocks"`
}

// @Summary     Получение адреса, который максимально
// @Description Получение адреса, который максимально изменился за count_of_blocks блоков
// @Description По умолчанию count_of_blocks = 100
// @Tags  	    StatsOfChanging
// @Param count_of_blocks query integer false "Количество последних блоков"
// @Success     200 {object} entity.BiggestChange "Адрес найден"
// @Failure     400 "Ошибка в запросе"
// @Failure     500 "Не удалось выполнить запрос"
// @Failure     500 "Таймаут запроса"
// @Router      /get_biggest_change [get] .
func (r *statsOfChangingRoutes) getBiggestChange(c *gin.Context) {
	var input getBiggestChangeRequest

	if err := c.ShouldBind(&input); err != nil {
		r.l.Error("http - v1 - getBiggestChange", sl.Err(err))
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	if input.CountOfBlocks == 0 {
		input.CountOfBlocks = _defaultCountOfBlocks
	}

	res, err := r.sc.GetAddressWithBiggestChange(c.Request.Context(), input.CountOfBlocks)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			c.AbortWithStatus(http.StatusGatewayTimeout)

			return
		}

		r.l.Error("http - v1 - getBiggestChange", sl.Err(err))
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	c.JSON(http.StatusOK, res)
}
