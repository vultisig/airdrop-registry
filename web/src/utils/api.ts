import axios from "axios";

import { toCamelCase, toSnakeCase } from "utils/case-converter";
import { Currency } from "utils/constants";
import { Balance, CoinProps, Derivation, VaultProps } from "utils/interfaces";

//import paths from "routes/constant-paths";

const api = axios.create({
  baseURL: import.meta.env.VITE_SERVER_ADDRESS,
  headers: { accept: "application/json" },
});

api.interceptors.request.use(
  (config) => {
    config.data = toSnakeCase(config.data);

    return config;
  },
  (error) => {
    return Promise.reject(error.response);
  }
);

api.interceptors.response.use(
  (response) => {
    response.data = toCamelCase(response.data);

    return response;
  },
  ({ response }) => {
    return Promise.reject(response.data.error);
  }
);

export default {
  airdrop: {
    join: async (params: VaultProps) => {
      return await api.post("vault/join-airdrop", params);
    },
    exit: async (params: VaultProps) => {
      return await api.post("vault/exit-airdrop", params);
    },
  },
  balance: {
    cosmos: async (path: string) => {
      return await api.get<Balance.Cosmos.Props>(path);
    },
    evm: async (path: string, params: Balance.EVM.Params) => {
      return await api.post<Balance.EVM.Props>(path, params);
    },
    polkadot: async (path: string, params: Balance.Polkadot.Params) => {
      return await api.post<Balance.Polkadot.Props>(path, params);
    },
    solana: async (path: string, params: Balance.Solana.Params) => {
      return await api.post<Balance.Solana.Props>(path, params);
    },
    utxo: async (path: string) => {
      return await api.get<Balance.UTXO.Props>(path);
    },
  },
  coin: {
    add: async (vault: VaultProps, coin: CoinProps) => {
      return await api.post<{ coinId: number }>(
        `coin/${vault.publicKeyEcdsa}/${vault.publicKeyEddsa}`,
        {
          address: coin.address,
          chain: coin.chain,
          cmcId: coin.cmcId,
          contractAddress: coin.contractAddress,
          decimals: coin.decimals,
          hexPublicKey: coin.hexPublicKey,
          isNativeToken: coin.isNativeToken,
          logo: coin.logo,
          ticker: coin.ticker,
        },
        { headers: { "x-hex-chain-code": vault.hexChainCode } }
      );
    },
    cmc: async (address: string) => {
      return await api.get<{
        data: { [cmcId: string]: { id: number } };
      }>(
        `https://api.vultisig.com/cmc/v1/cryptocurrency/info?address=${address}&skip_invalid=true&aux=status`
      );
    },
    del: async (vault: VaultProps, coin: CoinProps) => {
      return await api.delete(
        `coin/${vault.publicKeyEcdsa}/${vault.publicKeyEddsa}/${coin.ID}`,
        { headers: { "x-hex-chain-code": vault.hexChainCode } }
      );
    },
    values: async (ids: number[], currency: Currency) => {
      return await api.get<{
        data: {
          [id: string]: { quote: { [currency: string]: { price: number } } };
        };
      }>(
        `https://api.vultisig.com/cmc/v2/cryptocurrency/quotes/latest?id=${ids
          .filter((id) => id > 0)
          .join(",")}&skip_invalid=true&aux=is_active&convert=${currency}`
      );
    },
  },
  vault: {
    add: async (params: VaultProps) => {
      return await api.post<VaultProps>("vault", params);
    },
    get: async ({ publicKeyEcdsa, publicKeyEddsa }: VaultProps) => {
      return await api.get<VaultProps>(
        `vault/${publicKeyEcdsa}/${publicKeyEddsa}`
      );
    },
  },
  derivePublicKey: async (params: Derivation.Params) => {
    return await api.post<Derivation.Props>("derive-public-key", params);
  },
  oneInch: async (id: number) => {
    return await api.get<{
      tokens: {
        [address: string]: {
          decimals: number;
          logoURI: string;
          name: string;
          symbol: string;
        };
      };
    }>(`https://api.vultisig.com/1inch/swap/v6.0/${id}/tokens`);
  },
};
