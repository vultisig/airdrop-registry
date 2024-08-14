import { FC } from "react";
import { Button } from "antd";
import { useNavigate } from "react-router-dom";
import { Truncate } from "@re-dev/react-truncate";

import {
  CaretRightOutlined,
  CopyOutlined,
  CubeOutlined,
  PlusFilled,
  QRCodeOutlined,
} from "utils/icons";

const Component: FC = () => {
  const navigate = useNavigate();

  return (
    <>
      <div className="asset-page">
        <div className="breadcrumb">
          <Button type="link" className="back" onClick={() => navigate(-1)}>
            <CaretRightOutlined />
          </Button>
          <h1>Ethereum</h1>
        </div>
        <div className="content">
          <div className="chain">
            <div className="type">
              <img src="/images/chain-ethereum.png" alt="ethereum" />
              Ethereum
            </div>
            <div className="key">
              <Truncate end={10} middle>
                0x0cb1D4a24292bB89862f599Ac5B10F42b6DE07e4
              </Truncate>
            </div>
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
          </div>
          <div className="asset">
            <span className="token">ETH</span>
            <span className="balance">1.1</span>
            <span className="value">$60,899</span>
          </div>
          <div className="asset">
            <span className="token">USDT</span>
            <span className="balance">1,000</span>
            <span className="value">$1,000</span>
          </div>
          <div className="asset">
            <span className="token">WBTC</span>
            <span className="balance">0.1</span>
            <span className="value">$4,000</span>
          </div>
        </div>
        <Button type="link" className="add">
          <PlusFilled /> Choose Tokens
        </Button>
      </div>
    </>
  );
};

export default Component;
