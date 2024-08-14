import axios from "axios";

//import paths from "routes/constant-paths";

interface PublicKeys {
  ecdsa: string;
  eddsa: string;
}

interface HexChainCode extends PublicKeys {
  hexChainCode: string;
}

const api = axios.create({
  baseURL: import.meta.env.VITE_SERVER_ADDRESS,
  headers: { accept: "application/json" },
});

api.interceptors.request.use(
  (config) => {
    return config;
  },
  (error) => {
    return Promise.reject(error.response);
  }
);

api.interceptors.response.use(
  ({ data }) => {
    return data;
  },
  ({ response: { status } }) => {
    //if (data?.message) message.error(data?.message);

    switch (status) {
      case 401:
        break;
      case 403:
        break;
      default:
        break;
    }

    return Promise.reject(status);
  }
);

const get = async (url: string, params = {}) => {
  return await api.get(url, { params });
};

const post = async (url: string, params = {}) => {
  return await api.post(url, params);
};

export default {
  stepOne: async (params: HexChainCode) => {
    return post("vault", params);
  },
  stepTwo: ({ ecdsa, eddsa }: PublicKeys) => {
    return get(`vault/:${ecdsa}/:${eddsa}/address`);
  },
  stepThree: ({ ecdsa, eddsa }: PublicKeys) => {
    return post(`vault/:${ecdsa}/:${eddsa}/address`);
  },
  stepFour: ({ ecdsa, eddsa }: PublicKeys) => {
    return get(`vault/:${ecdsa}/:${eddsa}/balance`);
  },
};
