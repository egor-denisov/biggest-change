package v1

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/egor-denisov/biggest-change/internal/entity"
	mock "github.com/egor-denisov/biggest-change/internal/usecase/mocks"
	"github.com/egor-denisov/biggest-change/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	"github.com/golang/mock/gomock"
)

var errSomethingWentWrong = errors.New("something went wrong")

func Test_getBiggestChange(t *testing.T) {
	for _, test := range testsGetBiggestChange {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			usecase := mock.NewMockStatsOfChanging(c)
			test.mockBehavior(usecase)

			handler := statsOfChangingRoutes{
				sc: usecase,
				l:  logger.SetupLogger("debug"),
			}
			// Init Endpoint
			r := gin.New()
			r.GET("/", handler.getBiggestChange)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/"+test.query, nil)
			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

type mockBehavior func(m *mock.MockStatsOfChanging)

var testsGetBiggestChange = []struct {
	name                 string
	mockBehavior         mockBehavior
	query                string
	expectedStatusCode   int
	expectedResponseBody string
}{
	{
		name:  "valid request",
		query: `?count_of_blocks=50`,
		mockBehavior: func(m *mock.MockStatsOfChanging) {
			res := &entity.BiggestChange{
				Address:       "0x1",
				Amount:        "0x100",
				LastBlock:     "0x123",
				CountOfBlocks: 50,
				IsRecieved:    true,
			}
			m.EXPECT().GetAddressWithBiggestChange(gomock.Any(), uint(50)).Return(res, nil)
		},
		expectedStatusCode:   http.StatusOK,
		expectedResponseBody: `{"address":"0x1","amount":"0x100","lastBlock":"0x123","countOfBlocks":50,"isRecieved":true}`,
	},
	{
		name:  "default count of blocks",
		query: ``,
		mockBehavior: func(m *mock.MockStatsOfChanging) {
			res := &entity.BiggestChange{
				Address:       "0x1",
				Amount:        "0x100",
				LastBlock:     "0x123",
				CountOfBlocks: int64(100),
				IsRecieved:    true,
			}
			m.EXPECT().GetAddressWithBiggestChange(gomock.Any(), uint(0)).Return(res, nil)
		},
		expectedStatusCode:   http.StatusOK,
		expectedResponseBody: `{"address":"0x1","amount":"0x100","lastBlock":"0x123","countOfBlocks":100,"isRecieved":true}`,
	},
	{
		name:                 "Bad request",
		query:                `?count_of_blocks=hello`,
		mockBehavior:         func(_ *mock.MockStatsOfChanging) {},
		expectedStatusCode:   http.StatusBadRequest,
		expectedResponseBody: ``,
	},
	{
		name:  "Timeout",
		query: ``,
		mockBehavior: func(m *mock.MockStatsOfChanging) {
			m.EXPECT().GetAddressWithBiggestChange(gomock.Any(), uint(0)).
				Return(nil, context.DeadlineExceeded)
		},
		expectedStatusCode:   http.StatusGatewayTimeout,
		expectedResponseBody: ``,
	},
	{
		name:  "Something went wrong",
		query: ``,
		mockBehavior: func(m *mock.MockStatsOfChanging) {
			m.EXPECT().GetAddressWithBiggestChange(gomock.Any(), uint(0)).
				Return(nil, errSomethingWentWrong)
		},
		expectedStatusCode:   http.StatusInternalServerError,
		expectedResponseBody: ``,
	},
}
