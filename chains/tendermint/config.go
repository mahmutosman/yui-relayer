package tendermint

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hyperledger-labs/yui-relayer/core"
)

var _ core.ChainConfig = (*ChainConfig)(nil)

func (c ChainConfig) Build() (core.Chain, error) {
	return &Chain{
		config: c,
	}, nil
}

func (c ChainConfig) Validate() error {
	isEmpty := func(s string) bool {
		return strings.TrimSpace(s) == ""
	}

	var errs []error
	if isEmpty(c.Key) {
		errs = append(errs, fmt.Errorf("config attribute \"key\" is empty"))
	}
	if isEmpty(c.ChainId) {
		errs = append(errs, fmt.Errorf("config attribute \"chain_id\" is empty"))
	}
	if isEmpty(c.RpcAddr) {
		errs = append(errs, fmt.Errorf("config attribute \"rpc_addr\" is empty"))
	}
	if isEmpty(c.AccountPrefix) {
		errs = append(errs, fmt.Errorf("config attribute \"account_prefix\" is empty"))
	}
	if c.GasAdjustment <= 0 {
		errs = append(errs, fmt.Errorf("config attribute \"gas_adjustment\" is too small: %v", c.GasAdjustment))
	}
	if isEmpty(c.GasPrices) {
		errs = append(errs, fmt.Errorf("config attribute \"gas_prices\" is empty"))
	}
	if c.AverageBlockTimeMsec == 0 {
		errs = append(errs, fmt.Errorf("config attribute \"average_block_time_msec\" is zero"))
	}
	if c.MaxRetryForCommit == 0 {
		errs = append(errs, fmt.Errorf("config attribute \"max_retry_for_commit\" is zero"))
	}

	// errors.Join returns nil if len(errs) == 0
	return errors.Join(errs...)
}

var _ core.ProverConfig = (*ProverConfig)(nil)

func (c ProverConfig) Build(chain core.Chain) (core.Prover, error) {
	chain_, ok := chain.(*Chain)
	if !ok {
		return nil, fmt.Errorf("chain type must be %T, not %T", &Chain{}, chain)
	}
	return NewProver(chain_, c), nil
}

func (c ProverConfig) Validate() error {
	if _, err := time.ParseDuration(c.TrustingPeriod); err != nil {
		return fmt.Errorf("config attribute \"trusting_period\" is invalid: %v", err)
	}
	if c.RefreshThresholdRate <= 0 {
		return fmt.Errorf("config attribute \"refresh_threshold_rate\" is too small: %v", c.RefreshThresholdRate)
	}
	return nil
}

func (c ProverConfig) GetTrustingPeriod() time.Duration {
	if d, err := time.ParseDuration(c.TrustingPeriod); err != nil {
		panic(err)
	} else {
		return d
	}
}
