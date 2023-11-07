package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys.
var (
	KeyMintDenom                            = []byte("MintDenom")
	KeyGenesisBlockProvisions               = []byte("GenesisBlockProvisions")
	KeyReductionPeriodInBlocks              = []byte("ReductionPeriodInBlocks")
	KeyReductionFactor                      = []byte("ReductionFactor")
	KeyPoolAllocationRatio                  = []byte("PoolAllocationRatio")
	KeyDeveloperRewardsReceiver             = []byte("DeveloperRewardsReceiver")
	KeyMintingRewardsDistributionStartBlock = []byte("MintingRewardsDistributionStartBlock")
	KeyUsageIncentiveAddress                = []byte("UsageIncentiveAddress")
	KeyGrantsProgramAddress                 = []byte("GrantsProgramAddress")
	KeyTeamReserveAddress                   = []byte("TeamReserveAddress")
)

// ParamTable for minting module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams returns new mint module parameters initialized to the given values.
func NewParams(
	mintDenom string, genesisBlockProvisions sdk.Dec,
	ReductionFactor sdk.Dec, reductionPeriodInBlocks int64, distrProportions DistributionProportions,
	weightedDevRewardsReceivers []MonthlyVestingAddress, MintingRewardsDistributionStartBlock int64,
) Params {
	return Params{
		MintDenom:                            mintDenom,
		GenesisBlockProvisions:               genesisBlockProvisions,
		ReductionPeriodInBlocks:              reductionPeriodInBlocks,
		ReductionFactor:                      ReductionFactor,
		DistributionProportions:              distrProportions,
		WeightedDeveloperRewardsReceivers:    weightedDevRewardsReceivers,
		MintingRewardsDistributionStartBlock: MintingRewardsDistributionStartBlock,
	}
}

func addressTable() map[string]string {
	addressJSON := `{
		"adress1": "furya1cdz48lfj6u97cffqkp0ny4m3ht9rc6qcwgxjmh",
		"adress2": "furya1t68acn7j24yk0gylk2uyhkyflyzjsz3q7mttel",
		"adress3": "furya182zyyjuywpmqjdk6yq3mmnew36t5vgfa360kcm",
		"adress4": "furya1w2ykdkld554ak3rqdgznc5gux4vfd353nsfqxj",
		"adress5": "furya13hq8uw2fn4zhzs28340el959mjcyv7edklylgq",
		"adress6": "furya13fxzw5eeuqqvpjxycqvecgz27ducenenmdqtk7",
		"adress7": "furya1u5eve8nuzgu4x7dy7943q780dttzz33vuhws6k",
		"adress8": "furya1tud38wp5azhxamxnce9l64p3fnq398fzhdezr8",
		"adress9": "furya1a97a63vemnsqg24lamau67snnk93ap2wlm8nvs",
		"adress10": "furya1wngr6vhk9dv20c7jxt6jhrg0f47a4frwf6hcvv",
		"adress11": "furya1ppvvdmff2sfpq3kgmweuz6yvpw6r666rmr63fh",
		"adress12": "furya16a7uwstzdqgdeuanlcvmq8qcngauh7qvhy7ljs",
		"adress13": "furya1k3w35554ql9ntl5dusupkhf0s2g65t43w0rr6f",
		"adress14": "furya16xl8z7z8vxhcpwt57k4vkefuhq7pccvcgpc5vf",
		"adress15": "furya1p6hwsyqj3234r5e6efm53ywc87hjtpf85hn506",
		"adress16": "furya19gt7zdc3k2rj3xz4k0p0l4972x7swmcd9szx65",
		"adress17": "furya17zy7v8zyagctd06q46h95akssg7aaecdvmvrwd",
		"adress18": "furya1utvz4v08ullqlnpstaa6funmk4x5j6euvxk4e4",
		"adress19": "furya1uvegqucz0aqjnq7rftw2w4rwm5pt4puq2q5ht0",
		"adress20": "furya1ljz8etk5dklrq7nptr2r0a3lj56et9ru6q4u3f",
		"adress21": "furya1cdz48lfj6u97cffqkp0ny4m3ht9rc6qcwgxjmh",
		"adress22": "furya1t68acn7j24yk0gylk2uyhkyflyzjsz3q7mttel",
		"adress23": "furya182zyyjuywpmqjdk6yq3mmnew36t5vgfa360kcm",
		"adress24": "furya1w2ykdkld554ak3rqdgznc5gux4vfd353nsfqxj",
		"adress25": "furya13hq8uw2fn4zhzs28340el959mjcyv7edklylgq",
		"adress26": "furya13fxzw5eeuqqvpjxycqvecgz27ducenenmdqtk7",
		"adress27": "furya1u5eve8nuzgu4x7dy7943q780dttzz33vuhws6k",
		"adress28": "furya1wngr6vhk9dv20c7jxt6jhrg0f47a4frwf6hcvv"
	}`

	var addressMap map[string]string
	json.Unmarshal([]byte(addressJSON), &addressMap)
	return addressMap
}

