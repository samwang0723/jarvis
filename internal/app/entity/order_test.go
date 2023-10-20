package entity

import (
	"flag"
	"os"
	"testing"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	leak := flag.Bool("leak", false, "use leak detector")

	if *leak {
		goleak.VerifyTestMain(m)

		return
	}

	os.Exit(m.Run())
}

func TestCalculateProfitLoss(t *testing.T) {
	t.Parallel()

	tests := []struct {
		order        Order
		currentPrice float32
		expectedPL   float32
		expectedPLP  float32
	}{
		{
			order: Order{
				BuyQuantity:      2,
				BuyPrice:         84.9,
				BuyExchangeDate:  "20231011",
				SellQuantity:     2,
				SellPrice:        88.6,
				SellExchangeDate: "20231012",
			},
			expectedPL:  6745,
			expectedPLP: 3.97,
		},
		{
			order: Order{
				BuyQuantity:      2,
				BuyPrice:         84.9,
				BuyExchangeDate:  "20231011",
				SellQuantity:     2,
				SellPrice:        88.6,
				SellExchangeDate: "20231011",
			},
			expectedPL:  7011,
			expectedPLP: 4.13,
		},
	}

	for _, test := range tests {
		order := test.order
		order.CalculateProfitLoss()
		if order.ProfitLoss != test.expectedPL {
			t.Errorf("expected ProfitLoss %v but got %v", test.expectedPL, order.ProfitLoss)
		}
		if order.ProfitLossPercent != test.expectedPLP {
			t.Errorf("expected ProfitLossPercent %v but got %v", test.expectedPLP, order.ProfitLossPercent)
		}
	}
}

func TestCalculateUnrealizedProfitLoss(t *testing.T) {
	t.Parallel()

	tests := []struct {
		order        Order
		currentPrice float32
		expectedPL   float32
		expectedPLP  float32
	}{
		{
			order: Order{
				BuyQuantity:  2,
				BuyPrice:     84.9,
				SellQuantity: 0,
				SellPrice:    0,
			},
			currentPrice: 88.6,
			expectedPL:   6745,
			expectedPLP:  3.97,
		},
		{
			order: Order{
				BuyQuantity:  2,
				BuyPrice:     84.9,
				SellQuantity: 1,
				SellPrice:    88.7,
			},
			currentPrice: 86.6,
			expectedPL:   4851,
			expectedPLP:  2.86,
		},
	}

	for _, test := range tests {
		order := test.order
		order.CalculateUnrealizedProfitLoss(test.currentPrice)
		if order.ProfitLoss != test.expectedPL {
			t.Errorf("expected ProfitLoss %v but got %v", test.expectedPL, order.ProfitLoss)
		}
		if order.ProfitLossPercent != test.expectedPLP {
			t.Errorf("expected ProfitLossPercent %v but got %v", test.expectedPLP, order.ProfitLossPercent)
		}
	}
}
