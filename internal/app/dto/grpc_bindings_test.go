package dto_test

import (
	"flag"
	"os"
	"testing"

	"github.com/samwang0723/jarvis/internal/app/dto"
	"github.com/samwang0723/jarvis/internal/app/pb"
	"github.com/stretchr/testify/assert"
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

func ptrToString(s string) *string {
	return &s
}

func TestListDailyCloseRequestFromPB(t *testing.T) {
	t.Parallel()

	type args struct {
		in *pb.ListDailyCloseRequest
	}

	var endDate string
	endDate = "20240704"

	tests := []struct {
		name string
		args args
		want *dto.ListDailyCloseRequest
	}{
		{
			name: "NilInput",
			args: args{
				in: nil,
			},
			want: nil,
		},
		{
			name: "ValidInput",
			args: args{
				in: &pb.ListDailyCloseRequest{
					Offset: 10,
					Limit:  20,
					SearchParams: &pb.ListDailyCloseSearchParams{
						StockID: "2330",
						Start:   "20240701",
						End:     endDate,
					},
				},
			},
			want: &dto.ListDailyCloseRequest{
				Offset: 10,
				Limit:  20,
				SearchParams: &dto.ListDailyCloseSearchParams{
					StockID: "2330",
					Start:   "20240701",
					End:     &endDate,
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := dto.ListDailyCloseRequestFromPB(tt.args.in)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestListStockRequestFromPB(t *testing.T) {
	tests := []struct {
		name string
		in   *pb.ListStockRequest
		want *dto.ListStockRequest
	}{
		{
			name: "nil input",
			in:   nil,
			want: nil,
		},
		{
			name: "valid input",
			in: &pb.ListStockRequest{
				Offset: 10,
				Limit:  20,
				SearchParams: &pb.ListStockSearchParams{
					StockIDs: []string{"AAPL", "GOOGL"},
					Country:  "USA",
					Name:     "Tech",
					Category: "Technology",
				},
			},
			want: &dto.ListStockRequest{
				Offset: 10,
				Limit:  20,
				SearchParams: &dto.ListStockSearchParams{
					StockIDs: &[]string{"AAPL", "GOOGL"},
					Country:  "USA",
					Name:     ptrToString("Tech"),
					Category: ptrToString("Technology"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dto.ListStockRequestFromPB(tt.in)
			assert.Equal(t, tt.want, got)
		})
	}
}
