import { FC, useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { Button, Spin } from "antd";
import { Truncate } from "@re-dev/react-truncate";

import { useVaultContext } from "context";
import { Coin } from "context/interfaces";
import constantPaths from "routes/constant-paths";

import {
  CaretRightOutlined,
  CopyOutlined,
  CubeOutlined,
  QRCodeOutlined,
} from "icons";

interface InitialState {
  assets: Coin.Params[];
}

const Component: FC<Coin.Params> = ({
  address,
  balance,
  chain,
  ticker,
  value,
}) => {
  const initialState: InitialState = { assets: [] };
  const [state, setState] = useState(initialState);
  const { assets } = state;
  const { getBalance, vault } = useVaultContext();

  const componentDidUpdate = (): void => {
    const assets = vault
      ? vault.coins.filter((coin) => coin.chain === chain)
      : [];

    setState((prevState) => ({ ...prevState, assets }));

    getBalance(chain, ticker)
      .then(() => {})
      .catch(() => {});
  };

  const componentDidMount = (): void => {};

  useEffect(componentDidMount, []);
  useEffect(componentDidUpdate, [address]);

  return (
    <div className="chain-item">
      <div className="type">
        <img
          src={`/coins/${chain.toLocaleLowerCase()}.svg`}
          alt="bitcoin"
          className="logo"
        />
        <span className="name">{chain}</span>
        <span className="text">{ticker}</span>
      </div>
      <div className="key">
        {address ? (
          <Truncate end={10} middle>
            {address}
          </Truncate>
        ) : (
          <Spin />
        )}
      </div>
      <span className={`asset${assets.length > 1 ? " multi" : ""}`}>
        {assets.length > 1 ? `${assets.length} assets` : balance || <Spin />}
      </span>
      <span className="amount">{value ? `$${value}` : <Spin />}</span>
      <div className="actions">
        <Button type="link">
          <CopyOutlined />
        </Button>
        <Button type="link">
          <QRCodeOutlined />
        </Button>
        <Button type="link">
          <CubeOutlined />
        </Button>
      </div>
      <Link
        to={`${constantPaths.balance}/${chain.toLocaleLowerCase()}`}
        className="arrow"
      >
        <CaretRightOutlined />
      </Link>
    </div>
  );
};

export default Component;
