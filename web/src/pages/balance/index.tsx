import { FC, useEffect, useState } from "react";
import { Button, Select, Spin } from "antd";
import { Truncate } from "@re-dev/react-truncate";

import {
  CaretRightOutlined,
  CopyOutlined,
  CubeOutlined,
  PlusFilled,
  QRCodeOutlined,
  RefreshOutlined,
} from "utils/icons";
import VaultManager from "utils/vault-manager";

interface InitialState {
  chains: VaultManager.CoinMeta[];
}

const Chain: FC<VaultManager.CoinMeta> = ({
  icon,
  chainID: id,
  name,
  ticker,
}) => {
  const initialState = { address: "" };
  const [state, setState] = useState(initialState);
  const { address } = state;

  const componentDidMount = () => {
    VaultManager.getAddress(id)
      .then((address) => {
        setState((prevState) => ({ ...prevState, address }));
      })
      .catch(() => {});
  };

  useEffect(componentDidMount, []);

  return (
    <div className="chain">
      <div className="type">
        <img src={icon} alt="bitcoin" className="logo" />
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
      <span className="asset">12,000.12</span>
      <span className="amount">$65,899</span>
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
      <Button type="link" className="arrow">
        <CaretRightOutlined />
      </Button>
    </div>
  );
};

const Component: FC = () => {
  const initialState: InitialState = { chains: [] };
  const [state, setState] = useState(initialState);
  const { chains } = state;

  const componentDidMount = () => {
    VaultManager.init()
      .then(() => {
        const _chains = [
          VaultManager.Chain.BITCOIN,
          VaultManager.Chain.THORCHAIN,
          VaultManager.Chain.BSCCHAIN,
          VaultManager.Chain.ETHEREUM,
          VaultManager.Chain.SOLANA,
        ];

        VaultManager.getChains(_chains)
          .then((chains) => {
            setState((prevState) => ({ ...prevState, chains }));
          })
          .catch(() => {});
      })
      .catch(() => {});
  };

  useEffect(componentDidMount, []);

  return (
    <div className="balance-page">
      <div className="breadcrumb">
        <Select
          rootClassName="vault-select"
          popupClassName="vault-select-popup"
          defaultValue={0}
          options={[
            { label: "Main Vault", value: 0 },
            { label: "Test Vault", value: 1 },
          ]}
        />
        <Button type="link">
          <RefreshOutlined />
        </Button>
      </div>
      <div className="balance">
        <span className="title">Total Balance</span>
        <span className="value">$365,899.00</span>
      </div>
      {chains.map(({ chainID: id, ...res }) => (
        <Chain key={id} {...{ ...res, chainID: id }} />
      ))}
      <Button type="link" className="add">
        <PlusFilled /> Choose Chains
      </Button>
    </div>
  );
};

export default Component;
