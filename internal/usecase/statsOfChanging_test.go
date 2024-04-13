package usecase

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/egor-denisov/biggest-change/internal/entity"
	mock "github.com/egor-denisov/biggest-change/internal/usecase/mocks"

	"github.com/go-playground/assert"
	"github.com/golang/mock/gomock"
)

var errSomethingWentWrong = errors.New("something went wrong")

func Test_GetAddressWithBiggestChange(t *testing.T) {
	for _, test := range tests_GetAddressWithBiggestChange {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock.NewMockStatsOfChangingWebAPI(c)
			test.mockBehavior(service, test.countOfBlocks)

			// Call function
			biggestChange, err := New(service).GetAddressWithBiggestChange(context.Background(), test.countOfBlocks)

			assert.Equal(t, test.expectedResult, biggestChange)
			assert.Equal(t, errors.Is(err, test.expectedError), true)
		})
	}
}

type mockBehavior func(m *mock.MockStatsOfChangingWebAPI, countOfBlocks uint)

var tests_GetAddressWithBiggestChange = []struct {
	name           string
	mockBehavior   mockBehavior
	countOfBlocks  uint
	expectedResult *entity.BiggestChange
	expectedError  error
}{
	{
		name: "Success - Single Block",
		mockBehavior: func(m *mock.MockStatsOfChangingWebAPI, countOfBlocks uint) {
			m.EXPECT().GetCurrentBlockNumber(context.Background()).Return(big.NewInt(200), nil)
			m.EXPECT().GetTransactionsByBlockNumber(gomock.Any(), big.NewInt(200)).Return([]*entity.Transaction{
				{From: "0x1", To: "0x2", Value: big.NewInt(500), Gas: big.NewInt(50), GasPrice: big.NewInt(2)},
				{From: "0x3", To: "0x4", Value: big.NewInt(100), Gas: big.NewInt(20), GasPrice: big.NewInt(2)},
			}, nil)
		},
		countOfBlocks: 1,
		expectedResult: &entity.BiggestChange{
			Address:       "0x1",
			Amount:        "0x258",
			IsRecieved:    false,
			LastBlock:     "0xc8",
			CountOfBlocks: 1,
		},
		expectedError: nil,
	},
	{
		name: "Success - Many Blocks",
		mockBehavior: func(m *mock.MockStatsOfChangingWebAPI, countOfBlocks uint) {
			m.EXPECT().GetCurrentBlockNumber(context.Background()).Return(big.NewInt(200), nil)
			for i := uint(0); i < countOfBlocks; i++ {
				blockNum := big.NewInt(200 - int64(i))
				m.EXPECT().GetTransactionsByBlockNumber(gomock.Any(), blockNum).Return([]*entity.Transaction{
					{From: "0x1", To: "0x3", Value: big.NewInt(300 * int64(i+1)), Gas: big.NewInt(30), GasPrice: big.NewInt(3)},
					{From: "0x4", To: "0x3", Value: big.NewInt(100 * int64(i+1)), Gas: big.NewInt(10), GasPrice: big.NewInt(3)},
				}, nil)
			}
		},
		countOfBlocks: 3,
		expectedResult: &entity.BiggestChange{
			Address:       "0x3",
			Amount:        "0x960",
			IsRecieved:    true,
			LastBlock:     "0xc8",
			CountOfBlocks: 3,
		},
		expectedError: nil,
	},
	{
		name: "Zero Blocks",
		mockBehavior: func(m *mock.MockStatsOfChangingWebAPI, countOfBlocks uint) {
			m.EXPECT().GetCurrentBlockNumber(context.Background()).Return(big.NewInt(200), nil)
		},
		countOfBlocks: 0,
		expectedResult: &entity.BiggestChange{
			Amount:        "0x0",
			IsRecieved:    false,
			LastBlock:     "0xc8",
			CountOfBlocks: 0,
		},
		expectedError: nil,
	},
	{
		name: "Error - Getting block number",
		mockBehavior: func(m *mock.MockStatsOfChangingWebAPI, countOfBlocks uint) {
			m.EXPECT().GetCurrentBlockNumber(context.Background()).Return(nil, errSomethingWentWrong)
		},
		countOfBlocks:  1,
		expectedResult: nil,
		expectedError:  errSomethingWentWrong,
	},
	{
		name: "Error - Getting change map",
		mockBehavior: func(m *mock.MockStatsOfChangingWebAPI, countOfBlocks uint) {
			m.EXPECT().GetCurrentBlockNumber(context.Background()).Return(big.NewInt(200), nil)
			m.EXPECT().GetTransactionsByBlockNumber(gomock.Any(), big.NewInt(200)).Return(nil, errSomethingWentWrong)
		},
		countOfBlocks:  1,
		expectedResult: nil,
		expectedError:  errSomethingWentWrong,
	},
}
