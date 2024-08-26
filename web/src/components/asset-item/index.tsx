import { FC, useEffect } from "react";

import { Coin } from "utils/interfaces";

const Component: FC<Coin.Props> = ({ balance, ticker, value }) => {
  const componentDidUpdate = (): void => {};

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
      <span className="balance">
        {balance.toString().split(".")[1]?.length > 8
          ? balance.toFixed(8)
          : balance}
      </span>
      <span className="value">
        {balance ? `$${(balance * value).toFixed(2)}` : `$0.00`}
      </span>
    </div>
  );
};

export default Component;
