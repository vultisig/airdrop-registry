import { FC } from "react";

import { useVaultContext } from "context";
import { CoinProps } from "utils/interfaces";

const Component: FC<CoinProps> = ({ balance, logo, ticker, value }) => {
  const { currency } = useVaultContext();

  return (
    <div className="asset-item">
      <div className="token">
        <img
          src={logo || `/coins/${ticker.toLocaleLowerCase()}.svg`}
          alt="bitcoin"
          className="logo"
        />
        <span className="name">{ticker}</span>
      </div>
      <span className="balance">{balance.toBalanceFormat()}</span>
      <span className="value">
        {balance
          ? (balance * value).toValueFormat(currency)
          : (0).toValueFormat(currency)}
      </span>
    </div>
  );
};

export default Component;
