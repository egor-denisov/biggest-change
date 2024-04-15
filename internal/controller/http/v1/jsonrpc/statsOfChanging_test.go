package jsonrpc

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"

	"github.com/egor-denisov/biggest-change/internal/entity"
	mock "github.com/egor-denisov/biggest-change/internal/usecase/mocks"
	"github.com/egor-denisov/biggest-change/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	"github.com/golang/mock/gomock"
)

func Test_GetBiggestChange(t *testing.T) {
	for _, test := range tests_GetBiggestChange {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			webapi := mock.NewMockStatsOfChanging(c)
			test.mockBehavior(webapi)

			rpcServer := rpc.NewServer()
			rpcServer.RegisterCodec(json.NewCodec(), "application/json")

			s := NewStatsOfChangingService(logger.SetupLogger("debug"), webapi)
			rpcServer.RegisterService(s, "JsonRpc")

			// Init Endpoint
			r := gin.New()
			r.POST("/", gin.WrapH(rpcServer))

			// Create Request
			req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(test.requestBody))
			assert.Equal(t, err, nil)

			req.Header.Set("Content-Type", "application/json")

			// Make Request
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, strings.TrimSpace(w.Body.String()), strings.TrimSpace(test.expectedResponseBody))
		})
	}
}

func getBodyRequestByCountOfBlock(countOfBlocks int) string {
	return fmt.Sprintf(
		`
		{
			"id": "1",
			"jsonrpc": "2.0",
			"method": "JsonRpc.GetBiggestChange",
			"params": [{
				"countOfBlocks": %d
			}]
		}
		`, countOfBlocks,
	)
}

type mockBehavior func(m *mock.MockStatsOfChanging)

var tests_GetBiggestChange = []struct {
	name                 string
	requestBody          string
	mockBehavior         mockBehavior
	expectedResponseBody string
	expectedError        error
}{
	{
		name:        "Success",
		requestBody: getBodyRequestByCountOfBlock(50),
		mockBehavior: func(m *mock.MockStatsOfChanging) {
			result := &entity.BiggestChange{
				Address:       "0x1",
				Amount:        "0x100",
				LastBlock:     "0x123",
				CountOfBlocks: 50,
				IsRecieved:    true,
			}
			m.EXPECT().GetAddressWithBiggestChange(gomock.Any(), uint(50)).Return(result, nil)
		},
		expectedResponseBody: `{"result":{"address":"0x1","amount":"0x100","lastBlock":"0x123",` +
			`"countOfBlocks":50,"isRecieved":true},"error":null,"id":"1"}`,
	},
	{
		name:        "Default Count Of Blocks",
		requestBody: getBodyRequestByCountOfBlock(0),
		mockBehavior: func(m *mock.MockStatsOfChanging) {
			result := &entity.BiggestChange{
				Address:       "0x1",
				Amount:        "0x100",
				LastBlock:     "0x123",
				CountOfBlocks: int64(100),
				IsRecieved:    true,
			}
			m.EXPECT().GetAddressWithBiggestChange(gomock.Any(), uint(0)).Return(result, nil)
		},
		expectedResponseBody: `{"result":{"address":"0x1","amount":"0x100","lastBlock":"0x123",` +
			`"countOfBlocks":100,"isRecieved":true},"error":null,"id":"1"}`,
	},
	{
		name:        "Timeout Error Handling",
		requestBody: getBodyRequestByCountOfBlock(10),
		mockBehavior: func(m *mock.MockStatsOfChanging) {
			m.EXPECT().GetAddressWithBiggestChange(gomock.Any(), uint(10)).Return(nil, context.DeadlineExceeded)
		},
		expectedResponseBody: `{"result":null,"error":"process timeout","id":"1"}`,
	},
	{
		name:        "Else Error Handling",
		requestBody: getBodyRequestByCountOfBlock(10),
		mockBehavior: func(m *mock.MockStatsOfChanging) {
			m.EXPECT().GetAddressWithBiggestChange(gomock.Any(), uint(10)).Return(nil, entity.ErrInternalServer)
		},
		expectedResponseBody: `{"result":null,"error":"internal server error","id":"1"}`,
	},
}
