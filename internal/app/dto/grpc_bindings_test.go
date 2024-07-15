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
