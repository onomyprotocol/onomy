<!--
order: 0
title: DAO Overview
parent:
  title: "dao"
-->

# `dao`

## Abstract

This paper specifies the DAO (Decentralized Autonomous Organization) module of the onomy chain. The DAO module provides
the advanced model of funding and exchange using governance proposals. Additionally, it extends the existing
staking by the dynamic delegation of the staking tokens from the treasury to the validators.

## Contents

1. **[State](01_state.md)**
   - [Treasury](01_state.md#treasury)
   - [Params](01_state.md#params)
   - [Staking](01_state.md#staking)

2. **[State Transitions](02_state_transitions.md)**
   - [Delegation](02_state_transitions.md#delegation)
   - [WithdrawReward](02_state_transitions.md#withdrawreward)
   - [Voting](02_state_transitions.md#voting)
   - [Proposals](02_state_transitions.md#proposals)

5. **[End-Block](03_end_block.md)**

4. **[Parameters](04_params.md)**