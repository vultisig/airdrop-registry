import { FC, useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { Button, Spin } from "antd";
import { Truncate } from "@re-dev/react-truncate";

import { useVaultContext } from "context";
import { Coin } from "context/interfaces";
import {
  CaretRightOutlined,
  CopyOutlined,
  CubeOutlined,
  QRCodeOutlined,
} from "utils/icons";
import constantPaths from "routes/constant-paths";

interface InitialState {
  balance: string;
  value: string;
}

const assets = 1;

const Component: FC<Coin.Params> = ({ address, chain, decimals, ticker }) => {
  const initialState: InitialState = { balance: "", value: "" };
  const [state, setState] = useState(initialState);
  const { balance, value } = state;
  const { getBalance } = useVaultContext();

  const componentDidUpdate = (): void => {
    setState(initialState);

    if (assets === 1) {
      getBalance(chain, address)
        .then((balance) => {
          setState((prevState) => ({
            ...prevState,
            balance: (balance / Math.pow(10, decimals)).toString(),
          }));
        })
        .catch(() => {});
    }
  };

  const componentDidMount = (): void => {};

  useEffect(componentDidMount, []);
  useEffect(componentDidUpdate, [address]);

  return (
    <div className="chain-item">
      <div className="type">
        <img
          src={`/coins/${chain.toLocaleLowerCase()}-${ticker.toLocaleLowerCase()}.svg`}
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
      <span className={`asset${assets > 1 ? " multi" : ""}`}>
        {assets > 1 ? `${assets} assets` : balance || <Spin />}
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
