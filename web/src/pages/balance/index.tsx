import { FC, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Button, Select } from "antd";

import { useVaultContext } from "context";
import { ChainProps } from "context/interfaces";
import { PlusFilled, RefreshOutlined } from "utils/icons";
import constantPaths from "routes/constant-paths";

import ChainItem from "components/chain-item";

interface InitialState {
  chains: ChainProps[];
}

const Component: FC = () => {
  const initialState: InitialState = { chains: [] };
  const [state, setState] = useState(initialState);
  const { chains } = state;
  const { changeVault, vault, vaults } = useVaultContext();
  const navigate = useNavigate();

  const handleSelect = (uid: string) => {
    uid === "vault" ? navigate(constantPaths.landing) : changeVault(uid);
  };

  const componentDidUpdate = () => {
    if (vault) {
      if (Array.isArray(vault.coins) && vault.coins.length) {
        const chains: ChainProps[] = vault.coins
          ?.filter((coin) => coin.isNativeToken)
          .map(({ address, decimals, chain, ticker }) => ({
            address,
            assets:
              vault.coins?.filter((coin) => coin.chain === chain).length || 0,
            decimals,
            name: chain,
            ticker,
          }));

        setState((prevState) => ({ ...prevState, chains }));
      } else {
        setState((prevState) => ({ ...prevState, chains: [] }));
      }
    }
  };

  const componentDidMount = () => {};

  useEffect(componentDidMount, []);
  useEffect(componentDidUpdate, [vault?.uid]);

  return (
    <div className="balance-page">
      <div className="breadcrumb">
        <Select
          onChange={handleSelect}
          options={[
            ...vaults.map(({ name, uid }) => ({ label: name, value: uid })),
            { label: "Add new vault", value: "vault" },
          ]}
          popupClassName="vault-select-popup"
          rootClassName="vault-select"
          value={vault?.uid}
        />
        <Button type="link">
          <RefreshOutlined />
        </Button>
      </div>
      <div className="balance">
        <span className="title">Total Balance</span>
        <span className="value">$365,899.00</span>
      </div>
      {chains.map(({ name, ...res }) => (
        <ChainItem key={name} {...{ ...res, name }} />
      ))}
      <Button type="link" className="add">
        <PlusFilled /> Choose Chains
      </Button>
    </div>
  );
};

export default Component;
