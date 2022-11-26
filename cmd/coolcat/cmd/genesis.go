package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/spf13/cobra"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"

	minttypes "github.com/DigitalKitchenLabs/coolcat/v1/x/mint/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	appParams "github.com/DigitalKitchenLabs/coolcat/v1/app/params"
	alloctypes "github.com/DigitalKitchenLabs/coolcat/v1/x/alloc/types"
	catdroptypes "github.com/DigitalKitchenLabs/coolcat/v1/x/catdrop/types"
)

type GenesisParams struct {
	AirdropSupply sdk.Int

	StrategicReserveAccounts []banktypes.Balance

	ConsensusParams *tmproto.ConsensusParams

	GenesisTime         time.Time
	NativeCoinMetadatas []banktypes.Metadata

	StakingParams      stakingtypes.Params
	DistributionParams distributiontypes.Params
	GovParams          govtypes.Params

	CrisisConstantFee sdk.Coin

	SlashingParams slashingtypes.Params

	AllocParams alloctypes.Params
	ClaimParams catdroptypes.Params
	MintParams  minttypes.Params

	WasmParams wasmtypes.Params
}

func PrepareGenesisCmd(defaultNodeHome string, mbm module.BasicManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prepare-genesis [network] [chainID] [file]",
		Short: "Prepare a genesis file with initial setup",
		Long: `Prepare a genesis file with initial setup.
				Example:
				coolcatd prepare-genesis testnet kitten-02 snapshot.json
				- Check genesis file at: ~/.coolcatd/config/genesis.json
				`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.Codec

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config
			config.SetRoot(clientCtx.HomeDir)

			// read genesis file
			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			// get genesis params
			genesisParams := MainnetGenesisParams()
			switch args[0] {
			case "localnet":
				genesisParams = LocalnetGenesisParams()
			case "testnet":
				genesisParams = TestnetGenesisParams()
			case "devnet":
				genesisParams = DevnetGenesisParams()
			}
			// get genesis params
			chainID := args[1]

			// read snapshot.json and parse into struct
			snapshotFile, _ := ioutil.ReadFile(args[2])
			snapshot := CoolCatSnapshot{}
			err = json.Unmarshal(snapshotFile, &snapshot)
			if err != nil {
				panic(err)
			}

			// run Prepare Genesis
			appState, genDoc, err = PrepareGenesis(clientCtx, appState, genDoc, genesisParams, chainID, snapshot)
			if err != nil {
				return fmt.Errorf("failed to prepare genesis: %w", err)
			}

			// validate genesis state
			if err = mbm.ValidateGenesis(cdc, clientCtx.TxConfig, appState); err != nil {
				return fmt.Errorf("error validating genesis file: %s", err.Error())
			}

			// save genesis
			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON
			err = genutil.ExportGenesisFile(genDoc, genFile)
			return err
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// fill with data
func PrepareGenesis(
	clientCtx client.Context,
	appState map[string]json.RawMessage,
	genDoc *tmtypes.GenesisDoc,
	genesisParams GenesisParams,
	chainID string,
	snapshot CoolCatSnapshot,
) (map[string]json.RawMessage, *tmtypes.GenesisDoc, error) {
	cdc := clientCtx.Codec

	// chain params genesis
	genDoc.GenesisTime = genesisParams.GenesisTime
	genDoc.ChainID = chainID
	genDoc.ConsensusParams = genesisParams.ConsensusParams

	// IBC transfer module genesis
	ibcGenState := ibctransfertypes.DefaultGenesisState()
	ibcGenState.Params.SendEnabled = true
	ibcGenState.Params.ReceiveEnabled = true
	ibcGenStateBz, err := cdc.MarshalJSON(ibcGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal IBC transfer genesis state: %w", err)
	}
	appState[ibctransfertypes.ModuleName] = ibcGenStateBz

	// mint module genesis
	mintGenState := minttypes.DefaultGenesisState()
	mintGenState.Params = genesisParams.MintParams

	mintGenStateBz, err := cdc.MarshalJSON(mintGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal mint genesis state: %w", err)
	}
	appState[minttypes.ModuleName] = mintGenStateBz

	// staking module genesis
	stakingGenState := stakingtypes.GetGenesisStateFromAppState(cdc, appState)
	stakingGenState.Params = genesisParams.StakingParams
	stakingGenStateBz, err := cdc.MarshalJSON(stakingGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal staking genesis state: %w", err)
	}
	appState[stakingtypes.ModuleName] = stakingGenStateBz

	// distribution module genesis
	distributionGenState := distributiontypes.DefaultGenesisState()
	distributionGenState.Params = genesisParams.DistributionParams
	distributionGenStateBz, err := cdc.MarshalJSON(distributionGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal distribution genesis state: %w", err)
	}
	appState[distributiontypes.ModuleName] = distributionGenStateBz

	// gov module genesis
	govGenState := govtypes.DefaultGenesisState()
	govGenState.DepositParams = genesisParams.GovParams.DepositParams
	govGenState.TallyParams = genesisParams.GovParams.TallyParams
	govGenState.VotingParams = genesisParams.GovParams.VotingParams
	govGenStateBz, err := cdc.MarshalJSON(govGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal gov genesis state: %w", err)
	}
	appState[govtypes.ModuleName] = govGenStateBz

	// crisis module genesis
	crisisGenState := crisistypes.DefaultGenesisState()
	crisisGenState.ConstantFee = genesisParams.CrisisConstantFee
	crisisGenStateBz, err := cdc.MarshalJSON(crisisGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal crisis genesis state: %w", err)
	}
	appState[crisistypes.ModuleName] = crisisGenStateBz

	// slashing module genesis
	slashingGenState := slashingtypes.DefaultGenesisState()
	slashingGenState.Params = genesisParams.SlashingParams
	slashingGenStateBz, err := cdc.MarshalJSON(slashingGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal slashing genesis state: %w", err)
	}
	appState[slashingtypes.ModuleName] = slashingGenStateBz

	// auth accounts
	authGenState := authtypes.GetGenesisStateFromAppState(cdc, appState)
	accs, err := authtypes.UnpackAccounts(authGenState.Accounts)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get accounts from any: %w", err)
	}

	// ---
	// bank module genesis
	bankGenState := banktypes.GetGenesisStateFromAppState(cdc, appState)
	bankGenState.Params.DefaultSendEnabled = true
	bankGenState.DenomMetadata = genesisParams.NativeCoinMetadatas
	balances := bankGenState.Balances

	// catdrop module genesis
	claimGenState := catdroptypes.GetGenesisStateFromAppState(cdc, appState)
	claimGenState.Params = genesisParams.ClaimParams
	claimRecords := make([]catdroptypes.ClaimRecord, 0, len(snapshot.Accounts))
	claimsTotal := sdk.ZeroInt()
	// check from preexisint accounts in genesis
	preExistingAccounts := make(map[string]bool)
	for _, b := range balances {
		preExistingAccounts[b.Address] = true
	}
	for addr, acc := range snapshot.Accounts {
		claimRecord := catdroptypes.ClaimRecord{
			Address:                addr,
			InitialClaimableAmount: sdk.NewCoins(sdk.NewCoin(appParams.BaseCoinUnit, acc.AirdropAmount)),
			ActionCompleted:        []bool{false, false, false, false},
		}
		claimsTotal = claimsTotal.Add(acc.AirdropAmount)
		claimRecords = append(claimRecords, claimRecord)
		// skip account addition if existent
		exists := preExistingAccounts[addr]
		if exists {
			continue
		}
		balances = append(balances, banktypes.Balance{
			Address: addr,
			Coins:   sdk.NewCoins(sdk.NewInt64Coin(appParams.BaseCoinUnit, 1_000_000)),
		})

		address, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			return nil, nil, err
		}
		// add base account
		// Add the new account to the set of genesis accounts
		baseAccount := authtypes.NewBaseAccount(address, nil, 0, 0)
		if err := baseAccount.Validate(); err != nil {
			return nil, nil, fmt.Errorf("failed to validate new genesis account: %w", err)
		}
		accs = append(accs, baseAccount)
	}
	claimGenState.ClaimRecords = claimRecords
	claimGenState.ModuleAccountBalance = sdk.NewCoin(appParams.BaseCoinUnit, claimsTotal)
	claimGenStateBz, err := cdc.MarshalJSON(claimGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal claim genesis state: %w", err)
	}
	appState[catdroptypes.ModuleName] = claimGenStateBz

	// save accounts
	balances = append(balances, banktypes.Balance{
			Address: "ccat1z7qzdmefyqtk05m5qd6dhm77n6hjkfh2kv7eqx",
			Coins:   sdk.NewCoins(sdk.NewInt64Coin(appParams.BaseCoinUnit, 1_500_000_000_000_000)),
	})
	balances = append(balances, banktypes.Balance{
			Address: "ccat15myda56ut2hxakzdazuyu3ys94q6g9k0rwzjgx",
			Coins:   sdk.NewCoins(sdk.NewInt64Coin(appParams.BaseCoinUnit, 500_000_000_000_000)),
	})
	balances = append(balances, banktypes.Balance{
			Address: "ccat1nnp8wz7s5k8p38v7rdnptcl8e4acd8u0pjn6f6",
			Coins:   sdk.NewCoins(sdk.NewInt64Coin(appParams.BaseCoinUnit, 1_000_000_000_000_000)),
	})
	balances = append(balances, banktypes.Balance{
			Address: "ccat1cuflqsykgs9h3whhwmp2e97ftwl2vur80q4zah",
			Coins:   sdk.NewCoins(sdk.NewInt64Coin(appParams.BaseCoinUnit, 3_500_000_000_000_000)),
	})

	poolAddress, err := sdk.AccAddressFromBech32("ccat1cuflqsykgs9h3whhwmp2e97ftwl2vur80q4zah")
		if err != nil {
			return nil, nil, err
		}
	poolAccount := authtypes.NewBaseAccount(poolAddress, nil, 0, 0)
		if err := poolAccount.Validate(); err != nil {
			return nil, nil, fmt.Errorf("failed to validate new genesis account: %w", err)
		}
	accs = append(accs, poolAccount)

	airdropTwoAddress, err := sdk.AccAddressFromBech32("ccat1nnp8wz7s5k8p38v7rdnptcl8e4acd8u0pjn6f6")
		if err != nil {
			return nil, nil, err
		}
	airdropTwoAccount := authtypes.NewBaseAccount(airdropTwoAddress, nil, 0, 0)
		if err := airdropTwoAccount.Validate(); err != nil {
			return nil, nil, fmt.Errorf("failed to validate new genesis account: %w", err)
		}
	accs = append(accs, airdropTwoAccount)

	reserveAddress, err := sdk.AccAddressFromBech32("ccat15myda56ut2hxakzdazuyu3ys94q6g9k0rwzjgx")
		if err != nil {
			return nil, nil, err
		}
	reserveAccount := authtypes.NewBaseAccount(reserveAddress, nil, 0, 0)
		if err := airdropTwoAccount.Validate(); err != nil {
			return nil, nil, fmt.Errorf("failed to validate new genesis account: %w", err)
		}
	accs = append(accs, reserveAccount)


	// auth module genesis
	accs = authtypes.SanitizeGenesisAccounts(accs)
	genAccs, err := authtypes.PackAccounts(accs)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert accounts into any's: %w", err)
	}
	authGenState.Accounts = genAccs
	authGenStateBz, err := cdc.MarshalJSON(&authGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal auth genesis state: %w", err)
	}
	appState[authtypes.ModuleName] = authGenStateBz

	// save balances
	bankGenState.Balances = banktypes.SanitizeGenesisBalances(balances)
	bankGenStateBz, err := cdc.MarshalJSON(bankGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal bank genesis state: %w", err)
	}
	appState[banktypes.ModuleName] = bankGenStateBz

	// alloc module genesis
	allocGenState := alloctypes.GetGenesisStateFromAppState(cdc, appState)
	allocGenState.Params = genesisParams.AllocParams
	allocGenStateBz, err := cdc.MarshalJSON(allocGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal alloc genesis state: %w", err)
	}
	appState[alloctypes.ModuleName] = allocGenStateBz

	// wasm
	// wasm module genesis
	wasmGenState := &wasm.GenesisState{
		Params: genesisParams.WasmParams,
	}
	wasmGenStateBz, err := cdc.MarshalJSON(wasmGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal wasm genesis state: %w", err)
	}
	appState[wasm.ModuleName] = wasmGenStateBz

	return appState, genDoc, nil
}

