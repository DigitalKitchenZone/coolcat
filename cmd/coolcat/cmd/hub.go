package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/spf13/cobra"
)

type HubSnapshot struct {
	Accounts map[string]HubSnapshotAccount `json:"accounts"`
}

// HubSnapshotAccount provide fields of snapshot per account
type HubSnapshotAccount struct {
	AtomAddress       string `json:"atom_address"`
	AtomStaker        bool   `json:"atom_staker"`
	OutsideTopTwenty bool   `json:"atom_bonus"`
}

func isIn(s string, ss []string) bool {
	for _, t := range ss {
		if t == s {
			return true
		}
	}
	return false
}

// ExportHubSnapshotCmd generates a snapshot.json from a provided Cosmos Hub genesis export.
func ExportHubSnapshotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export-hub-snapshot [airdrop-to-denom] [input-genesis-file] [output-snapshot-json]",
		Short: "Export snapshot from a provided Cosmos Hub genesis export",
		Long: `Export snapshot from a provided Cosmos Hub genesis export
Example:
	starsd export-hub-snapshot genesis.json hub-snapshot.json
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
			snapshotAccs := make(map[string]HubSnapshotAccount)

			cdc := clientCtx.Codec

			appState, _, error := genutiltypes.GenesisStateFromGenFile(genesisFile)
			if error != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			exchanges := strings.Split(strings.TrimSpace(os.Getenv("EXCHANGES")), ",")
			fmt.Println("exchanges", len(exchanges))
			if len(exchanges) == 0 || strings.TrimSpace(os.Getenv("EXCHANGES")) == "" {
				panic("provide list of addresses")
			}
			stakingGenState := stakingtypes.GetGenesisStateFromAppState(cdc, appState)

			// Make a map from validator operator address to the validator type
			validators := make(map[string]stakingtypes.Validator)
			for _, validator := range stakingGenState.Validators {
				validators[validator.OperatorAddress] = validator
			}

			sortedValidators := make(stakingtypes.ValidatorsByVotingPower, 0, len(validators))
			for _, sortedValidator := range validators {
				sortedValidators = append(sortedValidators, sortedValidator)
			}

			sort.Slice(sortedValidators, func(i, j int) bool {
        		// return sortedValidators[i].Tokens > sortedValidators[j].Tokens
				return sortedValidators[i].Tokens.GT(sortedValidators[j].Tokens) 
    		})

			withoutTopTwenty := make(map[sdk.Int]stakingtypes.Validator)
			for i, k := range sortedValidators {
				if(i <= 20){
					continue
				}
				withoutTopTwenty[k.Tokens] = k
			}

			finalValidators := make(map[string]stakingtypes.Validator)
			for _, finalVal := range withoutTopTwenty {
				finalValidators[finalVal.OperatorAddress] = finalVal
			}
	 
			amounts := make(map[string]sdk.Dec)
			stakers := 0
			bonusReceivers := 0
			for _, delegation := range stakingGenState.Delegations {
				if isIn(delegation.ValidatorAddress, exchanges) {
					continue
				}

				val, ok := validators[delegation.ValidatorAddress]
				if !ok {
					panic(fmt.Sprintf("missing validator %s ", delegation.GetValidatorAddr()))
				}

				outsideTwentyVal, ok := finalValidators[delegation.ValidatorAddress]
				if !ok {
					outsideTwentyVal.OperatorAddress = "top20"
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
					acc = HubSnapshotAccount{
						AtomAddress: address,
						AtomStaker: false,
						OutsideTopTwenty: false,
					}
				}

				staker := false
				if newAmount.GTE(sdk.NewDecFromIntWithPrec(sdk.NewInt(777), 2)) {
					acc.AtomStaker = true
					staker = true
					stakers++
				}

				outsideTopTwenty := false
				if newAmount.GTE(sdk.NewDecFromIntWithPrec(sdk.NewInt(777), 2)) && outsideTwentyVal.OperatorAddress != "top20" {
					acc.OutsideTopTwenty = true
					outsideTopTwenty = true
					bonusReceivers++
				}
				if staker || outsideTopTwenty {
					snapshotAccs[address] = acc
				}
			}

			snapshot := HubSnapshot{
				Accounts: snapshotAccs,
			}

			fmt.Println("=== Cosmos Hub export ===")
			fmt.Printf("accounts: %d\n", len(snapshotAccs))
			fmt.Printf("stakers: %d\n", stakers)
			fmt.Printf("outsideTwenty: %d\n", bonusReceivers)

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