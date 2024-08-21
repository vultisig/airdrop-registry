import { FC, useEffect, useState } from "react";
import { Button, Card, Empty, Spin } from "antd";
import { Link, useNavigate, useParams } from "react-router-dom";
import { Truncate } from "@re-dev/react-truncate";

import { useVaultContext } from "context";
import { chooseToken } from "context/constants";
import { Coin } from "context/interfaces";
import constantPaths from "routes/constant-paths";

import AssetItem from "components/asset-item";
import ChooseCoin from "modals/choose-coin";

import {
  CaretRightOutlined,
  CopyOutlined,
  CubeOutlined,
  PlusFilled,
  QRCodeOutlined,
} from "icons";
import constantModals from "modals/constant-modals";

interface InitialState {
  coin?: Coin.Params;
}

const Component: FC = () => {
  const initialState: InitialState = {};
  const [state, setState] = useState(initialState);
  const { coin } = state;
  const { chainKey } = useParams();
  const { vault } = useVaultContext();
  const navigate = useNavigate();

  const componentDidUpdate = () => {
    if (chainKey && vault) {
      const coin = vault.coins.find(
        (coin) =>
          coin.isNativeToken && coin.chain.toLocaleLowerCase() === chainKey
      );

      if (coin) {
        setState((prevState) => ({ ...prevState, coin }));
      } else {
        navigate(constantPaths.balance);
      }
    } else {
      navigate(constantPaths.balance);
    }
  };

  useEffect(componentDidUpdate, [chainKey]);

  return (
    <>
      <div className="asset-page">
        {coin ? (
          <>
            <div className="breadcrumb">
              <Button type="link" className="back" onClick={() => navigate(-1)}>
                <CaretRightOutlined />
              </Button>
              <h1>{coin.chain}</h1>
            </div>
            <div className="content">
              <div className="chain">
                <div className="type">
                  <img
                    src={`/coins/${coin.chain.toLocaleLowerCase()}.svg`}
                    alt={coin.chain}
                  />
                  {coin.chain}
                </div>
                <div className="key">
                  {coin.address ? (
                    <Truncate end={10} middle>
                      {coin.address}
                    </Truncate>
                  ) : (
                    <Spin />
                  )}
                </div>
                <span className="amount">$0</span>
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
              {vault ? (
                vault?.coins.length ? (
                  vault?.coins
                    .filter(
                      (coin) => coin.chain.toLocaleLowerCase() === chainKey
                    )
                    .map(({ ticker, ...res }) => (
                      <AssetItem key={ticker} {...{ ...res, ticker }} />
                    ))
                ) : (
                  <Card className="empty">
                    <Empty description="Choose a asset..." />
                  </Card>
                )
              ) : (
                <Spin />
              )}
            </div>
            {chooseToken[coin.chain] && (
              <Link to={`#${constantModals.CHOOSE_TOKEN}`} className="add">
                <PlusFilled /> Choose Tokens
              </Link>
            )}
          </>
        ) : (
          <Spin className="center-spin" />
        )}
      </div>

      <ChooseCoin />
    </>
  );
};

export default Component;
