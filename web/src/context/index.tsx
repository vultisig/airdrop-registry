import {
  FC,
  ReactNode,
  useState,
  createContext,
  useEffect,
  useContext,
} from "react";
import { Modal, Spin } from "antd";
import { initWasm, WalletCore } from "@trustwallet/wallet-core";

import {
  Chain,
  Currency,
  Language,
  OneInch,
  balanceAPI,
  defTokens,
  errorKey,
  storageKey,
} from "utils/constants";
import { CoinProps, CoinRef, TokenProps, VaultProps } from "utils/interfaces";
import i18n from "i18n/config";
import api from "utils/api";

import SplashScreen from "components/splash-screen";

interface VaultContext {
  addVault: (vault: VaultProps) => Promise<void>;
  changeCurrency: (vault: VaultProps, currency: Currency) => void;
  fetchTokens: (coin: CoinProps) => Promise<void>;
  setVault: (vault: VaultProps) => void;
  useVault: (vault: VaultProps) => void;
  toggleCoin: (coin: TokenProps, vault: VaultProps) => Promise<void>;
  tokens: TokenProps[];
  currency: Currency;
  vault?: VaultProps;
  vaults: VaultProps[];
}

interface InitialState {
  tokens: TokenProps[];
  core?: WalletCore;
  currency: Currency;
  loaded: boolean;
  loading: boolean;
  vaults: VaultProps[];
  vault?: VaultProps;
  wcRefrence: CoinRef;
}

const VaultContext = createContext<VaultContext | undefined>(undefined);

