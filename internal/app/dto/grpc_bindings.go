// Copyright 2021 Wei (Sam) Wang <sam.wang.0723@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package dto

import (
	"encoding/json"

	"github.com/samwang0723/jarvis/internal/app/businessmodel"
	"github.com/samwang0723/jarvis/internal/app/entity"
	pb "github.com/samwang0723/jarvis/internal/app/pb"

	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func ListStockRequestFromPB(in *pb.ListStockRequest) *ListStockRequest {
	if in == nil {
		return nil
	}
	out := &ListStockRequest{
		Offset:       in.Offset,
		Limit:        in.Limit,
		SearchParams: ListStockSearchParamsFromPB(in.SearchParams),
	}

	return out
}

func ListStockSearchParamsFromPB(in *pb.ListStockSearchParams) *ListStockSearchParams {
	if in == nil {
		return nil
	}

	out := &ListStockSearchParams{}

	stockIDs := in.StockIDs
	if stockIDs != nil {
		out.StockIDs = &stockIDs
	}
	country := in.Country
	if len(country) > 0 {
		out.Country = country
	}
	name := in.Name
	if len(name) > 0 {
		out.Name = &name
	}
	category := in.Category
	if len(category) > 0 {
		out.Category = &category
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

	var averages []*pb.Average
	for _, obj := range in.Averages {
		averages = append(averages, AverageToPB(obj))
	}

	return &pb.ListDailyCloseResponse{
		Offset:     in.Offset,
		Limit:      in.Limit,
		TotalCount: in.TotalCount,
		Entries:    entries,
		Averages:   averages,
	}
}

func MapToProtobufStructFloat32(m map[int]float32) *structpb.Struct {
	b, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	s := &structpb.Struct{}
	err = protojson.Unmarshal(b, s)
	if err != nil {
		return nil
	}
	return s
}

func MapToProtobufStructUint64(m map[int]uint64) *structpb.Struct {
	b, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	s := &structpb.Struct{}
	err = protojson.Unmarshal(b, s)
	if err != nil {
		return nil
	}
	return s
}

func AverageToPB(in *businessmodel.Average) *pb.Average {
	if in == nil {
		return nil
	}

	return &pb.Average{
		StockID: in.StockID,
		Ma:      MapToProtobufStructFloat32(in.MA),
		Mv:      MapToProtobufStructUint64(in.MV),
	}
}

func DailyCloseToPB(in *entity.DailyClose) *pb.DailyClose {
	if in == nil {
		return nil
	}

	pbID := in.ID
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
		pbCreatedAt = timestamppb.New(*in.CreatedAt)
	}

	var pbUpdatedAt *timestamp.Timestamp
	if in.UpdatedAt != nil {
		pbUpdatedAt = timestamppb.New(*in.UpdatedAt)
	}

	var pbDeletedAt *timestamp.Timestamp
	if in.DeletedAt != nil {
		pbDeletedAt = timestamppb.New(*in.DeletedAt)
	}

	return &pb.DailyClose{
		Id:           pbID.Uint64(),
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

func ListStockResponseToPB(in *ListStockResponse) *pb.ListStockResponse {
	if in == nil {
		return nil
	}

	var entries []*pb.Stock
	for _, obj := range in.Entries {
		entries = append(entries, StockToPB(obj))
	}

	return &pb.ListStockResponse{
		Offset:     in.Offset,
		Limit:      in.Limit,
		TotalCount: in.TotalCount,
		Entries:    entries,
	}
}

func StockToPB(in *entity.Stock) *pb.Stock {
	if in == nil {
		return nil
	}

	pbID := in.ID
	pbStockID := in.StockID
	pbName := in.Name
	pbCategory := in.Category
	pbCountry := in.Country

	var pbCreatedAt *timestamp.Timestamp
	if in.CreatedAt != nil {
		pbCreatedAt = timestamppb.New(*in.CreatedAt)
	}

	var pbUpdatedAt *timestamp.Timestamp
	if in.UpdatedAt != nil {
		pbUpdatedAt = timestamppb.New(*in.UpdatedAt)
	}

	var pbDeletedAt *timestamp.Timestamp
	if in.DeletedAt != nil {
		pbDeletedAt = timestamppb.New(*in.DeletedAt)
	}

	return &pb.Stock{
		Id:        pbID.Uint64(),
		StockID:   pbStockID,
		Name:      pbName,
		Category:  pbCategory,
		Country:   pbCountry,
		CreatedAt: pbCreatedAt,
		UpdatedAt: pbUpdatedAt,
		DeletedAt: pbDeletedAt,
	}
}

func ListCategoriesResponseToPB(in *ListCategoriesResponse) *pb.ListCategoriesResponse {
	if in == nil {
		return nil
	}

	var entries []string
	entries = append(entries, in.Entries...)

	return &pb.ListCategoriesResponse{
		Entries: entries,
	}
}

func GetStakeConcentrationRequestFromPB(in *pb.GetStakeConcentrationRequest) *GetStakeConcentrationRequest {
	if in == nil {
		return nil
	}
	return &GetStakeConcentrationRequest{
		StockID: in.StockID,
		Date:    in.Date,
	}
}

func GetStakeConcentrationResponseToPB(in *entity.StakeConcentration) *pb.GetStakeConcentrationResponse {
	if in == nil {
		return nil
	}

	pbID := in.ID
	pbStockID := in.StockID
	pbDate := in.Date
	pbSumBuyShares := in.SumBuyShares
	pbSumSellShares := in.SumSellShares
	pbAvgBuyPrice := in.AvgBuyPrice
	pbAvgSellPrice := in.AvgSellPrice
	pbConcentration1 := in.Concentration_1
	pbConcentration5 := in.Concentration_5
	pbConcentration10 := in.Concentration_10
	pbConcentration20 := in.Concentration_20
	pbConcentration60 := in.Concentration_60

	var pbCreatedAt *timestamp.Timestamp
	if in.CreatedAt != nil {
		pbCreatedAt = timestamppb.New(*in.CreatedAt)
	}

	var pbUpdatedAt *timestamp.Timestamp
	if in.UpdatedAt != nil {
		pbUpdatedAt = timestamppb.New(*in.UpdatedAt)
	}

	var pbDeletedAt *timestamp.Timestamp
	if in.DeletedAt != nil {
		pbDeletedAt = timestamppb.New(*in.DeletedAt)
	}

	return &pb.GetStakeConcentrationResponse{
		StakeConcentration: &pb.StakeConcentration{
			Id:               pbID.Uint64(),
			StockID:          pbStockID,
			Date:             pbDate,
			SumBuyShares:     pbSumBuyShares,
			SumSellShares:    pbSumSellShares,
			AvgBuyPrice:      pbAvgBuyPrice,
			AvgSellPrice:     pbAvgSellPrice,
			Concentration_1:  pbConcentration1,
			Concentration_5:  pbConcentration5,
			Concentration_10: pbConcentration10,
			Concentration_20: pbConcentration20,
			Concentration_60: pbConcentration60,
			CreatedAt:        pbCreatedAt,
			UpdatedAt:        pbUpdatedAt,
			DeletedAt:        pbDeletedAt,
		},
	}
}

func ListThreePrimaryRequestFromPB(in *pb.ListThreePrimaryRequest) *ListThreePrimaryRequest {
	if in == nil {
		return nil
	}
	out := &ListThreePrimaryRequest{
		Offset:       in.Offset,
		Limit:        in.Limit,
		SearchParams: ListThreePrimarySearchParamsFromPB(in.SearchParams),
	}

	return out
}

func ListThreePrimarySearchParamsFromPB(in *pb.ListThreePrimarySearchParams) *ListThreePrimarySearchParams {
	if in == nil {
		return nil
	}

	out := &ListThreePrimarySearchParams{
		Start: in.Start,
	}

	out.StockID = in.StockID
	end := in.End
	if len(end) > 0 {
		out.End = &end
	}
	return out
}

func ListThreePrimaryResponseToPB(in *ListThreePrimaryResponse) *pb.ListThreePrimaryResponse {
	if in == nil {
		return nil
	}

	var entries []*pb.ThreePrimary
	for _, obj := range in.Entries {
		entries = append(entries, ThreePrimaryToPB(obj))
	}

	return &pb.ListThreePrimaryResponse{
		Offset:     in.Offset,
		Limit:      in.Limit,
		TotalCount: in.TotalCount,
		Entries:    entries,
	}
}

func ThreePrimaryToPB(in *entity.ThreePrimary) *pb.ThreePrimary {
	if in == nil {
		return nil
	}

	pbID := in.ID
	pbStockID := in.StockID
	pbDate := in.Date
	pbForeignTradeShares := in.ForeignTradeShares
	pbTrustTradeShares := in.TrustTradeShares
	pbDealerTradeShares := in.DealerTradeShares
	pbHedgingTradeShares := in.HedgingTradeShares

	var pbCreatedAt *timestamp.Timestamp
	if in.CreatedAt != nil {
		pbCreatedAt = timestamppb.New(*in.CreatedAt)
	}

	var pbUpdatedAt *timestamp.Timestamp
	if in.UpdatedAt != nil {
		pbUpdatedAt = timestamppb.New(*in.UpdatedAt)
	}

	var pbDeletedAt *timestamp.Timestamp
	if in.DeletedAt != nil {
		pbDeletedAt = timestamppb.New(*in.DeletedAt)
	}

	return &pb.ThreePrimary{
		Id:                 pbID.Uint64(),
		StockID:            pbStockID,
		Date:               pbDate,
		ForeignTradeShares: pbForeignTradeShares,
		TrustTradeShares:   pbTrustTradeShares,
		DealerTradeShares:  pbDealerTradeShares,
		HedgingTradeShares: pbHedgingTradeShares,
		CreatedAt:          pbCreatedAt,
		UpdatedAt:          pbUpdatedAt,
		DeletedAt:          pbDeletedAt,
	}
}
