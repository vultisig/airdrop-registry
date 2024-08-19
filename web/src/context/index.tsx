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

import { Chain, ErrorKey, balanceAPIs } from "context/constants";
import { Coin, FileProps, QRCodeProps, Vault } from "context/interfaces";
import { toCamelCase } from "utils/case-converter";
import api from "utils/api";

import SplashScreen from "components/splash-screen";

interface VaultContext {
  addVault: (vault: Vault.Params) => Promise<void>;
  changeVault: (uid: string) => void;
  getAddress: (value: Chain) => Promise<string>;
  getBalance: (address: string, chain: Chain) => Promise<number>;
  qrReader: (file: File) => Promise<QRCodeProps>;
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

  const getECDSAAddress = (chain: Chain, prefix?: string): Promise<string> => {
    return new Promise((resolve, reject) => {
      const coin = coinRef[chain];

      if (coin && core && vault) {
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

  const getEDDSAAdress = (chain: Chain): Promise<string> => {
    return new Promise((resolve, reject) => {
      const coin = coinRef[chain];

      if (coin && core && vault) {
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

  const getAddress = (chain: Chain): Promise<string> => {
    return new Promise((resolve, reject) => {
      switch (chain) {
        // EDDSA
        case Chain.POLKADOT:
        case Chain.SOLANA:
        case Chain.SUI: {
          getEDDSAAdress(chain)
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
          getECDSAAddress(chain, "maya")
            .then((address) => {
              resolve(address);
            })
            .catch(() => {
              reject();
            });

          break;
        }
        default: {
          getECDSAAddress(chain)
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

  const getBalance = (address: string, chain: Chain): Promise<number> => {
    return new Promise((resolve, reject) => {
      const path = balanceAPIs[chain];

      switch (chain) {
        // Cosmos
        case Chain.DYDX:
        case Chain.GAIACHAIN:
        case Chain.KUJIRA:
        case Chain.MAYACHAIN:
        case Chain.THORCHAIN: {
          api.balance
            .cosmos(`${path}/${address}`)
            .then(({ data: { balances } }) => {
              const [balance] = balances;

              if (balance && balance.amount) {
                resolve(parseInt(balance.amount));
              } else {
                resolve(0);
              }
            })
            .catch(() => {
              resolve(0);
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
              method: "eth_getBalance",
              params: [address, "latest"],
              id: 1,
            })
            .then(({ data: { result } }) => {
              resolve(parseInt(result, 16));
            })
            .catch(() => {
              resolve(0);
            });

          break;
        }
        case Chain.POLKADOT: {
          api.balance
            .polkadot(path, { key: address })
            .then(({ data: { data } }) => {
              if (data && data.account && data.account.balance) {
                const balance = data.account.balance.replace(".", "");

                resolve(parseInt(balance));
              } else {
                resolve(0);
              }
            })
            .catch(() => {
              resolve(0);
            });

          break;
        }
        case Chain.SOLANA: {
          api.balance
            .solana(path, {
              jsonrpc: "2.0",
              method: "getBalance",
              params: [address],
              id: 1,
            })
            .then(({ data: { result } }) => {
              resolve(parseInt(result, 16));
            })
            .catch(() => {
              resolve(0);
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
            .utxo(`${path}/${address}?state=latest`)
            .then(({ data: { data } }) => {
              if (
                data &&
                data[address] &&
                data[address].address &&
                typeof data[address].address.balance === "number"
              ) {
                resolve(data[address].address.balance);
              } else {
                resolve(0);
              }
            })
            .catch(() => {
              resolve(0);
            });

          break;
        }
        default:
          reject();
          break;
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
            .catch(() => {
              reject();
            });
        })
        .catch(() => {
          reject(ErrorKey.ALREADY_REGISTERED);
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
          setVault({ ...vault, ...data });

          resolve();
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
          reject(ErrorKey.INVALID_QRCODE);
        }
      };

      image.onerror = () => {
        reject(ErrorKey.INVALID_FILE);
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
        reject(ErrorKey.INVALID_FILE);
      };

      if (imageFormats.indexOf(file.type) >= 0) {
        reader.readAsDataURL(file);
      } else {
        reject(ErrorKey.INVALID_EXTENSION);
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

  const componentDidMount = () => {
    initWasm()
      .then((core) => {
        const storage = localStorage.getItem("vaults");
        const vaults: Vault.Props[] = storage ? JSON.parse(storage) : [];

        setState((prevState) => ({
          ...prevState,
          coinRef: {
            [Chain.ARBITRUM]: core.CoinType.arbitrum,
            [Chain.AVALANCHE]: core.CoinType.avalancheCChain,
            [Chain.BASE]: core.CoinType.base,
            [Chain.BITCOIN]: core.CoinType.bitcoin,
            [Chain.BITCOINCASH]: core.CoinType.bitcoinCash,
            [Chain.BLAST]: core.CoinType.blast,
            [Chain.BSCCHAIN]: core.CoinType.binance,
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

        if (Array.isArray(vaults) && vaults.length) {
          const promises = vaults.map((vault) => getVault(vault));

          Promise.all(promises)
            .then(() => {})
            .catch(() => {})
            .finally(() => {
              setState((prevState) => ({ ...prevState, loaded: true }));
            });
        } else {
          setState((prevState) => ({ ...prevState, loaded: true }));
        }
      })
      .catch(() => {});
  };

  useEffect(componentDidMount, []);

  return (
    <VaultContext.Provider
      value={{
        addVault,
        changeVault,
        getAddress,
        getBalance,
        qrReader,
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
