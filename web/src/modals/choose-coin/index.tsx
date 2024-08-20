import { FC, useEffect, useState } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import { Drawer, List, Spin, Switch } from "antd";

import { useVaultContext } from "context";
import { Chain, coins } from "context/constants";
import constantModals from "modals/constant-modals";
import { Coin } from "context/interfaces";

interface InitialState {
  chooseChain: boolean;
  loading: Chain | null;
  visible: boolean;
}

const Component: FC = () => {
  const initialState: InitialState = {
    chooseChain: false,
    loading: null,
    visible: false,
  };
  const [state, setState] = useState(initialState);
  const { chooseChain, loading, visible } = state;
  const { toggleCoin, vault } = useVaultContext();
  const { hash } = useLocation();
  const navigate = useNavigate();

  const handleToggle = (coin: Coin.Meta) => {
    if (loading === null) {
      setState((prevState) => ({ ...prevState, loading: coin.chain }));

      toggleCoin(coin)
        .then(() => {})
        .catch(() => {})
        .finally(() => {
          setState((prevState) => ({ ...prevState, loading: null }));
        });
    }
  };

  const componentDidUpdate = () => {
    switch (hash) {
      case `#${constantModals.CHOOSE_CHAIN}`: {
        setState((prevState) => ({
          ...prevState,
          chooseChain: true,
          visible: true,
        }));

        break;
      }
      case `#${constantModals.CHOOSE_COIN}`: {
        setState((prevState) => ({
          ...prevState,
          chooseChain: false,
          visible: false,
        }));

        break;
      }
      default: {
        setState(initialState);

        break;
      }
    }
  };

  useEffect(componentDidUpdate, [hash]);

  return (
    <Drawer
      footer={false}
      onClose={() => navigate(-1)}
      title={visible ? (chooseChain ? "Choose a Chain" : "Choose a Asset") : ""}
      maskClosable={false}
      open={visible}
      width={320}
    >
      {visible ? (
        <List
          dataSource={coins.filter((coin) => coin.isNative === chooseChain)}
          renderItem={(item) => {
            const checked = vault
              ? vault?.coins.findIndex((coin) => coin.chain === item.chain) >= 0
              : false;

            return (
              <List.Item
                key={item.chain}
                extra={
                  <Switch
                    checked={checked}
                    loading={item.chain === loading}
                    onClick={() => handleToggle(item)}
                  />
                }
              >
                <List.Item.Meta
                  avatar={
                    <img
                      src={`/coins/${item.chain.toLocaleLowerCase()}-${item.ticker.toLocaleLowerCase()}.svg`}
                      style={{ height: 48, width: 48 }}
                    />
                  }
                  title={item.chain}
                  description={item.ticker}
                />
              </List.Item>
            );
          }}
        />
      ) : (
        <Spin className="center-spin" />
      )}
    </Drawer>
  );
};

export default Component;
