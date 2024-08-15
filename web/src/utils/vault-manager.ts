import { initWasm, WalletCore } from "@trustwallet/wallet-core";
import type { CoinType } from "@trustwallet/wallet-core/dist/src/wallet-core";

import api from "utils/api";

namespace VaultManager {
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

  export interface CoinMeta {
    address?: string;
    icon: string;
    chainID: Chain;
    name: string;
    ticker: string;
    isNative: boolean;
  }

  export enum Chain {
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

  export const coins: CoinMeta[] = [
    {
      icon: "/chains/arbitrum.svg",
      chainID: Chain.ARBITRUM,
      name: "Arbitrum",
      ticker: "ARB",
      isNative: true,
    },
    {
      icon: "/chains/avax.svg",
      chainID: Chain.AVALANCHE,
      name: "Avalanche",
      ticker: "AVAX",
      isNative: true,
    },
    {
      icon: "/chains/eth.svg",
      chainID: Chain.BASE,
      name: "ethereum",
      ticker: "ETH",
      isNative: true,
    },
    {
      icon: "/chains/btc.svg",
      chainID: Chain.BITCOIN,
      name: "Bitcoin",
      ticker: "BTC",
      isNative: true,
    },
    {
      icon: "/chains/bch.svg",
      chainID: Chain.BITCOIN_CASH,
      name: "BitcoinCash",
      ticker: "BCH",
      isNative: true,
    },
    {
      icon: "/chains/eth.svg",
      chainID: Chain.BLAST,
      name: "Blast",
      ticker: "ETH",
      isNative: true,
    },
    {
      icon: "/chains/bnb.svg",
      chainID: Chain.BSCCHAIN,
      name: "BinanceCoin",
      ticker: "BNB",
      isNative: true,
    },
    {
      icon: "/chains/cro.svg",
      chainID: Chain.CRONOSCHAIN,
      name: "cronosChain",
      ticker: "CRO",
      isNative: true,
    },
    {
      icon: "/chains/dash.svg",
      chainID: Chain.DASH,
      name: "Dash",
      ticker: "DASH",
      isNative: true,
    },
    {
      icon: "/chains/doge.svg",
      chainID: Chain.DOGECOIN,
      name: "Dogecoin",
      ticker: "DOGE",
      isNative: true,
    },
    {
      icon: "/chains/dydx.svg",
      chainID: Chain.DYDX,
      name: "dydxChain",
      ticker: "DYDX",
      isNative: true,
    },
    {
      icon: "/chains/eth.svg",
      chainID: Chain.ETHEREUM,
      name: "Ethereum",
      ticker: "ETH",
      isNative: true,
    },
    {
      icon: "/chains/atom.svg",
      chainID: Chain.GAIACHAIN,
      name: "gaiaChain",
      ticker: "ATOM",
      isNative: true,
    },
    {
      icon: "/chains/kuji.svg",
      chainID: Chain.KUJIRA,
      name: "Kujira",
      ticker: "KUJI",
      isNative: true,
    },
    {
      icon: "/chains/ltc.svg",
      chainID: Chain.LITECOIN,
      name: "Litecoin",
      ticker: "LTC",
      isNative: true,
    },
    {
      icon: "/chains/cacao.svg",
      chainID: Chain.MAYACHAIN,
      name: "mayaChain",
      ticker: "CACAO",
      isNative: true,
    },
    {
      icon: "/chains/eth.svg",
      chainID: Chain.OPTIMISM,
      name: "Optimism",
      ticker: "ETH",
      isNative: true,
    },
    {
      icon: "/chains/dot.svg",
      chainID: Chain.POLKADOT,
      name: "Polkadot",
      ticker: "DOT",
      isNative: true,
    },
    {
      icon: "/chains/matic.svg",
      chainID: Chain.POLYGON,
      name: "Polygon",
      ticker: "MATIC",
      isNative: true,
    },
    {
      icon: "/chains/sol.svg",
      chainID: Chain.SOLANA,
      name: "Solana",
      ticker: "SOL",
      isNative: true,
    },
    {
      icon: "/chains/sui.svg",
      chainID: Chain.SUI,
      name: "Sui",
      ticker: "SUI",
      isNative: true,
    },
    {
      icon: "/chains/rune.svg",
      chainID: Chain.THORCHAIN,
      name: "THORChain",
      ticker: "RUNE",
      isNative: true,
    },
    {
      icon: "/chains/zksync.svg",
      chainID: Chain.ZKSYNC,
      name: "zkSync",
      ticker: "ZK",
      isNative: true,
    },
  ];

  let core: WalletCore;
  let vault: Vault;
  let coin: CoinType[] = [];

