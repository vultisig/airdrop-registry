import type { CoinType } from "@trustwallet/wallet-core/dist/src/wallet-core";

import { Chain } from "context/constants";

export namespace Balance {
  export namespace Cosmos {
    export interface Props {
      balances: {
        denom: string;
        amount: string;
      }[];
    }
  }

  export namespace EVM {
    export interface Params {
      jsonrpc: string;
      method: string;
      params: [string, string];
      id: number;
    }

    export interface Props {
      id: number;
      jsonrpc: string;
      result: string;
    }
  }

  export namespace Polkadot {
    export interface Params {
      key: string;
    }

    export interface Props {
      data: { account: { balance: string } };
    }
  }

  export namespace Solana {
    export interface Params {
      jsonrpc: string;
      method: string;
      params: [string];
      id: number;
    }

    export interface Props {
      id: number;
      jsonrpc: string;
      result: string;
    }
  }

  export namespace UTXO {
    export interface Props {
      data: any;
    }
  }
}

export namespace Coin {
  export interface MetaData {
    chain: Chain;
    contractAddress: string;
    decimals: number;
    isNative: boolean;
    providerId: string;
    ticker: string;
  }

  export interface Props {
    address: string;
    chain: Chain;
    decimals: number;
    hexPublicKey: string;
    isNativeToken: boolean;
    priceProviderId: string;
    ticker: string;
  }

  export interface Reference {
    [Chain.ARBITRUM]?: CoinType;
    [Chain.AVALANCHE]?: CoinType;
    [Chain.BASE]?: CoinType;
    [Chain.BITCOIN]?: CoinType;
    [Chain.BITCOINCASH]?: CoinType;
    [Chain.BLAST]?: CoinType;
    [Chain.BSCCHAIN]?: CoinType;
    [Chain.CRONOSCHAIN]?: CoinType;
    [Chain.DASH]?: CoinType;
    [Chain.DOGECOIN]?: CoinType;
    [Chain.DYDX]?: CoinType;
    [Chain.ETHEREUM]?: CoinType;
    [Chain.GAIACHAIN]?: CoinType;
    [Chain.KUJIRA]?: CoinType;
    [Chain.LITECOIN]?: CoinType;
    [Chain.MAYACHAIN]?: CoinType;
    [Chain.OPTIMISM]?: CoinType;
    [Chain.POLKADOT]?: CoinType;
    [Chain.POLYGON]?: CoinType;
    [Chain.SOLANA]?: CoinType;
    [Chain.SUI]?: CoinType;
    [Chain.THORCHAIN]?: CoinType;
    [Chain.ZKSYNC]?: CoinType;
  }
}

export namespace Derivation {
  export interface Params {
    publicKeyEcdsa: string;
    hexChainCode: string;
    derivePath: string;
  }

  export interface Props {
    publicKey: string;
  }
}

export namespace Vault {
  export interface Params {
    uid: string;
    name: string;
    publicKeyEcdsa: string;
    publicKeyEddsa: string;
    hexChainCode: string;
  }

  export interface Props {
    coins?: Coin.Props[];
    name: string;
    hexChainCode: string;
    joinAirdrop: boolean;
    publicKeyEcdsa: string;
    publicKeyEddsa: string;
    totalPoints: number;
    uid: string;
  }
}

export interface BalanceAPI {
  [Chain.ARBITRUM]: string;
  [Chain.AVALANCHE]: string;
  [Chain.BASE]: string;
  [Chain.BITCOIN]: string;
  [Chain.BITCOINCASH]: string;
  [Chain.BLAST]: string;
  [Chain.BSCCHAIN]: string;
  [Chain.CRONOSCHAIN]: string;
  [Chain.DASH]: string;
  [Chain.DOGECOIN]: string;
  [Chain.DYDX]: string;
  [Chain.ETHEREUM]: string;
  [Chain.GAIACHAIN]: string;
  [Chain.KUJIRA]: string;
  [Chain.LITECOIN]: string;
  [Chain.MAYACHAIN]: string;
  [Chain.OPTIMISM]: string;
  [Chain.POLKADOT]: string;
  [Chain.POLYGON]: string;
  [Chain.SOLANA]: string;
  [Chain.SUI]: string;
  [Chain.THORCHAIN]: string;
  [Chain.ZKSYNC]: string;
}

export interface ChainProps {
  address: string;
  assets: number;
  decimals: number;
  name: Chain;
  ticker: string;
}

export interface FileProps {
  data: string;
  name: string;
}

export interface QRCodeProps {
  file: FileProps;
  vault: Vault.Params;
}