func parseMonthlyVesting() []MonthlyVestingAddress {
	records := [][]string{}
	lines := strings.Split(vestingStr, "\n")
	for _, line := range lines {
		records = append(records, strings.Split(line, ","))
	}

	addressMap := addressTable()
	vAddrs := []MonthlyVestingAddress{}
	for _, addr := range records[0] {
		vAddrs = append(vAddrs, MonthlyVestingAddress{
			Address:        addressMap[addr],
			MonthlyAmounts: []sdk.Int{},
		})
	}

	for _, line := range records[1:] {
		for index, amountStr := range line {
			amountDec := sdk.MustNewDecFromStr(amountStr)
			amountInt := amountDec.Mul(sdk.NewDec(1000_000)).TruncateInt()
			vAddrs[index].MonthlyAmounts = append(vAddrs[index].MonthlyAmounts, amountInt)
		}
	}

	return vAddrs
}

// DefaultParams returns the default minting module parameters.
func DefaultParams() Params {
	return Params{
		MintDenom:               sdk.DefaultBondDenom,
		GenesisBlockProvisions:  sdk.NewDec(47000000),        //  300 million /  6307200 * 10 ^ 6
		ReductionPeriodInBlocks: 6307200,                     // 1 year - 86400 x 365 / 5
		ReductionFactor:         sdk.NewDecWithPrec(6666, 4), // 0.6666
		DistributionProportions: DistributionProportions{
			GrantsProgram:    sdk.NewDecWithPrec(10, 2), // 10%
			CommunityPool:    sdk.NewDecWithPrec(10, 2), // 10%
			UsageIncentive:   sdk.NewDecWithPrec(25, 2), // 25%
			Staking:          sdk.NewDecWithPrec(40, 2), // 40%
			DeveloperRewards: sdk.NewDecWithPrec(15, 2), // 15%
		},
		WeightedDeveloperRewardsReceivers:    parseMonthlyVesting(),
		UsageIncentiveAddress:                "furya1wngr6vhk9dv20c7jxt6jhrg0f47a4frwf6hcvv",
		GrantsProgramAddress:                 "furya1u5eve8nuzgu4x7dy7943q780dttzz33vuhws6k",
		TeamReserveAddress:                   "furya13fxzw5eeuqqvpjxycqvecgz27ducenenmdqtk7",
		MintingRewardsDistributionStartBlock: 0,
	}
}