  const setCoin = (): void => {
    coin[Chain.ARBITRUM] = core.CoinType.arbitrum;
    coin[Chain.AVALANCHE] = core.CoinType.avalancheCChain;
    coin[Chain.BASE] = core.CoinType.base;
    coin[Chain.BITCOIN] = core.CoinType.bitcoin;
    coin[Chain.BITCOIN_CASH] = core.CoinType.bitcoinCash;
    coin[Chain.BLAST] = core.CoinType.blast;
    coin[Chain.BSCCHAIN] = core.CoinType.binance;
    coin[Chain.CRONOSCHAIN] = core.CoinType.cronosChain;
    coin[Chain.DASH] = core.CoinType.dash;
    coin[Chain.DOGECOIN] = core.CoinType.dogecoin;
    coin[Chain.DYDX] = core.CoinType.dydx;
    coin[Chain.ETHEREUM] = core.CoinType.ethereum;
    coin[Chain.GAIACHAIN] = core.CoinType.cosmos;
    coin[Chain.KUJIRA] = core.CoinType.kujira;
    coin[Chain.LITECOIN] = core.CoinType.litecoin;
    coin[Chain.MAYACHAIN] = core.CoinType.thorchain;
    coin[Chain.OPTIMISM] = core.CoinType.optimism;
    coin[Chain.POLKADOT] = core.CoinType.polkadot;
    coin[Chain.POLYGON] = core.CoinType.polygon;
    coin[Chain.SOLANA] = core.CoinType.solana;
    coin[Chain.SUI] = core.CoinType.sui;
    coin[Chain.THORCHAIN] = core.CoinType.thorchain;
    coin[Chain.ZKSYNC] = core.CoinType.zksync;
  };

  const getVault = (): void => {
    const _vault = localStorage.getItem("vault");

    if (_vault) vault = JSON.parse(_vault);
  };

  const setVault = (vault: Vault): void => {
    localStorage.setItem("vault", JSON.stringify(vault));
  };

  const getECDSAAddress = (value: Chain, prefix?: string): Promise<string> => {
    return new Promise((resolve, reject) => {
      api
        .derivePublicKey({
          publicKeyEcdsa: vault.publicKeyEcdsa,
          hexChainCode: vault.hexChainCode,
          derivePath: core.CoinTypeExt.derivationPath(coin[value]),
        })
        .then(({ data }) => {
          const bytes = core.HexCoding.decode(data.publicKey);
          let address: string;

          const publicKey = core.PublicKey.createWithData(
            bytes,
            core.PublicKeyType.secp256k1
          );

          if (prefix) {
            address = core.AnyAddress.createBech32WithPublicKey(
              publicKey,
              coin[value],
              prefix
            )?.description();
          } else {
            address = core.AnyAddress.createWithPublicKey(
              publicKey,
              coin[value]
            )?.description();
          }

          address ? resolve(address) : reject();
        })
        .catch(() => {
          reject();
        });
    });
  };

  const getEDDSAAdress = (value: Chain): Promise<string> => {
    return new Promise((resolve, reject) => {
      const bytes = core.HexCoding.decode(vault.publicKeyEddsa);

      const eddsaKey = core.PublicKey.createWithData(
        bytes,
        core.PublicKeyType.ed25519
      );

      const address = core.AnyAddress.createWithPublicKey(
        eddsaKey,
        coin[value]
      )?.description();

      address ? resolve(address) : reject();
    });
  };

  export const register = (vault: Vault): Promise<void> => {
    return new Promise((resolve, reject) => {
      api
        .register(vault)
        .then(() => {
          resolve();
        })
        .catch(() => {
          reject();
        })
        .finally(() => {
          setVault(vault);
        });
    });
  };

  export const init = (): Promise<void> => {
    return new Promise((resolve, reject) => {
      initWasm()
        .then((_core) => {
          core = _core;

          getVault();

          setCoin();

          resolve();
        })
        .catch(() => {
          reject();
        });
    });
  };

  export const getChain = (id: Chain): Promise<CoinMeta> => {
    return new Promise((resolve, reject) => {
      if (core) {
        const _chain = coins.find((chain) => chain.chainID === id);

        _chain ? resolve(_chain) : reject();
      } else {
        reject();
      }
    });
  };

  export const getChains = (ids: Chain[]): Promise<CoinMeta[]> => {
    return new Promise((resolve, reject) => {
      if (core) {
        const _chains = coins.filter(
          (chain) => ids.indexOf(chain.chainID) >= 0
        );

        resolve(_chains);
      } else {
        reject();
      }
    });
  };

  export const getAddress = (value: Chain): Promise<string> => {
    return new Promise((resolve, reject) => {
      if (core) {
        switch (value) {
          // EDDSA
          case Chain.POLKADOT:
          case Chain.SOLANA:
          case Chain.SUI: {
            getEDDSAAdress(value)
              .then((address) => {
                resolve(address);
              })
              .catch(() => {});

            break;
          }
          // ECDSA
          case Chain.MAYACHAIN: {
            getECDSAAddress(value, "maya")
              .then((address) => {
                resolve(address);
              })
              .catch(() => {});

            break;
          }
          default: {
            getECDSAAddress(value)
              .then((address) => {
                resolve(address);
              })
              .catch(() => {});

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
