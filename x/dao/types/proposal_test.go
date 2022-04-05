package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/onomyprotocol/onomy/testutil/simapp"
	"github.com/onomyprotocol/onomy/x/dao/types"
)

func TestFundTreasuryProposal_ValidateBasic(t *testing.T) {
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
