import {
  FC,
  ReactNode,
  useState,
  createContext,
  useEffect,
  useContext,
} from "react";
import { initWasm, WalletCore } from "@trustwallet/wallet-core";
import jsQR from "jsqr";

import { Chain, balanceAPIs, coins, errorKey } from "context/constants";
import { Coin, FileProps, QRCodeProps, Vault } from "context/interfaces";
import { toCamelCase } from "utils/case-converter";
import api from "utils/api";

import SplashScreen from "components/splash-screen";

interface VaultContext {
  addVault: (vault: Vault.Params) => Promise<void>;
  changeVault: (uid: string) => void;
  getBalance: (chain: Chain, ticker: string) => Promise<void>;
  qrReader: (file: File) => Promise<QRCodeProps>;
  toggleCoin: (coin: Coin.Meta) => Promise<void>;
  vault?: Vault.Props;
  vaults: Vault.Props[];
}

interface InitialState {
  coinRef: Coin.Reference;
  core?: WalletCore;
  loaded: boolean;
  vaults: Vault.Props[];
  vault?: Vault.Props;
}

const VaultContext = createContext<VaultContext | undefined>(undefined);

const Component: FC<{ children: ReactNode }> = ({ children }) => {
  const initialState: InitialState = { coinRef: {}, loaded: false, vaults: [] };
  const [state, setState] = useState(initialState);
  const { coinRef, core, vault, vaults, loaded } = state;

  const getECDSAAddress = (
    chain: Chain,
    vault: Vault.Params,
    prefix?: string
  ): Promise<string> => {
    return new Promise((resolve, reject) => {
      const coin = coinRef[chain];

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

  const getEDDSAAdress = (
    chain: Chain,
    vault: Vault.Params
  ): Promise<string> => {
    return new Promise((resolve, reject) => {
      const coin = coinRef[chain];

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

  const getAddress = (chain: Chain, vault: Vault.Params): Promise<string> => {
    return new Promise((resolve, reject) => {
      switch (chain) {
        // EDDSA
        case Chain.POLKADOT:
        case Chain.SOLANA:
        case Chain.SUI: {
          getEDDSAAdress(chain, vault)
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
          getECDSAAddress(chain, vault, "maya")
            .then((address) => {
              resolve(address);
            })
            .catch(() => {
              reject();
            });

          break;
        }
        default: {
          getECDSAAddress(chain, vault)
            .then((address) => {
              resolve(address);
            })
            .catch(() => {
              reject();
            });

          break;
        }
      }
    });
  };

  const getBalance = (chain: Chain, ticker: string): Promise<void> => {
    return new Promise((resolve, reject) => {
      const uid = Math.floor(Math.random() * 10000);
      const path = balanceAPIs[chain];
      const coin = vault?.coins.find(
        (coin) => coin.chain === chain && coin.ticker === ticker
      );

      if (coin) {
        switch (chain) {
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
                  setBalance(coin, parseInt(balances[0].amount));
                } else {
                  setBalance(coin, 0);
                }

                resolve();
              })
              .catch(() => {
                resolve();
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
                setBalance(coin, parseInt(result, 16));

                resolve();
              })
              .catch(() => {
                resolve();
              });

            break;
          }
          case Chain.POLKADOT: {
            api.balance
              .polkadot(path, { key: coin.address })
              .then(({ data: { data } }) => {
                if (data && data.account && data.account.balance) {
                  const balance = data.account.balance.replace(".", "");

                  setBalance(coin, parseInt(balance));
                } else {
                  setBalance(coin, 0);
                }

                resolve();
              })
              .catch(() => {
                resolve();
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
              .then(({ data: { result } }) => {
                setBalance(coin, parseInt(result, 16));

                resolve();
              })
              .catch(() => {
                resolve();
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
                  setBalance(coin, data[coin.address].address.balance);
                } else {
                  setBalance(coin, 0);
                }

                resolve();
              })
              .catch(() => {
                resolve();
              });

            break;
          }
          default:
            reject();
            break;
        }
      } else {
        reject();
      }
    });
  };

  const setBalance = (coin: Coin.Params, value: number): void => {
    setState((prevState) => {
      if (prevState.vault) {
        const coins = prevState.vault.coins.map((item) => ({
          ...item,
          balance:
            item.chain === coin.chain && item.ticker === coin.ticker
              ? (value / Math.pow(10, item.decimals)).toString()
              : item.balance,
        }));

        return { ...prevState, vault: { ...prevState.vault, coins } };
      } else {
        return prevState;
      }
    });
  };

  const addCoin = (
    coin: Coin.Meta,
    vault: Vault.Params
  ): Promise<Coin.Params> => {
    return new Promise((resolve, reject) => {
      const _coin: Coin.Params = {
        address: "",
        balance: "",
        chain: coin.chain,
        contractAddress: coin.contractAddress,
        decimals: coin.decimals,
        hexPublicKey:
          coin.hexPublicKey === "ECDSA"
            ? vault.publicKeyEcdsa
            : vault.publicKeyEddsa,
        isNativeToken: coin.isNative,
        priceProviderId: coin.providerId,
        ticker: coin.ticker,
        value: "",
      };

      getAddress(coin.chain, vault)
        .then((address) => {
          _coin.address = address;

          api.coin
            .add(vault, _coin)
            .then(({ data: { coinId } }) => {
              resolve({ ..._coin, ID: coinId });
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

  const delCoin = (coin: Coin.Params, vault: Vault.Params): Promise<void> => {
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

  const toggleCoin = (coin: Coin.Meta): Promise<void> => {
    return new Promise((resolve, reject) => {
      if (vault) {
        const _coin = vault.coins.find(
          ({ chain, ticker }) => coin.chain === chain && coin.ticker === ticker
        );

        if (_coin) {
          delCoin(_coin, vault)
            .then(() => {
              setVault({
                ...vault,
                coins: vault.coins.filter(
                  ({ chain, ticker }) =>
                    coin.chain !== chain || coin.ticker !== ticker
                ),
              });

              resolve();
            })
            .catch(() => {
              reject();
            });
        } else {
          addCoin(coin, vault)
            .then((coin) => {
              setVault({ ...vault, coins: [...vault.coins, coin] });

              resolve();
            })
            .catch(() => {
              reject();
            });
        }
      } else {
        reject();
      }
    });
  };

  const addVault = (vault: Vault.Params): Promise<void> => {
    return new Promise((resolve, reject) => {
      api.vault
        .add(vault)
        .then(() => {
          getVault(vault)
            .then(() => {
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
                .then(() => {
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

  const changeVault = (value: string) => {
    const vault = vaults.find((item) => item.uid === value);

    if (vault) setState((prevState) => ({ ...prevState, vault }));
  };

  const getVault = (vault: Vault.Params): Promise<void> => {
    return new Promise((resolve, reject) => {
      api.vault
        .get(vault)
        .then(({ data }) => {
          if (data.coins.length) {
            setVault({ ...vault, ...data });

            resolve();
          } else {
            const promises = coins
              .filter((coin) => coin.isDefault)
              .map((coin) => addCoin(coin, vault));

            Promise.all(promises)
              .then((coins) => {
                setVault({ ...vault, ...data, coins });

                resolve();
              })
              .catch(() => {});
          }
        })
        .catch(() => {
          reject();
        });
    });
  };

  const setVault = (vault: Vault.Props): void => {
    const storage = localStorage.getItem("vaults");
    const vaults: Vault.Props[] = storage ? JSON.parse(storage) : [];

    if (Array.isArray(vaults) && vaults.length) {
      setState((prevState) => {
        const vaults = [
          ...prevState.vaults.filter((item) => item.uid !== vault.uid),
          vault,
        ];

        localStorage.setItem("vaults", JSON.stringify(vaults));

        return { ...prevState, vault, vaults };
      });
    } else {
      localStorage.setItem("vaults", JSON.stringify([vault]));

      setState((prevState) => ({ ...prevState, vault, vaults: [vault] }));
    }
  };

  const readQRCode = (data: string): Promise<Vault.Params> => {
    return new Promise((resolve, reject) => {
      const canvas = document.createElement("canvas");
      const ctx = canvas.getContext("2d");
      const image = new Image();

      image.src = data;

      image.onload = () => {
        canvas.width = image.width;
        canvas.height = image.height;

        ctx?.drawImage(image, 0, 0, image.width, image.height);

        const imageData = ctx?.getImageData(
          0,
          0,
          image.width,
          image.height
        )?.data;

        if (imageData) {
          const qrData = jsQR(imageData, image.width, image.height);

          if (qrData) {
            const vaultData: Vault.Params = toCamelCase(
              JSON.parse(qrData.data)
            );

            resolve(vaultData);
          } else {
            reject();
          }
        } else {
          reject(errorKey.INVALID_QRCODE);
        }
      };

      image.onerror = () => {
        reject(errorKey.INVALID_FILE);
      };
    });
  };

  const readImage = (file: File): Promise<FileProps> => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      const imageFormats: string[] = [
        "image/jpg",
        "image/jpeg",
        "image/png",
        "image/bmp",
      ];

      reader.onload = () => {
        resolve({
          data: (reader.result || "").toString(),
          name: file.name,
        });
      };

      reader.onerror = () => {
        reject(errorKey.INVALID_FILE);
      };

      if (imageFormats.indexOf(file.type) >= 0) {
        reader.readAsDataURL(file);
      } else {
        reject(errorKey.INVALID_EXTENSION);
      }
    });
  };

  const qrReader = (file: File): Promise<QRCodeProps> => {
    return new Promise((resolve, reject) => {
      readImage(file)
        .then((file) => {
          readQRCode(file.data)
            .then((vault) => {
              resolve({ file, vault });
            })
            .catch((error) => {
              reject({ file, error });
            });
        })
        .catch((error) => {
          reject(error);
        });
    });
  };

  const componentDidUpdate = () => {
    if (core) {
      const storage = localStorage.getItem("vaults");
      const vaults: Vault.Props[] = storage ? JSON.parse(storage) : [];

      if (Array.isArray(vaults) && vaults.length) {
        const promises = vaults.map((vault) => getVault(vault));

        Promise.all(promises)
          .then(() => {})
          .catch(() => {})
          .finally(() => {
            setState((prevState) => ({
              ...prevState,
              vault: prevState.vaults[0],
              loaded: true,
            }));
          });
      } else {
        setState((prevState) => ({ ...prevState, loaded: true }));
      }
    }
  };

  const componentDidMount = () => {
    initWasm()
      .then((core) => {
        setState((prevState) => ({
          ...prevState,
          coinRef: {
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
        changeVault,
        getBalance,
        qrReader,
        toggleCoin,
        vault,
        vaults,
      }}
    >
      {loaded ? children : <SplashScreen />}
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