// params only
func MainnetGenesisParams() GenesisParams {
	genParams := GenesisParams{}

	genParams.AirdropSupply = sdk.NewInt(3_500_000_000_000_000)              // 3.5B CCAT
	genParams.GenesisTime = time.Date(2022, 11, 24, 0, 0, 0, 0, time.UTC) // August 22, 2022 - 00:00 UTC

	genParams.NativeCoinMetadatas = []banktypes.Metadata{
		{
			Description: "The native token of the CoolCat ecosystem",
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    appParams.BaseCoinUnit,
					Exponent: 0,
					Aliases:  nil,
				},
				{
					Denom:    appParams.HumanCoinUnit,
					Exponent: appParams.DenomExponent,
					Aliases:  nil,
				},
			},
			Name:    "CoolCat",
			Base:    appParams.BaseCoinUnit,
			Display: appParams.HumanCoinUnit,
			Symbol:  "CCAT",
		},
	}

	// alloc
	genParams.AllocParams = alloctypes.DefaultParams()
	genParams.AllocParams.DistributionProportions = alloctypes.DistributionProportions{
		CommunityPool:    sdk.NewDecWithPrec(10, 2), // 10%
	}

	// mint
	genParams.MintParams = minttypes.DefaultParams()
	genParams.MintParams.MintDenom = appParams.BaseCoinUnit
	genParams.MintParams.BlocksPerYear = uint64(6311520)

	genParams.StakingParams = stakingtypes.DefaultParams()
	genParams.StakingParams.UnbondingTime = time.Hour * 24 * 7 * 2 // 2 weeks
	genParams.StakingParams.MaxValidators = 100
	genParams.StakingParams.BondDenom = genParams.NativeCoinMetadatas[0].Base
	// MinCommissionRate is enforced in ante-handler

	genParams.DistributionParams = distributiontypes.DefaultParams()
	genParams.DistributionParams.BaseProposerReward = sdk.MustNewDecFromStr("0.02")
	genParams.DistributionParams.BonusProposerReward = sdk.MustNewDecFromStr("0.05")
	genParams.DistributionParams.CommunityTax = sdk.MustNewDecFromStr("0.02")
	genParams.DistributionParams.WithdrawAddrEnabled = true

	genParams.GovParams = govtypes.DefaultParams()
	genParams.GovParams.DepositParams.MaxDepositPeriod = time.Hour * 24 * 3 // 3 days
	genParams.GovParams.DepositParams.MinDeposit = sdk.NewCoins(sdk.NewCoin(
		genParams.NativeCoinMetadatas[0].Base,
		sdk.NewInt(2_000_000_000_000),
	))
	genParams.GovParams.TallyParams.Quorum = sdk.MustNewDecFromStr("0.33") // 33%
	genParams.GovParams.VotingParams.VotingPeriod = time.Hour * 24 * 5    // 5 days

	genParams.CrisisConstantFee = sdk.NewCoin(
		genParams.NativeCoinMetadatas[0].Base,
		sdk.NewInt(1_000_000_000),
	)

	genParams.SlashingParams = slashingtypes.DefaultParams()
	genParams.SlashingParams.SignedBlocksWindow = int64(25000)                       // ~41 hr at 6 second blocks
	genParams.SlashingParams.MinSignedPerWindow = sdk.MustNewDecFromStr("0.05")      // 5% minimum liveness
	genParams.SlashingParams.DowntimeJailDuration = time.Minute * 10                      // 10 minutes jail period
	genParams.SlashingParams.SlashFractionDoubleSign = sdk.MustNewDecFromStr("0.05") // 5% double sign slashing
	genParams.SlashingParams.SlashFractionDowntime = sdk.MustNewDecFromStr("0.0001") // 0.01% liveness slashing

	genParams.ClaimParams = catdroptypes.Params{
		AirdropEnabled: 	false,
		AirdropStartTime:   genParams.GenesisTime.Add(time.Hour * 24 * 3), 		// 3 days after genesis
		DurationUntilDecay: time.Hour * 24 * 30,                            // 30 days
		DurationOfDecay:    time.Hour * 24 * 90,                            // 90 days
		ClaimDenom:         genParams.NativeCoinMetadatas[0].Base,
	}

	genParams.ConsensusParams = tmtypes.DefaultConsensusParams()
	genParams.ConsensusParams.Block.MaxBytes = 25 * 1024 * 1024 // 26,214,400 for cosmwasm
	genParams.ConsensusParams.Block.MaxGas = 500_000_000
	genParams.ConsensusParams.Evidence.MaxAgeDuration = genParams.StakingParams.UnbondingTime
	genParams.ConsensusParams.Evidence.MaxAgeNumBlocks = int64(genParams.StakingParams.UnbondingTime.Seconds()) / 3
	genParams.ConsensusParams.Version.AppVersion = 1

	genParams.WasmParams = wasmtypes.DefaultParams()
	genParams.WasmParams.CodeUploadAccess = wasmtypes.AllowNobody
	genParams.WasmParams.InstantiateDefaultPermission = wasmtypes.AccessTypeNobody

	return genParams
}

