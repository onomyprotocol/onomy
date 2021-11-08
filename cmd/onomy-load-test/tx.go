//go:build tmload
// +build tmload

package main

import (
	"github.com/cosmos/cosmos-sdk/client/tx"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	_ "github.com/onomyprotocol/onomy/app"
)

const (
	coin          = "footoken"
	coinAmount    = 1
	gasLimit      = 200000
	txFees        = 5
	memo          = ""
	timeoutHeight = 5
	chainID       = "onomy"
	accountNo     = 0
	sequence      = 0
)

// GenTx return transactions bytes .
func GenTx() ([]byte, error) {
	encCfg := simapp.MakeTestEncodingConfig()

	// Create a new TxBuilder
	txBuilder := encCfg.TxConfig.NewTxBuilder()
	priv1, _, addr1 := testdata.KeyTestPubAddr()
	_, _, addr2 := testdata.KeyTestPubAddr()
	msg1 := banktypes.NewMsgSend(addr1, addr2, sdk.NewCoins(sdk.NewInt64Coin(coin, coinAmount)))
	err := txBuilder.SetMsgs(msg1)
	if err != nil {
		return nil, err
	}

	txBuilder.SetGasLimit(gasLimit)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewInt64Coin(coin, txFees)))
	txBuilder.SetMemo(memo)
	txBuilder.SetTimeoutHeight(timeoutHeight)
	//------------------------------- ------Signing a Transaction---------------------------------//

	privs := []cryptotypes.PrivKey{priv1}
	var sigsV2 []signing.SignatureV2
	for _, priv := range privs {
		sigV2 := signing.SignatureV2{
			PubKey: priv.PubKey(),
			Data: &signing.SingleSignatureData{
				SignMode:  encCfg.TxConfig.SignModeHandler().DefaultMode(),
				Signature: nil,
			},

			Sequence: 0,
		}
		sigsV2 = append(sigsV2, sigV2)
	}
	err = txBuilder.SetSignatures(sigsV2...)
	if err != nil {
		return nil, err
	}
	sigsV2 = []signing.SignatureV2{}

	for _, priv := range privs {
		signerData := xauthsigning.SignerData{
			ChainID:       chainID,
			AccountNumber: accountNo,
			Sequence:      sequence,
		}
		sigV2, err := tx.SignWithPrivKey(
			encCfg.TxConfig.SignModeHandler().DefaultMode(), signerData,
			txBuilder, priv, encCfg.TxConfig, 0)
		if err != nil {
			return nil, err
		}
		sigsV2 = append(sigsV2, sigV2)
	}

	err = txBuilder.SetSignatures(sigsV2...)
	if err != nil {
		return nil, err
	}

	txBytes, err := encCfg.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}
	return txBytes, nil
}
