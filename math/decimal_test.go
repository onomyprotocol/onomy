package types

import (
	"math"
	"reflect"
	"testing"
)

func TestDec_Pow(t *testing.T) {
	tests := []struct {
		name   string
		d 	   Dec
		pow    int64
		want   Dec
	}{
		{"0 value", NewDec(0), 2, NewDec(0) },
		{"negative power", NewDec(2), -2, NewDecFromFloat64(0.25) },
		{"1 with power 2", NewDec(1), 2, NewDec(1) },
		{"2 with power 3", NewDec(2), 3, NewDec(8) },
		{"0 power", NewDec(2), 0, NewDec(1) },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Pow(tt.pow); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDec_Fact(t *testing.T) {
	tests := []struct {
		name    string
		d       Dec
		want    Dec
		wantErr bool
	}{
		{"0", NewDec(0), NewDec(1), false },
		{"1", NewDec(1), NewDec(1), false },
		{"2", NewDec(2), NewDec(2), false },
		{"3", NewDec(3), NewDec(6), false },
		{"10", NewDec(10), NewDec(3628800), false },
		{"-1", NewDec(-1), NewDec(0), true },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.d.Fact()
			if (err != nil) != tt.wantErr {
				t.Errorf("Fact() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fact() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDec_Exp(t *testing.T) {
	precision := struct {
		from Dec
		to Dec
	}{
		NewDecFromFloat64(99.9), NewDecFromFloat64(100.1),
	}
	tests := []struct {
		name   string
		d      Dec
		want   Dec
	}{
		{"0 value", NewDec(0), NewDec(1) },
		{"1 value", NewDec(1), NewDecFromFloat64(2.718281828459045) },
		{"2 value", NewDec(2), NewDecFromFloat64(7.38905609893065) },
		{"10 value", NewDec(10), NewDecFromFloat64(22026.465794806718) },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Exp(); tt.want.MulInt64(100).Quo(got).LT(precision.from) || tt.want.MulInt64(100).Quo(got).GT(precision.to) {
				t.Errorf("Exp() = %v, prec %v want from %v to %v", got, tt.want.MulInt64(100).Quo(got), precision.from, precision.to)
			}
		})
	}
}

func TestDec_Ln(t *testing.T) {
	precision := struct {
		from Dec
		to Dec
	}{
		NewDecFromFloat64(99.9), NewDecFromFloat64(100.1),
	}

	tests := []struct {
		name    string
		d       Dec
		want    Dec
		wantErr bool
	}{
		{"0", NewDec(0), NewDec(0), true },
		{"1", NewDec(1), NewDec(0), false },
		{"1", NewDecFromFloat64(math.E), NewDec(1), false },
		{"2", NewDec(2), NewDecFromFloat64(0.69314718055995), false },
		{"3", NewDec(3), NewDecFromFloat64(1.0986122886681), false },
		{"10", NewDec(10), NewDecFromFloat64(2.302585092994), false },
		{"-1", NewDec(-1), NewDec(0), true },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.Ln()
			if (err != nil) != tt.wantErr {
				t.Errorf("Ln() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got.GT(ZeroDec()) && (tt.want.MulInt64(100).Quo(got).LT(precision.from) || tt.want.MulInt64(100).Quo(got).GT(precision.to)) {
				t.Errorf("Ln() = %v, prec %v want from %v to %v", got, tt.want.MulInt64(100).Quo(got), precision.from, precision.to)
			}
		})
	}
}

func TestDec_RoundInflation(t *testing.T) {
	tests := []struct {
		name   string
		d      Dec
		a      Dec
		b      Dec
		c      Dec
		want   Dec
	}{
		{"1", NewDec(0), NewDec(100), NewDec(150000000), NewDec(50000000), NewDec(1) },
		{"2", NewDec(10000000), NewDec(100), NewDec(150000000), NewDec(50000000), NewDec(2) },
		{"3", NewDec(20000000), NewDec(100), NewDec(150000000), NewDec(50000000), NewDec(3) },
		{"4", NewDec(100000000), NewDec(100), NewDec(150000000), NewDec(50000000), NewDec(61) },
		{"5", NewDec(150000000), NewDec(100), NewDec(150000000), NewDec(50000000), NewDec(100) },
		{"6", NewDec(270000000), NewDec(100), NewDec(150000000), NewDec(50000000), NewDec(6) },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.RoundInflation(tt.a, tt.b, tt.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Inflation() = %v, want %v", got, tt.want)
			}
		})
	}
}