// params only
func TestnetGenesisParams() GenesisParams {
	genParams := MainnetGenesisParams()

	genParams.GenesisTime = time.Date(2022, 11, 24, 0, 0, 0, 0, time.UTC) // Aug 22

	genParams.ClaimParams = catdroptypes.Params{
		AirdropEnabled: 	false,
		AirdropStartTime:   genParams.GenesisTime.Add(time.Hour * 24), 		// 1 day after genesis
		DurationUntilDecay: time.Hour * 24 * 30,                            // 30 days
		DurationOfDecay:    time.Hour * 24 * 90,                            // 90 days
		ClaimDenom:         genParams.NativeCoinMetadatas[0].Base,
	}

	// gov
	genParams.GovParams.DepositParams.MaxDepositPeriod = time.Hour * 1 // 1 hour
	genParams.GovParams.TallyParams.Quorum = sdk.MustNewDecFromStr("0.1") // 10%
	genParams.GovParams.VotingParams.VotingPeriod = time.Minute * 15    // 15 min
	genParams.GovParams.DepositParams.MinDeposit = sdk.NewCoins(sdk.NewCoin(
		genParams.NativeCoinMetadatas[0].Base,
		sdk.NewInt(1_000_000),
	))

	// alloc
	genParams.AllocParams.DistributionProportions = alloctypes.DistributionProportions{
		CommunityPool: sdk.NewDecWithPrec(10, 2), // 10%
	}

	return genParams
}

