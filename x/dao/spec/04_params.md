<!--
order: 4
-->

# Parameters

The DAO module contains the following parameters:

| Key              | Type   | Example                                                                                                      |
|------------------|--------|--------------------------------------------------------------------------------------------------------------|
| params           | object | {"pool_rate":"0.05", "max_proposal_rate":"0.05", "max_val_commission":"0.1", "withdraw_reward_period":51840} |
| treasury_balance | array  | [{"denom": "anom", "amount": "0"}]                                                                           |

__NOTE__: The "params" might be changed using `ParameterChangeProposal`. 
