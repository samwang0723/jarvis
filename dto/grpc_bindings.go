package dto

import (
	"samwang0723/jarvis/entity"
	pb "samwang0723/jarvis/pb"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func ListDailyCloseRequestFromPB(in *pb.ListDailyCloseRequest) *ListDailyCloseRequest {
	if in == nil {
		return nil
	}
	out := &ListDailyCloseRequest{
		Offset:       in.Offset,
		Limit:        in.Limit,
		SearchParams: ListDailyCloseSearchParamsFromPB(in.SearchParams),
	}

	return out
}

func ListDailyCloseSearchParamsFromPB(in *pb.ListDailyCloseSearchParams) *ListDailyCloseSearchParams {
	if in == nil {
		return nil
	}

	out := &ListDailyCloseSearchParams{
		Start: in.Start,
	}

	stockIDs := in.StockIDs
	if stockIDs != nil {
		out.StockIDs = &stockIDs
	}
	end := in.End
	if len(end) > 0 {
		out.End = &end
	}
	return out
}

func ListDailyCloseResponseToPB(in *ListDailyCloseResponse) *pb.ListDailyCloseResponse {
	if in == nil {
		return nil
	}

	var entries []*pb.DailyClose
	for _, obj := range in.Entries {
		entries = append(entries, DailyCloseToPB(obj))
	}

	return &pb.ListDailyCloseResponse{
		Offset:     in.Offset,
		Limit:      in.Limit,
		TotalCount: in.TotalCount,
		Entries:    entries,
	}
}

func DailyCloseToPB(in *entity.DailyClose) *pb.DailyClose {
	if in == nil {
		return nil
	}

	pbDailyCloseID := in.ID
	pbStockID := in.StockID
	pbDate := in.Date
	pbTradeShares := in.TradedShares
	pbTransactions := in.Transactions
	pbTurnover := in.Turnover
	pbOpen := in.Open
	pbClose := in.Close
	pbHigh := in.High
	pbLow := in.Low
	pbDiff := in.PriceDiff

	var pbCreatedAt *timestamp.Timestamp
	if in.CreatedAt != nil {
		pbCreatedAt, _ = ptypes.TimestampProto(*in.CreatedAt)
	}

	var pbUpdatedAt *timestamp.Timestamp
	if in.UpdatedAt != nil {
		pbUpdatedAt, _ = ptypes.TimestampProto(*in.UpdatedAt)
	}

	var pbDeletedAt *timestamp.Timestamp
	if in.DeletedAt != nil {
		pbDeletedAt, _ = ptypes.TimestampProto(*in.DeletedAt)
	}

	return &pb.DailyClose{
		DailyCloseID: pbDailyCloseID.Uint64(),
		StockID:      pbStockID,
		Date:         pbDate,
		TradeShares:  pbTradeShares,
		Transactions: pbTransactions,
		Turnover:     pbTurnover,
		Open:         pbOpen,
		Close:        pbClose,
		High:         pbHigh,
		Low:          pbLow,
		Diff:         pbDiff,
		CreatedAt:    pbCreatedAt,
		UpdatedAt:    pbUpdatedAt,
		DeletedAt:    pbDeletedAt,
	}
}
