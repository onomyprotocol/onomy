// Package cli contains dao module cli utils.
package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	govcli "github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/onomyprotocol/onomy/x/dao/types"
)

type proposalGeneric struct {
	Title       string
	Description string
	Deposit     string
}

// GetTxCmd returns the transaction commands for this module.
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s dao subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2, // nolint:gomnd
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdFundTreasuryProposal())
	cmd.AddCommand(CmdExchangeWithTreasuryProposal())
	cmd.AddCommand(CmdFundAccountProposal())

	return cmd
}

// CmdFundTreasuryProposal implements the command to submit a fund-treasury proposal.
func CmdFundTreasuryProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fund-treasury amount",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a fund treasury proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a fund treasury proposal along with an initial deposit.
Example:
$ %s tx gov submit-proposal fund-treasury 5000000000000000000anom --title="Test Proposal" --description="My awesome proposal" --deposit="10000000000000000000anom" --from mykey`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinsNormalized(args[0])
			if err != nil {
				return err
			}

			proposalGeneric, err := parseSubmitProposalFlags(cmd.Flags())
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(proposalGeneric.Deposit)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			content := types.NewFundTreasuryProposal(from, proposalGeneric.Title, proposalGeneric.Description, amount)

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	addProposalFlags(cmd)

	return cmd
}

// CmdExchangeWithTreasuryProposal implements the command to submit еру exchange-with-treasury proposal.
func CmdExchangeWithTreasuryProposal() *cobra.Command { // nolint:gocognit,gocyclo,cyclop // the command is long but not complicated
	cmd := &cobra.Command{
		Use:   "exchange-with-treasury coins-pairs",
		Args:  cobra.ExactArgs(1),
		Short: "Submit an exchange with treasury proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submitan exchange with treasury proposal along with an initial deposit.
Example (the anom in the example is ask coin, ousd is bid coin):
$ %s tx gov submit-proposal exchange-with-treasury "5000000000000000000anom/5000000000000000000ousd" --title="Test Proposal" --description="My awesome proposal" --deposit="10000000000000000000anom" --from mykey`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			exchangeCoins := make([]types.CoinsExchangePair, 0)
			pairsStrings := strings.Split(args[0], ",")
			for i := range pairsStrings {
				coinsString := strings.Split(pairsStrings[i], "/")
				if len(coinsString) != 2 { // nolint:gomnd // pair number
					return fmt.Errorf("coins pair %s is invalid", pairsStrings[i])
				}

				coinAsk, err := sdk.ParseCoinNormalized(coinsString[0])
				if err != nil {
					return err
				}
				if err := coinAsk.Validate(); err != nil {
					return err
				}
				coinBid, err := sdk.ParseCoinNormalized(coinsString[1])
				if err != nil {
					return err
				}
				if err := coinBid.Validate(); err != nil {
					return err
				}

				exchangeCoins = append(exchangeCoins, types.CoinsExchangePair{
					CoinAsk: coinAsk,
					CoinBid: coinBid,
				})
			}

			proposalGeneric, err := parseSubmitProposalFlags(cmd.Flags())
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(proposalGeneric.Deposit)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			content := types.NewExchangeWithTreasuryProposal(from, proposalGeneric.Title, proposalGeneric.Description, exchangeCoins)

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	addProposalFlags(cmd)

	return cmd
}

// CmdFundAccountProposal implements the command to submit a fund-account proposal.
func CmdFundAccountProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fund-account recipient amount",
		Args:  cobra.ExactArgs(2), // nolint:gomnd // the args count
		Short: "Submit a fund account proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a fund account proposal along with an initial deposit.
Example:
$ %s tx gov submit-proposal fund-account recipient-address 5000000000000000000anom --title="Test Proposal" --description="My awesome proposal" --deposit="10000000000000000000anom" --from mykey`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			recipient, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}

			proposalGeneric, err := parseSubmitProposalFlags(cmd.Flags())
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(proposalGeneric.Deposit)
			if err != nil {
				return err
			}

			content := types.NewFundAccountProposal(recipient, proposalGeneric.Title, proposalGeneric.Description, amount)

			from := clientCtx.GetFromAddress()
			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	addProposalFlags(cmd)

	return cmd
}

func parseSubmitProposalFlags(fs *pflag.FlagSet) (*proposalGeneric, error) {
	title, err := fs.GetString(govcli.FlagTitle)
	if err != nil {
		return nil, err
	}
	description, err := fs.GetString(govcli.FlagDescription)
	if err != nil {
		return nil, err
	}

	deposit, err := fs.GetString(govcli.FlagDeposit)
	if err != nil {
		return nil, err
	}

	return &proposalGeneric{
		Title:       title,
		Description: description,
		Deposit:     deposit,
	}, nil
}

func addProposalFlags(cmd *cobra.Command) {
	cmd.Flags().String(govcli.FlagTitle, "", "The proposal title")
	cmd.Flags().String(govcli.FlagDescription, "", "The proposal description")
	cmd.Flags().String(govcli.FlagDeposit, "", "The proposal deposit")
}