const Component: FC<{ children: ReactNode }> = ({ children }) => {
  const initialState: InitialState = {
    tokens: defTokens,
    currency: Currency.USD,
    loaded: false,
    loading: false,
    vaults: [],
    wcRefrence: {},
  };
  const [state, setState] = useState(initialState);
  const { core, currency, loaded, loading, tokens, vault, vaults, wcRefrence } =
    state;

  const changeCurrency = (vault: VaultProps, currency: Currency): void => {
    if (!loading) {
      setState((prevState) => ({ ...prevState, loading: true }));

      getValue(vault.coins, currency)
        .then((coins) => {
          setState((prevState) => ({
            ...prevState,
            currency,
            loading: false,
            vault: { ...vault, coins },
          }));

          localStorage.setItem(storageKey.CURRENCY, currency);
        })
        .catch(() => {
          setState((prevState) => ({ ...prevState, loading: false }));
        });
    }
  };

  const fetchTokens = (coin: CoinProps): Promise<void> => {
    return new Promise((resolve, reject) => {
      const token = defTokens.find((item) => item.chain === coin.chain);
      let id: number;

      switch (coin.chain) {
        case Chain.ARBITRUM:
          id = OneInch.ARBITRUM;
          break;
        case Chain.AVALANCHE:
          id = OneInch.AVALANCHE;
          break;
        case Chain.BASE:
          id = OneInch.BASE;
          break;
        case Chain.BLAST:
          id = OneInch.BLAST;
          break;
        case Chain.BSCCHAIN:
          id = OneInch.BSCCHAIN;
          break;
        case Chain.CRONOSCHAIN:
          id = OneInch.CRONOSCHAIN;
          break;
        case Chain.ETHEREUM:
          id = OneInch.ETHEREUM;
          break;
        case Chain.OPTIMISM:
          id = OneInch.OPTIMISM;
          break;
        case Chain.POLYGON:
          id = OneInch.POLYGON;
          break;
        default:
          id = 0;
          break;
      }

      if (id && token) {
        api
          .oneInch(id)
          .then(({ data }) => {
            setState((prevState) => {
              const tokens: TokenProps[] = [];

              Object.entries(data.tokens).forEach(([key, value]) => {
                const notFound =
                  prevState.tokens.findIndex(
                    (token) => token.contractAddress === key
                  ) < 0;

                if (notFound) {
                  tokens.push({
                    chain: coin.chain,
                    cmcId: 0,
                    contractAddress: key,
                    decimals: value.decimals,
                    hexPublicKey: token.hexPublicKey,
                    isDefault: false,
                    isLocally: false,
                    isNative: false,
                    logo: value.logoURI,
                    ticker: value.symbol,
                  });
                }
              });

              return { ...prevState, tokens: [...prevState.tokens, ...tokens] };
            });

            resolve();
          })
          .catch(() => {
            reject();
          });
      } else {
        reject();
      }
    });
  };

  const getECDSAAddress = (
    chain: Chain,
    vault: VaultProps,
    prefix?: string
  ): Promise<string> => {
    return new Promise((resolve, reject) => {
      const coin = wcRefrence[chain];

      if (coin && core) {
        api
          .derivePublicKey({
            publicKeyEcdsa: vault.publicKeyEcdsa,
            hexChainCode: vault.hexChainCode,
            derivePath: core.CoinTypeExt.derivationPath(coin),
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
                coin,
                prefix
              )?.description();
            } else {
              address = core.AnyAddress.createWithPublicKey(
                publicKey,
                coin
              )?.description();
            }

            address ? resolve(address) : reject();
          })
          .catch(() => {
            reject();
          });
      } else {
        reject();
      }
    });
  };

  const getEDDSAAdress = (chain: Chain, vault: VaultProps): Promise<string> => {
    return new Promise((resolve, reject) => {
      const coin = wcRefrence[chain];

      if (coin && core) {
        const bytes = core.HexCoding.decode(vault.publicKeyEddsa);

        const eddsaKey = core.PublicKey.createWithData(
          bytes,
          core.PublicKeyType.ed25519
        );

        const address = core.AnyAddress.createWithPublicKey(
          eddsaKey,
          coin
        )?.description();

        address ? resolve(address) : reject();
      } else {
        reject();
      }
    });
  };

  const getAddress = (coin: TokenProps, vault: VaultProps): Promise<string> => {
    return new Promise((resolve, reject) => {
      if (coin.isNative) {
        switch (coin.chain) {
          // EDDSA
          case Chain.POLKADOT:
          case Chain.SOLANA:
          case Chain.SUI: {
            getEDDSAAdress(coin.chain, vault)
              .then((address) => {
                resolve(address);
              })
              .catch(() => {
                reject();
              });

            break;
          }
          // ECDSA
          case Chain.MAYACHAIN: {
            getECDSAAddress(coin.chain, vault, "maya")
              .then((address) => {
                resolve(address);
              })
              .catch(() => {
                reject();
              });

            break;
          }
          case Chain.BITCOINCASH: {
            getECDSAAddress(coin.chain, vault)
              .then((address) => {
                resolve(address.replace("bitcoincash:", ""));
              })
              .catch(() => {
                reject();
              });

            break;
          }
          default: {
            getECDSAAddress(coin.chain, vault)
              .then((address) => {
                resolve(address);
              })
              .catch(() => {
                reject();
              });

            break;
          }
        }
      } else {
        const address = vault.coins.find(
          ({ chain, isNativeToken }) => isNativeToken && chain === coin.chain
        )?.address;

        address ? resolve(address) : reject();
      }
    });
  };

  const addCoin = (coin: TokenProps, vault: VaultProps): Promise<CoinProps> => {
    return new Promise((resolve, reject) => {
      const newCoin: CoinProps = {
        address: "",
        balance: 0,
        chain: coin.chain,
        cmcId: coin.cmcId,
        contractAddress: coin.contractAddress,
        decimals: coin.decimals,
        hexPublicKey:
          coin.hexPublicKey === "ECDSA"
            ? vault.publicKeyEcdsa
            : vault.publicKeyEddsa,
        ID: 0,
        isNativeToken: coin.isNative,
        logo: coin.logo,
        ticker: coin.ticker,
        value: 0,
      };

      getAddress(coin, vault)
        .then((address) => {
          newCoin.address = address;

          api.coin
            .add(vault, newCoin)
            .then(({ data: { coinId } }) => {
              const coin: CoinProps = { ...newCoin, ID: coinId };

              getBalance(coin)
                .then((coin) => {
                  getValue([coin], currency)
                    .then(([coin]) => {
                      resolve(coin);
                    })
                    .catch(() => {
                      reject();
                    });
                })
                .catch(() => {
                  reject();
                });
            })
            .catch(() => {
              reject();
            });
        })
        .catch(() => {
          reject();
        });
    });
  };

  const delCoin = (coin: CoinProps, vault: VaultProps): Promise<void> => {
    return new Promise((resolve, reject) => {
      api.coin
        .del(vault, coin)
        .then(() => {
          resolve();
        })
        .catch(() => {
          reject();
        });
    });
  };

  const toggleCoin = (coin: TokenProps, vault: VaultProps): Promise<void> => {
    return new Promise((resolve, reject) => {
      const selectedCoin = vault.coins.find(
        ({ chain, ticker }) => coin.chain === chain && coin.ticker === ticker
      );

      if (selectedCoin) {
        delCoin(selectedCoin, vault)
          .then(() => {
            setState((prevState) => {
              if (prevState.vault?.uid === vault.uid) {
                return {
                  ...prevState,
                  vault: {
                    ...prevState.vault,
                    coins: prevState.vault.coins.filter(
                      ({ chain, ticker }) =>
                        coin.chain !== chain || coin.ticker !== ticker
                    ),
                  },
                };
              }

              return prevState;
            });

            resolve();
          })
          .catch(() => {
            reject();
          });
      } else {
        getCMC(coin)
          .then((cmcId) => {
            addCoin({ ...coin, cmcId }, vault)
              .then((coin) => {
                setState((prevState) => {
                  if (prevState.vault?.uid === vault.uid) {
                    return {
                      ...prevState,
                      vault: {
                        ...prevState.vault,
                        coins: [...prevState.vault.coins, coin],
                      },
                    };
                  }

                  return prevState;
                });

                resolve();
              })
              .catch(() => {
                reject();
              });
          })
          .catch(() => {
            reject();
          });
      }
    });
  };

  const getCMC = (coin: TokenProps): Promise<number> => {
    return new Promise((resolve, reject) => {
      if (coin.isLocally) {
        resolve(coin.cmcId);
      } else {
        api.coin
          .cmc(coin.contractAddress)
          .then(({ data }) => {
            const [key] = Object.keys(data.data);

            key ? resolve(parseInt(key)) : reject();
          })
          .catch(() => {
            reject();
          });
      }
    });
  };

  const getBalance = (coin: CoinProps): Promise<CoinProps> => {
    return new Promise((resolve) => {
      const uid = Math.floor(Math.random() * 10000);
      const path = balanceAPI[coin.chain];

      switch (coin.chain) {
        // Cosmos
        case Chain.DYDX:
        case Chain.GAIACHAIN:
        case Chain.KUJIRA:
        case Chain.MAYACHAIN:
        case Chain.THORCHAIN: {
          api.balance
            .cosmos(`${path}/${coin.address}`)
            .then(({ data: { balances } }) => {
              if (balances.length && balances[0].amount) {
                resolve({
                  ...coin,
                  balance:
                    parseInt(balances[0].amount) / Math.pow(10, coin.decimals),
                });
              } else {
                resolve({ ...coin, balance: 0 });
              }
            })
            .catch(() => {
              resolve({ ...coin, balance: 0 });
            });

          break;
        }
        // EVM
        case Chain.ARBITRUM:
        case Chain.AVALANCHE:
        case Chain.BASE:
        case Chain.BLAST:
        case Chain.BSCCHAIN:
        case Chain.CRONOSCHAIN:
        case Chain.ETHEREUM:
        case Chain.OPTIMISM:
        case Chain.POLYGON: {
          api.balance
            .evm(path, {
              jsonrpc: "2.0",
              method: coin.isNativeToken ? "eth_getBalance" : "eth_call",
              params: [
                coin.isNativeToken
                  ? coin.address
                  : {
                      data: `0x70a08231000000000000000000000000${coin.address.replace(
                        "0x",
                        ""
                      )}`,
                      to: coin.contractAddress,
                    },
                "latest",
              ],
              id: uid,
            })
            .then(({ data: { result } }) => {
              resolve({
                ...coin,
                balance: parseInt(result, 16) / Math.pow(10, coin.decimals),
              });
            })
            .catch(() => {
              resolve({ ...coin, balance: 0 });
            });

          break;
        }
        case Chain.POLKADOT: {
          api.balance
            .polkadot(path, { key: coin.address })
            .then(({ data: { data } }) => {
              if (data && data.account && data.account.balance) {
                const balance = data.account.balance.replace(".", "");

                resolve({
                  ...coin,
                  balance: parseInt(balance) / Math.pow(10, coin.decimals),
                });
              } else {
                resolve({ ...coin, balance: 0 });
              }
            })
            .catch(() => {
              resolve({ ...coin, balance: 0 });
            });

          break;
        }
        case Chain.SOLANA: {
          api.balance
            .solana(path, {
              jsonrpc: "2.0",
              method: coin.isNativeToken
                ? "getBalance"
                : "getTokenAccountsByOwner",
              params: coin.isNativeToken
                ? [coin.address]
                : [
                    "address",
                    { mint: coin.contractAddress },
                    { encoding: "jsonParsed" },
                  ],
              id: 1,
            })
            .then(({ data }) => {
              resolve({
                ...coin,
                balance: data.result.value / Math.pow(10, coin.decimals),
              });
            })
            .catch(() => {
              resolve({ ...coin, balance: 0 });
            });

          break;
        }
        // UTXO
        case Chain.BITCOIN:
        case Chain.BITCOINCASH:
        case Chain.DASH:
        case Chain.DOGECOIN:
        case Chain.LITECOIN: {
          api.balance
            .utxo(`${path}/${coin.address}?state=latest`)
            .then(({ data: { data } }) => {
              if (
                data &&
                data[coin.address] &&
                data[coin.address].address &&
                typeof data[coin.address].address.balance === "number"
              ) {
                resolve({
                  ...coin,
                  balance:
                    data[coin.address].address.balance /
                    Math.pow(10, coin.decimals),
                });
              } else {
                resolve({ ...coin, balance: 0 });
              }
            })
            .catch(() => {
              resolve({ ...coin, balance: 0 });
            });

          break;
        }
        default:
          resolve({ ...coin, balance: 0 });
          break;
      }
    });
  };

  const getValue = (
    coins: CoinProps[],
    currency: Currency
  ): Promise<CoinProps[]> => {
    return new Promise((resolve, reject) => {
      const ids = tokens.map(({ cmcId }) => cmcId);

      api.coin
        .values(ids, currency)
        .then(({ data }) => {
          coins.forEach((coin) => {
            if (data?.data && data?.data[coin.cmcId]?.quote) {
              coin.value = data.data[coin.cmcId].quote[currency]?.price || 0;
            } else {
              coin.value = 0;
            }
          });

          resolve(coins);
        })
        .catch(() => {
          reject();
        });
    });
  };

  const addVault = (vault: VaultProps): Promise<void> => {
    return new Promise((resolve, reject) => {
      api.vault
        .add(vault)
        .then(() => {
          getVault(vault)
            .then((vault) => {
              setState((prevState) => {
                const vaults = [
                  vault,
                  ...prevState.vaults.filter(({ uid }) => uid !== vault.uid),
                ];

                setVaults(vaults);

                return { ...prevState, vault, vaults };
              });

              resolve();
            })
            .catch((error) => {
              reject(error);
            });
        })
        .catch((error) => {
          switch (error) {
            case errorKey.VAULT_ALREADY_REGISTERED: {
              getVault(vault)
                .then((vault) => {
                  setState((prevState) => {
                    const vaults = [
                      vault,
                      ...prevState.vaults.filter(
                        ({ uid }) => uid !== vault.uid
                      ),
                    ];

                    setVaults(vaults);

                    return { ...prevState, vault, vaults };
                  });

                  resolve();
                })
                .catch((error) => {
                  reject(error);
                });

              break;
            }
            default: {
              reject(error);

              break;
            }
          }
        });
    });
  };

  const getVault = (vault: VaultProps): Promise<VaultProps> => {
    return new Promise((resolve, reject) => {
      api.vault
        .get(vault)
        .then(({ data }) => {
          if (data.coins.length) {
            getValue(data.coins, currency)
              .then((coins) => {
                const promises = coins.map((coin) => getBalance(coin));

                Promise.all(promises)
                  .then((coins) => {
                    resolve({ ...vault, ...data, coins });
                  })
                  .catch(() => {
                    reject();
                  });
              })
              .catch(() => {
                reject();
              });
          } else {
            const promises = tokens
              .filter((coin) => coin.isDefault)
              .map((coin) => addCoin(coin, { ...vault, ...data }));

            Promise.all(promises)
              .then((coins) => {
                getValue(coins, currency)
                  .then((coins) => {
                    const promises = coins.map((coin) => getBalance(coin));

                    Promise.all(promises)
                      .then((coins) => {
                        resolve({ ...vault, ...data, coins });
                      })
                      .catch(() => {
                        reject();
                      });
                  })
                  .catch(() => {
                    reject();
                  });
              })
              .catch(() => {});
          }
        })
        .catch(() => {
          reject();
        });
    });
  };

  const getVaults = (vaults: VaultProps[]): Promise<VaultProps[]> => {
    return new Promise((resolve, reject) => {
      const [vault, ...remainingVaults] = vaults;

      if (vault) {
        getVault(vault)
          .then((vault) => {
            resolve([vault, ...remainingVaults]);
          })
          .catch(() => {
            getVaults(remainingVaults).then(resolve).catch(reject);
          });
      } else {
        reject();
      }
    });
  };

  const setVault = (vault: VaultProps): void => {
    setState((prevState) => ({
      ...prevState,
      vaults: prevState.vaults.map((item) => ({
        ...item,
        joinAirdrop:
          item.uid === vault.uid ? vault.joinAirdrop : item.joinAirdrop,
      })),
    }));
  };

  const setVaults = (vaults: VaultProps[]): void => {
    localStorage.setItem(storageKey.VAULTS, JSON.stringify(vaults));
  };

  const useVault = (vault: VaultProps): void => {
    if (!loading) {
      setState((prevState) => ({ ...prevState, loading: true }));

      getVault(vault)
        .then((vault) => {
          setState((prevState) => ({ ...prevState, loading: false, vault }));
        })
        .catch(() => {
          setState((prevState) => ({ ...prevState, loading: false }));
        });
    }
  };

  const componentDidUpdate = () => {
    if (core) {
      const handleError = () => {
        setVaults([]);

        setState((prevState) => ({ ...prevState, loaded: true }));
      };

      try {
        const data = localStorage.getItem(storageKey.VAULTS);
        const vaults: VaultProps[] = data ? JSON.parse(data) : [];

        if (Array.isArray(vaults) && vaults.length) {
          getVaults(vaults)
            .then((vaults) => {
              setState((prevState) => ({
                ...prevState,
                vault: vaults[0],
                vaults,
                loaded: true,
              }));

              setVaults(vaults);
            })
            .catch(() => {
              handleError();
            });
        } else {
          handleError();
        }
      } catch {
        handleError();
      }
    }
  };

  const componentDidMount = () => {
    initWasm()
      .then((core) => {
        let currency: Currency;
        let language: Language;

        switch (localStorage.getItem(storageKey.CURRENCY)) {
          case Currency.AUD:
            currency = Currency.AUD;
            break;
          case Currency.CNY:
            currency = Currency.CNY;
            break;
          case Currency.CAD:
            currency = Currency.CAD;
            break;
          case Currency.EUR:
            currency = Currency.EUR;
            break;
          case Currency.GPB:
            currency = Currency.GPB;
            break;
          case Currency.JPY:
            currency = Currency.JPY;
            break;
          case Currency.RUB:
            currency = Currency.RUB;
            break;
          case Currency.SEK:
            currency = Currency.SEK;
            break;
          case Currency.SGD:
            currency = Currency.SGD;
            break;
          default:
            currency = Currency.USD;
            break;
        }

        switch (localStorage.getItem(storageKey.LANGUAGE)) {
          case Language.CROATIA:
            language = Language.CROATIA;
            break;
          case Language.DUTCH:
            language = Language.DUTCH;
            break;
          case Language.GERMAN:
            language = Language.GERMAN;
            break;
          case Language.ITALIAN:
            language = Language.ITALIAN;
            break;
          case Language.PORTUGUESE:
            language = Language.PORTUGUESE;
            break;
          case Language.RUSSIAN:
            language = Language.RUSSIAN;
            break;
          case Language.SPANISH:
            language = Language.SPANISH;
            break;
          default:
            language = Language.ENGLISH;
            break;
        }

        i18n.changeLanguage(language);

        setState((prevState) => ({
          ...prevState,
          wcRefrence: {
            [Chain.ARBITRUM]: core.CoinType.arbitrum,
            [Chain.AVALANCHE]: core.CoinType.avalancheCChain,
            [Chain.BASE]: core.CoinType.base,
            [Chain.BITCOIN]: core.CoinType.bitcoin,
            [Chain.BITCOINCASH]: core.CoinType.bitcoinCash,
            [Chain.BLAST]: core.CoinType.blast,
            [Chain.BSCCHAIN]: core.CoinType.smartChain,
            [Chain.CRONOSCHAIN]: core.CoinType.cronosChain,
            [Chain.DASH]: core.CoinType.dash,
            [Chain.DOGECOIN]: core.CoinType.dogecoin,
            [Chain.DYDX]: core.CoinType.dydx,
            [Chain.ETHEREUM]: core.CoinType.ethereum,
            [Chain.GAIACHAIN]: core.CoinType.cosmos,
            [Chain.KUJIRA]: core.CoinType.kujira,
            [Chain.LITECOIN]: core.CoinType.litecoin,
            [Chain.MAYACHAIN]: core.CoinType.thorchain,
            [Chain.OPTIMISM]: core.CoinType.optimism,
            [Chain.POLKADOT]: core.CoinType.polkadot,
            [Chain.POLYGON]: core.CoinType.polygon,
            [Chain.SOLANA]: core.CoinType.solana,
            [Chain.SUI]: core.CoinType.sui,
            [Chain.THORCHAIN]: core.CoinType.thorchain,
            [Chain.ZKSYNC]: core.CoinType.zksync,
          },
          core,
          currency,
        }));
      })
      .catch(() => {});
  };

  useEffect(componentDidMount, []);
  useEffect(componentDidUpdate, [core]);

  return (
    <VaultContext.Provider
      value={{
        addVault,
        changeCurrency,
        fetchTokens,
        setVault,
        useVault,
        toggleCoin,
        currency,
        tokens,
        vault,
        vaults,
      }}
    >
      {loaded ? children : <SplashScreen />}
      <Modal
        className="modal-preloader"
        closeIcon={false}
        footer={false}
        maskClosable={false}
        open={loading}
        title="Loading..."
        centered
      >
        <Spin size="large" />
      </Modal>
    </VaultContext.Provider>
  );
};

export default Component;

export const useVaultContext = (): VaultContext => {
  const context = useContext(VaultContext);

  if (!context) {
    throw new Error("useVaultContext must be used within a VaultProvider");
  }

  return context;
};
