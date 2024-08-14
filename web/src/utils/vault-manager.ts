import { initWasm, WalletCore } from "@trustwallet/wallet-core";

import api from "utils/api";

namespace VaultManager {
  export enum ChainType {
    ARBITRUM,
    AVALANCHE,
    BASE,
    BITCOIN,
    BITCOIN_CASH,
    BLAST,
    BSCCHAIN,
    CRONOSCHAIN,
    DASH,
    DOGECOIN,
    DYDX,
    ETHEREUM,
    GAIACHAIN,
    KUJIRA,
    LITECOIN,
    MAYACHAIN,
    OPTIMISM,
    POLKADOT,
    POLYGON,
    SOLANA,
    SUI,
    THORCHAIN,
    ZKSYNC,
  }

  export interface Vault {
    uid: string;
    name: string;
    publicKeyEcdsa: string;
    publicKeyEddsa: string;
    hexChainCode: string;
  }

  export interface Derivation {
    publicKeyEcdsa: string;
    hexChainCode: string;
    derivePath: string;
  }

  let core: WalletCore;
  let vault: Vault;

  const derivePath = {
    [ChainType.BITCOIN]: "m/84'/0'/0'/0/0",
    [ChainType.ETHEREUM]: "m/44'/60'/0'/0/0",
    [ChainType.THORCHAIN]: "m/44'/931'/0'/0/0",
    [ChainType.MAYACHAIN]: "m/44'/931'/0'/0/0",
    [ChainType.ARBITRUM]: "m/44'/60'/0'/0/0",
    [ChainType.AVALANCHE]: "m/44'/60'/0'/0/0",
    [ChainType.BSCCHAIN]: "m/44'/60'/0'/0/0",
    [ChainType.BASE]: "m/44'/60'/0'/0/0",
    [ChainType.BITCOIN_CASH]: "m/44'/145'/0'/0/0",
    [ChainType.BLAST]: "m/44'/60'/0'/0/0",
    [ChainType.CRONOSCHAIN]: "m/44'/60'/0'/0/0",
    [ChainType.DASH]: "m/44'/5'/0'/0/0",
    [ChainType.DOGECOIN]: "m/44'/3'/0'/0/0",
    [ChainType.DYDX]: "m/44'/118'/0'/0/0",
    [ChainType.GAIACHAIN]: "m/44'/118'/0'/0/0",
    [ChainType.KUJIRA]: "m/44'/118'/0'/0/0",
    [ChainType.LITECOIN]: "m/84'/2'/0'/0/0",
    [ChainType.OPTIMISM]: "m/44'/60'/0'/0/0",
    [ChainType.POLYGON]: "m/44'/60'/0'/0/0",
    [ChainType.ZKSYNC]: "m/44'/60'/0'/0/0",
  };

  const getVault = (): Vault => {
    const vault = localStorage.getItem("vault");

    return vault && JSON.parse(vault);
  };

  const setVault = (vault: Vault): void => {
    localStorage.setItem("vault", JSON.stringify(vault));
  };

  export const register = (_vault: Vault): Promise<void> => {
    return new Promise((resolve, reject) => {
      api
        .register(_vault)
        .then(() => {
          resolve();
        })
        .catch(() => {
          reject();
        })
        .finally(() => {
          setVault(_vault);
        });
    });
  };

  export const initiate = (): Promise<void> => {
    return new Promise((resolve, reject) => {
      initWasm()
        .then((_core) => {
          core = _core;

          vault = getVault();

          resolve();
        })
        .catch(() => {
          reject();
        });
    });
  };

  export const getAddress = (chain: ChainType): Promise<string> => {
    return new Promise((resolve, reject) => {
      if (core) {
        switch (chain) {
          // ECDSA
          case ChainType.BITCOIN: {
            api
              .derivePublicKey({
                publicKeyEcdsa: vault.publicKeyEcdsa,
                hexChainCode: vault.hexChainCode,
                derivePath: derivePath[chain],
              })
              .then(({ data }) => {
                const bytes = core.HexCoding.decode(data.publicKey);

                const publicKey = core.PublicKey.createWithData(
                  bytes,
                  core.PublicKeyType.secp256k1
                );

                const address = core.AnyAddress.createWithPublicKey(
                  publicKey,
                  core.CoinType.bitcoin
                ).description();

                resolve(address);
              })
              .catch(() => {});

            break;
          }
          case ChainType.ETHEREUM: {
            api
              .derivePublicKey({
                publicKeyEcdsa: vault.publicKeyEcdsa,
                hexChainCode: vault.hexChainCode,
                derivePath: derivePath[chain],
              })
              .then(({ data }) => {
                const bytes = core.HexCoding.decode(data.publicKey);

                const publicKey = core.PublicKey.createWithData(
                  bytes,
                  core.PublicKeyType.secp256k1
                );

                const address = core.AnyAddress.createWithPublicKey(
                  publicKey,
                  core.CoinType.ethereum
                ).description();

                resolve(address);
              })
              .catch(() => {});

            break;
          }
          case ChainType.THORCHAIN: {
            api
              .derivePublicKey({
                publicKeyEcdsa: vault.publicKeyEcdsa,
                hexChainCode: vault.hexChainCode,
                derivePath: derivePath[chain],
              })
              .then(({ data }) => {
                const bytes = core.HexCoding.decode(data.publicKey);

                const publicKey = core.PublicKey.createWithData(
                  bytes,
                  core.PublicKeyType.secp256k1
                );

                const address = core.AnyAddress.createWithPublicKey(
                  publicKey,
                  core.CoinType.thorchain
                ).description();

                resolve(address);
              })
              .catch(() => {});

            break;
          }
          // EDDSA
          case ChainType.SOLANA: {
            const bytes = core.HexCoding.decode(vault.publicKeyEddsa);

            const eddsaKey = core.PublicKey.createWithData(
              bytes,
              core.PublicKeyType.ed25519
            );

            const address = core.AnyAddress.createWithPublicKey(
              eddsaKey,
              core.CoinType.solana
            ).description();

            resolve(address);

            break;
          }
          case ChainType.SUI: {
            const bytes = core.HexCoding.decode(vault.publicKeyEddsa);

            const eddsaKey = core.PublicKey.createWithData(
              bytes,
              core.PublicKeyType.ed25519
            );

            const address = core.AnyAddress.createWithPublicKey(
              eddsaKey,
              core.CoinType.sui
            ).description();

            resolve(address);

            break;
          }
          case ChainType.POLKADOT: {
            const bytes = core.HexCoding.decode(vault.publicKeyEddsa);

            const eddsaKey = core.PublicKey.createWithData(
              bytes,
              core.PublicKeyType.ed25519
            );

            const address = core.AnyAddress.createWithPublicKey(
              eddsaKey,
              core.CoinType.polkadot
            ).description();

            resolve(address);

            break;
          }
          default: {
            reject();
            break;
          }
        }
      } else {
        reject();
      }
    });
  };
}

export default VaultManager;
