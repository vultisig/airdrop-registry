import axios from "axios";

import CaseConverter from "utils/case-converter";
import VaultManager from "utils/vault-manager";

//import paths from "routes/constant-paths";

const api = axios.create({
  baseURL: import.meta.env.VITE_SERVER_ADDRESS,
  headers: { accept: "application/json" },
});

api.interceptors.request.use(
  (config) => {
    config.data = CaseConverter.toSnake(config.data);

    return config;
  },
  (error) => {
    return Promise.reject(error.response);
  }
);

api.interceptors.response.use(
  (response) => {
    response.data = CaseConverter.toCamel(response.data);

    return response;
  },
  ({ response }) => {
    //if (data?.message) message.error(data?.message);

    switch (response.status) {
      case 401:
        break;
      case 403:
        break;
      default:
        break;
    }

    return Promise.reject(response.status);
  }
);

export default {
  register: async (params: VaultManager.Vault) => {
    return await api.post("vault", params);
  },
  derivePublicKey: async (params: VaultManager.Derivation) => {
    return await api.post<{ publicKey: string }>("derive-public-key", params);
  },
};
