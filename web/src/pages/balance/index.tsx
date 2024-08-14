import { FC } from "react";
import { Button, Select } from "antd";
import { Truncate } from "@re-dev/react-truncate";

import {
  CaretRightOutlined,
  CopyOutlined,
  CubeOutlined,
  PlusFilled,
  QRCodeOutlined,
  RefreshOutlined,
} from "utils/icons";

const Component: FC = () => {
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
      <div className="chain">
        <div className="type">
          <img src="/images/chain-bitcoin.png" alt="bitcoin" className="logo" />
          <span className="name">Bitcoin</span>
          <span className="text">BTC</span>
        </div>
        <div className="key">
          <Truncate end={10} middle>
            bc1psrjtwm7682v6nhx2uwfgcfelrennd7pcvqq7v6w
          </Truncate>
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
      <div className="chain">
        <div className="type">
          <img
            src="/images/chain-ethereum.png"
            alt="ethereum"
            className="logo"
          />
          <span className="name">Ethereum</span>
          <span className="text">BTC</span>
        </div>
        <div className="key">
          <Truncate end={10} middle>
            0x0cb1D4a24292bB89862f599Ac5B10F42b6DE07e4
          </Truncate>
        </div>
        <span className="asset multi">3 assets</span>
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
      <Button type="link" className="add">
        <PlusFilled /> Choose Chains
      </Button>
    </div>
  );
};

export default Component;
