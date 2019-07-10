# Cosmos Tutorial

## Application Goals

The goal of the application you are building is to let users buy names and to set a value these names resolve to. The owner of a given name will be the current highest bidder. ***In this section, you will learn how these simple requirements translate to application design.***

The [Cosmos SDK](https://github.com/cosmos/cosmos-sdk/)[ ](https://github.com/cosmos/cosmos-sdk/) is designed to help you build state machines. The SDK is a **modular framework**, meaning applications are built by aggregating(합계, 총액) a collection of interoperable modules. Each module contains its own message/transaction processor, while ***the SDK is responsible for routing each message to its respective module.***

- `auth`: This module defines accounts and fees and gives access to these functionalities to the rest of your application.
- `bank`: This module enables the application to create and manage tokens and token balances.
- `staking` : This module enables the application to have validators that people can delegate to.
- `distribution` : This module give a functional way to passively distribute rewards between validators and delegators.
- `slashing` : This module disincentivizes people with value staked in the network, ie. Validators.
- `nameservice`: This module does not exist yet! It will handle the core logic for the `nameservice` application you are building. It is the main piece of software you have to work on to build your application.

### State

The state represents your application at a given moment.

- It tells how much token each account possesses
- what are the owners and price of  each name
- and to what value each name resolves to.

In the SDK, everything is stored in one store called the `multistore`. Any number of key/value stores (called `KVStores` in the Cosmos SDK) can be created in this multistore. For this application, we will use one store to map `name`s to its respective `whois`, a struct that holds a name's value, owner, and price.

### Message

Messages are contained in transactions. They trigger state transitions. Each module defines a list of messages and how to handle them. Here are the messages you need to implement the desired functionality for your nameservice application:

- `MsgSetName`: This message allows name owners to set a value for a given name.
- `MsgBuyName` : This message allows accounts to buy a name and become its owner. 
  - When someone buys a name, they are required to pay the previous owner of the name a price higher than the price the previous owner paid for it. If a name does not have a previous owner yet, they must burn a `MinPrice` amount.

When a transaction (included in a block) reaches a Tendermint node, it is passed to the application via the [ABCI](https://github.com/tendermint/tendermint/tree/master/abci) - Application Blockchain Interface and decoded to get the message.

The message is then routed to the  appropriate module and handled there according to the logic defined in the `Handler`. If the state needs to be updated, the `Handler` calls the `Keeper` to perform the update. You will learn more about these concepts in the next steps of this tutorial.