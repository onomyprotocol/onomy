<!--
order: 2
-->

# State Transitions

This document describes the state transition operations:

1. [Delegation](02_state_transitions.md#delegation)
2. [WithdrawReward](02_state_transitions.md#withdrawreward)
3. [Voting](02_state_transitions.md#voting)
4. [Proposals](02_state_transitions.md#proposals)

## Delegation

### Delegate

The DAO delegates its staking coins (NOMs) from the treasury to validators proportionately based on the self-bonded amount
of each validator with respect to the total self-bonded amount of all validators, in addition, if the validator commission rate
more than "staking_max_commission_rate" param the validator will be excluded from the DAO staking. The DAO
limits the delegation by the "staking_token_pool_rate" parameter. It means that "staking_token_pool_rate" of the total
amount of DAO's staking coins (NOMs) won't be staked.

### Rebalance

The rebalance is a process of the unbonding and delegation of the treasury staking coins (NOMs) based on
the [delegation](#Delegate) rules.

## WithdrawReward

Since the DAO is a delegator it can withdraw the delegator reward. Once the delegation reward is withdrawn it goes
to the treasury. The withdrawal reward is called based on the "withdraw_reward_period" parameter, specified in blocks.

## Voting

### Vote

Since the DAO is a delegator it can vote for the proposals otherwise the validator inherits its voting power. In
order to exclude the DAO's voting power, that DAO votes "Abstain" for all created proposals.

## Proposals

The DAO module extends the existing proposals with new types of proposals:

- `FundTreasuryProposal` - the user provides the amount to fund. If accepted the coins will be sent from the user's
  account to the treasury.
- `ExchangeWithTreasuryProposal` - the user provides the list of pairs to exchange with the treasury. If accepted the
  coins will be sent from the treasury to the user's account.
- `FundAccountProposal` - the user provides the amount to fund and recipient. If accepted the coins will be
  sent from the treasury to the recipient's account.

Additional proposals rule:

* Each funding or exchange proposal may not exceed more than the "staking_token_max_proposal_rate" of the total amount of
  DAO's staking coins (NOMs).
