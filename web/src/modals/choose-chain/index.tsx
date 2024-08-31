import { FC, useEffect, useState } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import { Drawer, Input, List, Spin, Switch } from "antd";

import { useVaultContext } from "context";
import { Chain, defTokens } from "utils/constants";
import { TokenProps } from "utils/interfaces";
import constantModals from "modals/constant-modals";

import { SearchOutlined } from "icons";

interface InitialState {
  loading: Chain | null;
  search: string;
  tokens: TokenProps[];
  visible: boolean;
}

const Component: FC = () => {
  const initialState: InitialState = {
    loading: null,
    search: "",
    tokens: [],
    visible: false,
  };
  const [state, setState] = useState(initialState);
  const { loading, search, tokens, visible } = state;
  const { toggleCoin, vault } = useVaultContext();
  const { hash } = useLocation();
  const navigate = useNavigate();

  const handleSearch = (value: string): void => {
    setState((prevState) => ({
      ...prevState,
      search: value.toLocaleLowerCase(),
    }));
  };

  const handleToggle = (coin: TokenProps): void => {
    if (vault && loading === null) {
      setState((prevState) => ({ ...prevState, loading: coin.chain }));

      toggleCoin(coin, vault)
        .then(() => {})
        .catch(() => {})
        .finally(() => {
          setState((prevState) => ({ ...prevState, loading: null }));
        });
    }
  };

  const componentDidUpdate = (): void => {
    switch (hash) {
      case `#${constantModals.CHOOSE_CHAIN}`: {
        setState((prevState) => ({
          ...prevState,
          tokens: defTokens.filter(({ isNative }) => isNative),
          visible: true,
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
      title={
        <Input
          onChange={(e) => handleSearch(e.target.value)}
          placeholder="Search..."
          prefix={<SearchOutlined />}
          value={search}
        />
      }
      closeIcon={false}
      open={visible}
      width={320}
    >
      {visible ? (
        <List
          dataSource={tokens.filter(
            ({ chain }) =>
              chain.toLocaleLowerCase().indexOf(search.toLocaleLowerCase()) >= 0
          )}
          renderItem={(item) => {
            const checked = vault
              ? vault?.coins.findIndex(
                  (coin) =>
                    coin.chain === item.chain && coin.ticker === item.ticker
                ) >= 0
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
                      src={`/coins/${item.chain.toLocaleLowerCase()}.svg`}
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
