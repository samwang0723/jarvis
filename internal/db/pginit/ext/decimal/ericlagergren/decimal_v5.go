package ericlagergren

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"reflect"

	"github.com/ericlagergren/decimal"
	"github.com/jackc/pgx/v5/pgtype"
)

type Decimal decimal.Big

func (d *Decimal) ScanNumeric(val pgtype.Numeric) error {
	if !val.Valid {
		return ScanNumericError{Val: "NULL"}
	}

	if val.NaN {
		return ScanNumericError{Val: "NaN"}
	}

	if val.InfinityModifier != pgtype.Finite {
		return ScanNumericError{Val: val.InfinityModifier}
	}

	err := d.Compose(0, val.Int.Sign() < 0, val.Int.Bytes(), val.Exp)
	if err != nil {
		return ComposeError{Val: val}
	}

	return nil
}

//nolint:gocritic // cannot use pointer
func (d Decimal) NumericValue() (pgtype.Numeric, error) {
	_, negative, coeff, exp := d.Decompose(nil)

	numericInt := new(big.Int).SetBytes(coeff)
	if negative {
		numericInt = new(big.Int).Neg(numericInt)
	}

	return pgtype.Numeric{Int: numericInt, Exp: exp, Valid: true}, nil
}

// ported from https://github.com/ericlagergren/decimal/blob/495c53812d05/decomposer.go#L43
func (d *Decimal) Decompose(
	buf []byte,
) (form byte, negative bool, coefficient []byte, exponent int32) {
	dd := (*decimal.Big)(d)

	negative = dd.Sign() < 0

	switch {
	case dd.IsInf(0):
		form = 1

		return
	case dd.IsNaN(0):
		form = 2

		return
	}

	if !dd.IsFinite() {
		panic("expected number to be finite")
	}

	exp := -dd.Scale()
	if exp > math.MaxInt32 {
		panic("exponent exceeds max size")
	}

	exponent = int32(exp)

	compact, unscaled := decimal.Raw(dd)

	if d.isCompact() {
		if cap(buf) >= 8 {
			coefficient = buf[:8]
		} else {
			coefficient = make([]byte, 8)
		}

		binary.BigEndian.PutUint64(coefficient, *compact)
	} else {
		coefficient = unscaled.Bytes() // This returns a big-endian slice.
	}

	return
}

// ported from https://github.com/ericlagergren/decimal/blob/495c53812d05/decomposer.go#L76
func (d *Decimal) Compose(form byte, negative bool, coefficient []byte, exponent int32) error {
	dd := (*decimal.Big)(d)

	switch form {
	default:
		return fmt.Errorf("unknown form: %v", form)
	case 0:
		// Finite form below.
	case 1:
		dd.SetInf(negative)

		return nil
	case 2:
		dd.SetNaN(false)

		return nil
	}

	bigc := &big.Int{}
	bigc.SetBytes(coefficient)

	dd.SetBigMantScale(bigc, -int(exponent))

	if negative {
		dd.Neg(dd)
	}

	return nil
}

func (d *Decimal) isCompact() bool {
	dd := decimal.Big(*d)
	compact, _ := decimal.Raw(&dd)

	return *compact != math.MaxUint64
}

func TryWrapNumericEncodePlan( //nolint: ireturn // ref shopspring
	value interface{},
) (pgtype.WrappedEncodePlanNextSetter, interface{}, bool) {
	if fmt.Sprintf("%T", value) == "decimal.Big" {
		return &wrapDecimalEncodePlan{}, Decimal(
			value.(decimal.Big),
		), true //nolint:forcetypeassert // ref shopspring
	}

	return nil, nil, false
}

type wrapDecimalEncodePlan struct {
	next pgtype.EncodePlan
}

func (plan *wrapDecimalEncodePlan) SetNext(next pgtype.EncodePlan) { plan.next = next }

func (plan *wrapDecimalEncodePlan) Encode(value interface{}, buf []byte) ([]byte, error) {
	return plan.next.Encode(
		Decimal(value.(decimal.Big)),
		buf,
	) //nolint:forcetypeassert,wrapcheck // ref shopspring
}

func TryWrapNumericScanPlan( //nolint: ireturn // ref shopspring
	target interface{},
) (pgtype.WrappedScanPlanNextSetter, interface{}, bool) {
	if fmt.Sprintf("%T", target) == "*decimal.Big" {
		return &wrapDecimalScanPlan{},
			(*Decimal)(target.(*decimal.Big)), true //nolint:forcetypeassert // ref shopspring
	}

	return nil, nil, false
}

type wrapDecimalScanPlan struct {
	next pgtype.ScanPlan
}

func (plan *wrapDecimalScanPlan) SetNext(next pgtype.ScanPlan) { plan.next = next }

func (plan *wrapDecimalScanPlan) Scan(src []byte, dst interface{}) error {
	return plan.next.Scan(
		src,
		(*Decimal)(dst.(*decimal.Big)),
	) //nolint:forcetypeassert,wrapcheck // ref shopspring
}

type NumericCodec struct {
	pgtype.NumericCodec
}

func (NumericCodec) DecodeValue(typeMap *pgtype.Map,
	oid uint32,
	format int16,
	src []byte,
) (interface{}, error) {
	if src == nil {
		return nil, nil //nolint: nilnil // a
	}

	var target decimal.Big
	scanPlan := typeMap.PlanScan(oid, format, &target)

	if scanPlan == nil {
		return nil, NoPlanError{}
	}

	err := scanPlan.Scan(src, &target)
	if err != nil {
		return nil, ScanError{Err: err}
	}

	return target, nil
}

// Register registers the ericlagergren/decimal integration with a pgtype.ConnInfo.
func Register(typeMap *pgtype.Map) {
	typeMap.TryWrapEncodePlanFuncs = append(
		[]pgtype.TryWrapEncodePlanFunc{TryWrapNumericEncodePlan},
		typeMap.TryWrapEncodePlanFuncs...)
	typeMap.TryWrapScanPlanFuncs = append([]pgtype.TryWrapScanPlanFunc{TryWrapNumericScanPlan},
		typeMap.TryWrapScanPlanFuncs...)

	typeMap.RegisterType(&pgtype.Type{
		Name:  "numeric",
		OID:   pgtype.NumericOID,
		Codec: NumericCodec{},
	})

	registerDefaultPgTypeVariants := func(name, arrayName string, value interface{}) {
		// T
		typeMap.RegisterDefaultPgType(value, name)

		// *T
		valueType := reflect.TypeOf(value)
		typeMap.RegisterDefaultPgType(reflect.New(valueType).Interface(), name)

		// []T
		sliceType := reflect.SliceOf(valueType)
		typeMap.RegisterDefaultPgType(reflect.MakeSlice(sliceType, 0, 0).Interface(), arrayName)

		// *[]T
		typeMap.RegisterDefaultPgType(reflect.New(sliceType).Interface(), arrayName)

		// []*T
		sliceOfPointerType := reflect.SliceOf(reflect.TypeOf(reflect.New(valueType).Interface()))
		typeMap.RegisterDefaultPgType(
			reflect.MakeSlice(sliceOfPointerType, 0, 0).Interface(),
			arrayName,
		)

		// *[]*T
		typeMap.RegisterDefaultPgType(reflect.New(sliceOfPointerType).Interface(), arrayName)
	}

	registerDefaultPgTypeVariants("numeric", "_numeric", decimal.Big{})
	registerDefaultPgTypeVariants("numeric", "_numeric", Decimal{})
}
