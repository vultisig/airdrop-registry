import axios from "axios";

import { Balance, Coin, Derivation, Vault } from "context/interfaces";
import { toCamelCase, toSnakeCase } from "utils/case-converter";

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
    add: async (vault: Vault.Params, params: Coin.Params) => {
      return await api.post<Coin.Props>(
        `coin/${vault.publicKeyEcdsa}/${vault.publicKeyEddsa}`,
        params,
        { headers: { "x-hex-chain-code": vault.hexChainCode } }
      );
    },
    del: async (vault: Vault.Params, coin: Coin.Params) => {
      return await api.delete(
        `coin/${vault.publicKeyEcdsa}/${vault.publicKeyEddsa}/${coin.ID}`, //${coin.chain}-${coin.ticker}-${coin.address}
        { headers: { "x-hex-chain-code": vault.hexChainCode } }
      );
    },
  },
  vault: {
    add: async (params: Vault.Params) => {
      return await api.post("vault", params);
    },
    get: async ({ publicKeyEcdsa, publicKeyEddsa }: Vault.Params) => {
      return await api.get<Vault.Props>(
        `vault/${publicKeyEcdsa}/${publicKeyEddsa}`
      );
    },
  },
  derivePublicKey: async (params: Derivation.Params) => {
    return await api.post<Derivation.Props>("derive-public-key", params);
  },
};