// Validate validates mint module parameters. Returns nil if valid,
// error otherwise
func (p Params) Validate() error {
	if err := validateMintDenom(p.MintDenom); err != nil {
		return err
	}
	if err := validateGenesisBlockProvisions(p.GenesisBlockProvisions); err != nil {
		return err
	}
	if err := validateReductionPeriodInBlocks(p.ReductionPeriodInBlocks); err != nil {
		return err
	}
	if err := validateReductionFactor(p.ReductionFactor); err != nil {
		return err
	}
	if err := validateDistributionProportions(p.DistributionProportions); err != nil {
		return err
	}

	if err := validateAddress(p.UsageIncentiveAddress); err != nil {
		return err
	}

	if err := validateAddress(p.GrantsProgramAddress); err != nil {
		return err
	}

	if err := validateAddress(p.TeamReserveAddress); err != nil {
		return err
	}

	if err := validateWeightedDeveloperRewardsReceivers(p.WeightedDeveloperRewardsReceivers); err != nil {
		return err
	}
	if err := validateMintingRewardsDistributionStartBlock(p.MintingRewardsDistributionStartBlock); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Implements params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {

	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMintDenom, &p.MintDenom, validateMintDenom),
		paramtypes.NewParamSetPair(KeyGenesisBlockProvisions, &p.GenesisBlockProvisions, validateGenesisBlockProvisions),
		paramtypes.NewParamSetPair(KeyReductionPeriodInBlocks, &p.ReductionPeriodInBlocks, validateReductionPeriodInBlocks),
		paramtypes.NewParamSetPair(KeyReductionFactor, &p.ReductionFactor, validateReductionFactor),
		paramtypes.NewParamSetPair(KeyPoolAllocationRatio, &p.DistributionProportions, validateDistributionProportions),
		paramtypes.NewParamSetPair(KeyDeveloperRewardsReceiver, &p.WeightedDeveloperRewardsReceivers, validateWeightedDeveloperRewardsReceivers),
		paramtypes.NewParamSetPair(KeyUsageIncentiveAddress, &p.UsageIncentiveAddress, validateAddress),
		paramtypes.NewParamSetPair(KeyGrantsProgramAddress, &p.GrantsProgramAddress, validateAddress),
		paramtypes.NewParamSetPair(KeyTeamReserveAddress, &p.TeamReserveAddress, validateAddress),
		paramtypes.NewParamSetPair(KeyMintingRewardsDistributionStartBlock, &p.MintingRewardsDistributionStartBlock, validateMintingRewardsDistributionStartBlock),
	}
}

func validateMintDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("mint denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func validateGenesisBlockProvisions(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.LT(sdk.ZeroDec()) {
		return fmt.Errorf("genesis block provision must be non-negative")
	}

	return nil
}

func validateReductionPeriodInBlocks(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("reduction period must be positive: %d", v)
	}

	return nil
}

func validateReductionFactor(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.GT(sdk.NewDec(1)) {
		return fmt.Errorf("reduction factor cannot be greater than 1")
	}

	if v.IsNegative() {
		return fmt.Errorf("reduction factor cannot be negative")
	}

	return nil
}

func validateDistributionProportions(i interface{}) error {
	v, ok := i.(DistributionProportions)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.GrantsProgram.IsNegative() {
		return errors.New("staking distribution ratio should not be negative")
	}

	if v.CommunityPool.IsNegative() {
		return errors.New("staking distribution ratio should not be negative")
	}

	if v.UsageIncentive.IsNegative() {
		return errors.New("community pool distribution ratio should not be negative")
	}

	if v.Staking.IsNegative() {
		return errors.New("staking distribution ratio should not be negative")
	}

	if v.DeveloperRewards.IsNegative() {
		return errors.New("developer rewards distribution ratio should not be negative")
	}

	totalProportions := v.GrantsProgram.Add(v.CommunityPool).Add(v.UsageIncentive).Add(v.Staking).Add(v.DeveloperRewards)

	if !totalProportions.Equal(sdk.NewDec(1)) {
		return errors.New("total distributions ratio should be 1")
	}

	return nil
}

func validateWeightedDeveloperRewardsReceivers(i interface{}) error {
	v, ok := i.([]MonthlyVestingAddress)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// fund community pool when rewards address is empty
	if len(v) == 0 {
		return nil
	}

	return nil
}

func validateMintingRewardsDistributionStartBlock(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0 {
		return fmt.Errorf("start block must be non-negative")
	}

	return nil
}

func validateAddress(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	_, err := sdk.AccAddressFromBech32(v)

	return err
}
