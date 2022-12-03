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

type JunoSnapshot struct {
	Accounts map[string]JunoSnapshotAccount `json:"accounts"`
}

// JunoSnapshotAccount provide fields of snapshot per account
type JunoSnapshotAccount struct {
	JunoAddress       string `json:"juno_address"`
	JunoStaker        bool   `json:"juno_staker"`
	PropSixteen       bool   `json:"prop_sixteen"`
}

// ExportjunoSnapshotCmd generates a snapshot.json from a provided Juno genesis export.
func ExportJunoSnapshotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export-juno-snapshot [juno-denom] [input-juno-genesis-file] [output-json]",
		Short: "Export snapshot from a provided Juno genesis export",
		Long: `Export snapshot from a provided Juno genesis export
Example:
	starsd export-juno-snapshot genesis.json juno-snapshot.json
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

			// Produce the map of address to total atom balance, both staked and unstaked
			snapshotAccs := make(map[string]JunoSnapshotAccount)

			cdc := clientCtx.Codec

			appState, _, error := genutiltypes.GenesisStateFromGenFile(genesisFile)
			if error != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			// export EXCHANGES="cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf,cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll,cosmosvaloper1nm0rrq86ucezaf8uj35pq9fpwr5r82clzyvtd8,cosmosvaloper1kn3wugetjuy4zetlq6wadchfhvu3x740ae6z6x,cosmosvaloper1te8nxpc2myjfrhaty0dnzdhs5ahdh5agzuym9v,cosmosvaloper19yy989ka5usws6gsd8vl94y7l6ssgdwsrnscjc,cosmosvaloper12w6tynmjzq4l8zdla3v4x0jt8lt4rcz5gk7zg2,cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42,cosmosvaloper1z66j0z75a9flwnez7sa8jxx46cqu4rfhd9q82w"
			// copy this line for cosmos hub exchanges

			// prop16votes := strings.Split(strings.TrimSpace(os.Getenv("PROP16VOTES")), ",")
			// fmt.Println("prop16 voters bonussed(?)")
			// if strings.TrimSpace(os.Getenv("PROP16VOTES")) == "" {
			// 	panic("provide list of prop 16 voter addresses")
			// }
			stakingGenState := stakingtypes.GetGenesisStateFromAppState(cdc, appState)

			// Make a map from validator operator address to the validator type
			validators := make(map[string]stakingtypes.Validator)
			for _, validator := range stakingGenState.Validators {
				validators[validator.OperatorAddress] = validator
			}
			amounts := make(map[string]sdk.Dec)
			stakers := 0
			bonusReceivers := 0
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
					acc = JunoSnapshotAccount{
						JunoAddress: address,
					}
				}
				staker := false
				if newAmount.GTE(sdk.NewDecFromIntWithPrec(sdk.NewInt(555), 2)) {
					acc.JunoStaker = true
					staker = true
					stakers++
				}
				// propSixteen := false
				// if isIn(delegation.DelegatorAddress, prop16votes) {
				// 	propSixteen = true
				// 	bonusReceivers++
				// }
				if staker {
					snapshotAccs[address] = acc
				}
			}

			snapshot := JunoSnapshot{
				Accounts: snapshotAccs,
			}

			fmt.Println("=== Juno export ===")
			fmt.Printf("accounts: %d\n", len(snapshotAccs))
			fmt.Printf("stakers: %d\n", stakers)
			fmt.Printf("prop16 voters: %d\n", bonusReceivers)

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