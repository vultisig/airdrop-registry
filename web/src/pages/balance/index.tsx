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

interface Chain {
  address?: string;
  icon: string;
  key: VaultManager.ChainType;
  name: string;
}
interface InitialState {
  chains: Chain[];
}

const Component: FC = () => {
  const initialState: InitialState = {
    chains: [
      {
        icon: "/images/chain-bitcoin.png",
        key: VaultManager.ChainType.BITCOIN,
        name: "Bitcoin",
      },
      {
        icon: "/images/chain-ethereum.png",
        key: VaultManager.ChainType.ETHEREUM,
        name: "Ethereum",
      },
      {
        icon: "/images/chain-solana.png",
        key: VaultManager.ChainType.SOLANA,
        name: "Solana",
      },
      {
        icon: "/images/chain-thor.png",
        key: VaultManager.ChainType.THORCHAIN,
        name: "Thorchain",
      },
    ],
  };
  const [state, setState] = useState(initialState);
  const { chains } = state;

  const componentDidMount = () => {
    VaultManager.initiate()
      .then(() => {
        chains.forEach((chain) => {
          VaultManager.getAddress(chain.key)
            .then((address) => {
              setState((prevState) => {
                const _chains = prevState.chains.map((obj) =>
                  obj.key === chain.key ? { ...obj, address } : obj
                );

                return { ...prevState, chains: _chains };
              });
            })
            .catch(() => {});
        });
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
      {chains.map(({ address, icon, key, name }) => (
        <div className="chain" key={key}>
          <div className="type">
            <img src={icon} alt="bitcoin" className="logo" />
            <span className="name">{name}</span>
            <span className="text">BTC</span>
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
      ))}
      <Button type="link" className="add">
        <PlusFilled /> Choose Chains
      </Button>
    </div>
  );
};

export default Component;
