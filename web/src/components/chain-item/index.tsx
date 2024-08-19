import { FC, useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { Button, Spin } from "antd";
import { Truncate } from "@re-dev/react-truncate";

import { useVaultContext } from "context";
import { ChainProps } from "context/interfaces";
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

const Component: FC<ChainProps> = ({
  address,
  assets,
  decimals,
  name,
  ticker,
}) => {
  const initialState: InitialState = { balance: "", value: "" };
  const [state, setState] = useState(initialState);
  const { balance, value } = state;
  const { getBalance } = useVaultContext();

  const componentDidMount = () => {
    if (assets === 1) {
      getBalance(address, name)
        .then((balance) => {
          setState((prevState) => ({
            ...prevState,
            balance: (balance / Math.pow(10, decimals)).toString(),
          }));
        })
        .catch(() => {});
    }
  };

  useEffect(componentDidMount, []);

  return (
    <div className="chain">
      <div className="type">
        <img
          src={`/coins/${ticker.toLocaleLowerCase()}.svg`}
          alt="bitcoin"
          className="logo"
        />
        <span className="name">{name}</span>
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
        to={`${constantPaths.balance}/${name.toLocaleLowerCase()}`}
        className="arrow"
      >
        <CaretRightOutlined />
      </Link>
    </div>
  );
};

export default Component;
