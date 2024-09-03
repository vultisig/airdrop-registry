# Vultisig Airdrop Registry

This repository contains the backend code for the Vultisig Airdrop Registry. This application allows users to securely register their vaults, track balances across multiple blockchain networks, and accumulate points for a future $VULT token airdrop. 

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Airdrop Process](#airdrop-process)
- [Endpoints](#endpoints)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Overview

Vultisig Airdrop Registry is designed to reward early adopters of the Vultisig wallet security standard. Users can register their vaults, track balances, and accumulate points that determine their share of the $VULT airdrop. The initial airdrop consists of 10,000,000 $VULT, distributed based on users' vault balance over time.

## Features

- **Vault Registration**: Upload a QR code containing vault public info (public keys (ECDSA, EDdsa), hex-chain-code, vault name and etc) to register your vault.
- **Balance Tracking**: Automatically track the balance of registered vaults across multiple chains using public information.
- **Coin Management**: Add or remove coins from your vault for tracking purposes.
- **Leaderboard**: View top vaults by points on the leaderboard.
- **Airdrop Points**: Accumulate points based on the value and duration of assets in your vault, which will determine your share of the $VULT airdrop.

## Airdrop Process

1. **Register Vault**: Users send their exported vault public keys to the airdrop registry.
2. **Balance Scan**: The registry scans for funds on supported chains and begins counting the vault's airdrop value.
3. **Ongoing Scans**: Regular scans are performed each cycle to update the vault’s balance and accumulated points.
4. **Points Accumulation**: Points are calculated based on the total vault value over time.
5. **Airdrop Distribution**: After 12 months, the accumulated points determine the user's share of the $VULT airdrop.

## Endpoints

### Health Check
- **GET** `/api/ping`: Check the health of the Vultisig Airdrop Registry service.

### Public Key Derivation
- **POST** `/api/derive-public-key`: Derive public keys from the vault information.

### Vault Management
- **POST** `/api/vault`: Register a new vault.
- **DELETE** `/api/vault`: Delete a registered vault.
- **GET** `/api/vault/:ecdsaPublicKey/:eddsaPublicKey`: Get details of a specific vault.
- **POST** `/api/vault/:ecdsaPublicKey/:eddsaPublicKey/alias`: Update the alias of a vault.
- **GET** `/api/vault/shared/:uid`: Get vault information by UID.
- **POST** `/api/vault/join-airdrop`: Register a vault for the airdrop.
- **POST** `/api/vault/exit-airdrop`: Unregister a vault from the airdrop.

### Coin Management
- **DELETE** `/api/coin/:ecdsaPublicKey/:eddsaPublicKey/:coinID`: Remove a coin from a vault.
- **POST** `/api/coin/:ecdsaPublicKey/:eddsaPublicKey`: Add a coin to a vault.
- **GET** `/api/coin/:ecdsaPublicKey/:eddsaPublicKey`: Get all coins for a vault.

## Usage
- **Register for Airdrop**: 
  - Use the `/api/vault/join-airdrop` endpoint to register your vault for the airdrop. This will start the process of tracking your vault's balance and accumulating points.
- **Track Balances**:
  - Users can see and track the balance of their vault across different blockchain networks.
  - Additionally, users can view their Liquidity Provider (LP) positions, Saver accounts, and Bond information in both Thorchain and MayaChain. This allows for comprehensive tracking of all assets held in the vault across these platforms.
- **Check Points**:
  - Use the `/api/vault/:ecdsaPublicKey/:eddsaPublicKey` endpoint to retrieve the details of your vault, including the latest airdrop points accumulated. This will help you monitor your progress and estimate your share of the upcoming $VULT airdrop.

- **Proof of Reserve**:
  - Users can share their vault with others as proof of reserve or for other purposes. This feature is useful for demonstrating the assets held within a vault without compromising security or exposing sensitive information. To share your vault, use the `/api/vault/shared/:uid` endpoint to generate a shareable link or details.


## Contributing
Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.


## License
Vultisig is licensed under the STMF License (Set the Memes Free).
For full terms and conditions, please refer to the [Vultisig License Documentation](https://docs.vultisig.com/other/licence).