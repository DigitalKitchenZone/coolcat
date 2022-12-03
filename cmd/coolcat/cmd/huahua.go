package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/spf13/cobra"
)

type HuahuaSnapshot struct {
	Accounts map[string]HuahuaSnapshotAccount `json:"accounts"`
}
type HuahuaSnapshotAccount struct {
	HuahuaAddress       string `json:"huahua_address"`
	HuahuaStaker        bool   `json:"huahua_staker"`
}

func ExportHuahuaSnapshotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export-huahua-snapshot [input-huahua-genesis-file] [output-json]",
		Short: "Export snapshot from a provided Chihuahua Chain genesis export",
		Long: `Export snapshot from a provided Chihuahua Chain genesis export
Example:
	coolcatd export-huahua-snapshot genesis_huahua.json huahua_snapshot.json
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			genesisFile := args[0]
			snapshotOutput := args[1]

			// Read genesis file
			genesisJSON, err := os.Open(genesisFile)
			if err != nil {
				return err
			}
			defer genesisJSON.Close()

			// Produce the map of address to total huahua balance, both staked and unstaked
			snapshotAccs := make(map[string]HuahuaSnapshotAccount)

			cdc := clientCtx.Codec

			appState, _, error := genutiltypes.GenesisStateFromGenFile(genesisFile)
			if error != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			stakingGenState := stakingtypes.GetGenesisStateFromAppState(cdc, appState)

			// Make a map from validator operator address to the validator type
			validators := make(map[string]stakingtypes.Validator)
			for _, validator := range stakingGenState.Validators {
				validators[validator.OperatorAddress] = validator
			}
			amounts := make(map[string]sdk.Dec)
			stakers := 0
			for _, delegation := range stakingGenState.Delegations {
				val, ok := validators[delegation.ValidatorAddress]
				if !ok {
					panic(fmt.Sprintf("missing validator %s ", delegation.GetValidatorAddr()))
				}

				address := delegation.DelegatorAddress
				delegationAmount := val.TokensFromShares(delegation.Shares).Quo(sdk.NewDec(1_000_000))
				current, ok := amounts[address]
				if !ok {
					current = sdk.ZeroDec()
				}
				newAmount := current.Add(delegationAmount)
				amounts[address] = newAmount

				acc, ok := snapshotAccs[address]
				if !ok {
					acc = HuahuaSnapshotAccount{
						HuahuaAddress: address,
					}
				}
				staker := false
				if newAmount.GTE(sdk.NewDecFromIntWithPrec(sdk.NewInt(694200), 0)) {
					acc.HuahuaStaker = true
					staker = true
					stakers++
				}
				if staker {
					snapshotAccs[address] = acc
				}
			}

			snapshot := HuahuaSnapshot{
				Accounts: snapshotAccs,
			}

			fmt.Println("=== Huahua export ===")
			fmt.Printf("accounts: %d\n", len(snapshotAccs))
			fmt.Printf("stakers: %d\n", stakers)

			// export snapshot json
			snapshotJSON, err := json.MarshalIndent(snapshot, "", "    ")
			if err != nil {
				return fmt.Errorf("failed to marshal snapshot: %w", err)
			}

			err = os.WriteFile(snapshotOutput, snapshotJSON, 0600)
			return err
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}