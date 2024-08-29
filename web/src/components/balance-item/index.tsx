import { FC, useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { Button, Spin, Tooltip, message } from "antd";
import { Truncate } from "@re-dev/react-truncate";

import { useVaultContext } from "context";
import { CoinProps } from "utils/interfaces";
import constantModals from "modals/constant-modals";
import constantPaths from "routes/constant-paths";

import {
  CaretRightOutlined,
  CopyOutlined,
  CubeOutlined,
  QRCodeOutlined,
} from "icons";
import QRCode from "modals/qr-code";
import { currencySymbol, exploreToken } from "utils/constants";

interface InitialState {
  assets: CoinProps[];
}

const Component: FC<CoinProps> = ({ address, balance, chain, ticker }) => {
  const initialState: InitialState = { assets: [] };
  const [state, setState] = useState(initialState);
  const { assets } = state;
  const { currency, vault } = useVaultContext();
  const [messageApi, contextHolder] = message.useMessage();

  const handleCopy = () => {
    navigator.clipboard
      .writeText(address)
      .then(() => {
        messageApi.open({
          type: "success",
          content: "Address copied to clipboard",
        });
      })
      .catch(() => {
        messageApi.open({
          type: "error",
          content: "Failed to copy address",
        });
      });
  };

  const componentDidUpdate = (): void => {
    const assets = vault
      ? vault.coins.filter((coin) => coin.chain === chain)
      : [];

    setState((prevState) => ({ ...prevState, assets }));
  };

  const componentDidMount = (): void => {};

  useEffect(componentDidMount, []);
  useEffect(componentDidUpdate, [address]);

  return (
    <>
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
          {assets.length > 1
            ? `${assets.length} assets`
            : balance.toString().split(".")[1]?.length > 8
            ? balance.toFixed(8)
            : balance}
        </span>
        <span className="amount">
          {vault
            ? `${currencySymbol[currency]}${vault.coins
                .filter((coin) => coin.chain === chain)
                .reduce((acc, coin) => acc + coin.balance * coin.value, 0)
                .toFixed(2)}`
            : `${currencySymbol[currency]}0.00`}
        </span>
        <div className="actions">
          <Tooltip title="Copy Address">
            <Button type="link" onClick={handleCopy}>
              <CopyOutlined />
            </Button>
          </Tooltip>
          <Tooltip title="View QRCode">
            <Link to={`#${constantModals.QR_CODE}`}>
              <QRCodeOutlined />
            </Link>
          </Tooltip>
          <Tooltip title="Link to Address">
            <a
              href={`${exploreToken[chain]}${address}`}
              rel="noopener noreferrer"
              target="_blank"
            >
              <CubeOutlined />
            </a>
          </Tooltip>
        </div>
        <Link
          to={`${constantPaths.balance}/${chain.toLocaleLowerCase()}`}
          className="arrow"
        >
          <CaretRightOutlined />
        </Link>
      </div>

      <QRCode address={address} />

      {contextHolder}
    </>
  );
};

export default Component;
