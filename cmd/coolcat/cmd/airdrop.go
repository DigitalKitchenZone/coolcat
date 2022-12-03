package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	catdroptypes "github.com/DigitalKitchenLabs/coolcat/v1/x/catdrop/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/cobra"
)

const Denom = "uccat"

func AddAirdropCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-airdrop [airdrop-snapshot-file]",
		Short: "Add balances of accounts to claim module.",
		Args:  cobra.ExactArgs(1),
		Long: `Add balances of accounts to claim module.
				Example:
				coolcatd add-airdrop /path/to/snapshot.json
				`,
		RunE: func(cmd *cobra.Command, args []string) error {

			var clientCtx = client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.Codec

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			// read snapshot
			snapshotFile := args[0]
			snapshotJSON, _ := os.ReadFile(snapshotFile)
			snapshot := CoolCatSnapshot{}
			json.Unmarshal([]byte(snapshotJSON), &snapshot)

			genFile := config.GenesisFile()

			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			fmt.Printf("Accounts %d\n", len(snapshot.Accounts))

			claimGenState := catdroptypes.GetGenesisStateFromAppState(cdc, appState)

			totalClaimedAmount := sdk.NewInt64Coin(Denom, 0)
			for address, acc := range snapshot.Accounts {
				// empty account check
				if acc.AirdropAmount.LTE(sdk.NewInt(0)) {
					panic("Empty account")
				}

				coin := sdk.NewCoin(Denom, acc.AirdropAmount)
				coins := sdk.NewCoins(coin)

				record := catdroptypes.ClaimRecord{
					Address:                address,
					InitialClaimableAmount: coins,
					ActionCompleted:        []bool{false, false},
				}
				claimGenState.ClaimRecords = append(claimGenState.ClaimRecords, record)
				totalClaimedAmount = totalClaimedAmount.Add(coin)
			}

			claimGenState.ModuleAccountBalance = totalClaimedAmount
			claimGenStateBz, err := cdc.MarshalJSON(claimGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal claim genesis state: %w", err)
			}
			appState[catdroptypes.ModuleName] = claimGenStateBz

			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON

			fmt.Printf("Saving genesis...")
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	return cmd
}