func DevnetGenesisParams() GenesisParams {
	genParams := MainnetGenesisParams()

	genParams.GenesisTime = time.Now()
		genParams.ClaimParams = catdroptypes.Params{
		AirdropEnabled: 	false,
		AirdropStartTime:   genParams.GenesisTime, 		// genesis time
		DurationUntilDecay: time.Hour * 24 * 30,                            // 30 days
		DurationOfDecay:    time.Hour * 24 * 90,                            // 90 days
		ClaimDenom:         genParams.NativeCoinMetadatas[0].Base,
	}

	genParams.GovParams.DepositParams.MaxDepositPeriod = time.Hour * 1 // 1 hour
	genParams.GovParams.TallyParams.Quorum = sdk.MustNewDecFromStr("0.1")    // 10%
	genParams.GovParams.TallyParams.Threshold = sdk.MustNewDecFromStr("0.5") // 50%
	genParams.GovParams.VotingParams.VotingPeriod = time.Minute * 5          // 5 min

	return genParams
}

func LocalnetGenesisParams() GenesisParams {
	genParams := MainnetGenesisParams()

	genParams.GenesisTime = time.Date(2022, 06, 1, 0, 0, 0, 0, time.UTC) // May 01, 2022 - 00:00 UTC
	genParams.ClaimParams.AirdropStartTime = genParams.GenesisTime

	genParams.GovParams.DepositParams.MaxDepositPeriod = time.Hour * 1 // 1 hour
	genParams.GovParams.DepositParams.MinDeposit = sdk.NewCoins(sdk.NewCoin(
		genParams.NativeCoinMetadatas[0].Base,
		sdk.NewInt(1),
	))
	genParams.GovParams.TallyParams.Quorum = sdk.MustNewDecFromStr("0.1")    // 10%
	genParams.GovParams.TallyParams.Threshold = sdk.MustNewDecFromStr("0.5") // 50%
	genParams.GovParams.VotingParams.VotingPeriod = time.Minute * 1          // 5 min
	return genParams
}
