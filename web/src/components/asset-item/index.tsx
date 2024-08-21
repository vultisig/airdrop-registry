import { FC, useEffect } from "react";
import { Spin } from "antd";

import { useVaultContext } from "context";
import { Coin } from "context/interfaces";

const Component: FC<Coin.Params> = ({ balance, chain, ticker, value }) => {
  const { getBalance } = useVaultContext();

  const componentDidUpdate = (): void => {
    getBalance(chain, ticker)
      .then(() => {})
      .catch(() => {});
  };

  useEffect(componentDidUpdate, [ticker]);

  return (
    <div className="asset-item">
      <div className="token">
        <img
          src={`/coins/${ticker.toLocaleLowerCase()}.svg`}
          alt="bitcoin"
          className="logo"
        />
        <span className="name">{ticker}</span>
      </div>
      <span className="balance">{balance || <Spin />}</span>
      <span className="value">{value ? `$${value}` : <Spin />}</span>
    </div>
  );
};

export default Component;
