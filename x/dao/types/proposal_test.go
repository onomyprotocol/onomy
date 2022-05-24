package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/onomyprotocol/onomy/testutil/simapp"
	"github.com/onomyprotocol/onomy/x/dao/types"
)

func TestFundTreasuryProposal_ValidateBasic(t *testing.T) { //nolint:dupl // test template
	const denom1 = "denom1"

	type fields struct {
		Sender      string
		Title       string
		Description string
		Amount      sdk.Coins
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "positive",
			fields: fields{
				Sender:      simapp.GenAccount().String(),
				Title:       "title",
				Description: "desc",
				Amount:      sdk.NewCoins(sdk.NewInt64Coin(denom1, 10)),
			},
		},
		{
			name: "negative_invalid_sender",
			fields: fields{
				Sender:      "invalid-sender",
				Title:       "title",
				Description: "desc",
				Amount:      sdk.NewCoins(sdk.NewInt64Coin(denom1, 10)),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			m := &types.FundTreasuryProposal{
				Sender:      tt.fields.Sender,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				Amount:      tt.fields.Amount,
			}
			if err := m.ValidateBasic(); (err != nil) != tt.wantErr {
				t.Errorf("ValidateBasic() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExchangeWithTreasuryProposal_ValidateBasic(t *testing.T) {
	const (
		denom1 = "denom1"
		denom2 = "denom2"
	)

	type fields struct {
		Sender      string
		Title       string
		Description string
		CoinsPairs  []types.CoinsExchangePair
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "positive",
			fields: fields{
				Sender:      simapp.GenAccount().String(),
				Title:       "title",
				Description: "desc",
				CoinsPairs: []types.CoinsExchangePair{
					{
						CoinAsk: sdk.NewInt64Coin(denom1, 10),
						CoinBid: sdk.NewInt64Coin(denom2, 5),
					},
					{
						CoinAsk: sdk.NewInt64Coin(denom2, 10),
						CoinBid: sdk.NewInt64Coin(denom1, 5),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "negative_zero_coin",
			fields: fields{
				Sender:      simapp.GenAccount().String(),
				Title:       "title",
				Description: "desc",
				CoinsPairs: []types.CoinsExchangePair{
					{
						CoinAsk: sdk.NewInt64Coin(denom1, 10),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "negative_invalid_sender",
			fields: fields{
				Sender:      "invalid-sender",
				Title:       "title",
				Description: "desc",
				CoinsPairs: []types.CoinsExchangePair{
					{
						CoinAsk: sdk.NewInt64Coin(denom2, 10),
						CoinBid: sdk.NewInt64Coin(denom1, 5),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			m := &types.ExchangeWithTreasuryProposal{
				Sender:      tt.fields.Sender,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				CoinsPairs:  tt.fields.CoinsPairs,
			}
			if err := m.ValidateBasic(); (err != nil) != tt.wantErr {
				t.Errorf("ValidateBasic() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFundAccountProposal_ValidateBasic(t *testing.T) { // nolint:dupl // test template
	const denom1 = "denom1"

	type fields struct {
		Recipient   string
		Title       string
		Description string
		Amount      sdk.Coins
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "positive",
			fields: fields{
				Recipient:   simapp.GenAccount().String(),
				Title:       "title",
				Description: "desc",
				Amount:      sdk.NewCoins(sdk.NewInt64Coin(denom1, 10)),
			},
		},
		{
			name: "negative_invalid_sender",
			fields: fields{
				Recipient:   "invalid-sender",
				Title:       "title",
				Description: "desc",
				Amount:      sdk.NewCoins(sdk.NewInt64Coin(denom1, 10)),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			m := &types.FundAccountProposal{
				Recipient:   tt.fields.Recipient,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				Amount:      tt.fields.Amount,
			}
			if err := m.ValidateBasic(); (err != nil) != tt.wantErr {
				t.Errorf("ValidateBasic() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
