import { FC, useEffect } from "react";
import { Link } from "react-router-dom";
import { Button, Card, Dropdown, Empty, Input, MenuProps, Spin } from "antd";

import { useVaultContext } from "context";
import constantModals from "modals/constant-modals";
import constantPaths from "routes/constant-paths";

import { PlusFilled, RefreshOutlined } from "icons";
import BalanceItem from "components/balance-item";
import ChooseChain from "modals/choose-coin";
import JoinAirdrop from "modals/join-airdrop";

const Component: FC = () => {
  const { changeVault, vault, vaults } = useVaultContext();

  const componentDidUpdate = () => {};

  const componentDidMount = () => {};

  useEffect(componentDidMount, []);
  useEffect(componentDidUpdate, [vault?.uid]);

  const items: MenuProps["items"] = [
    ...vaults.map(({ name, uid }) => ({
      label: name,
      key: uid,
      onClick: () => changeVault(uid),
    })),
    {
      type: "divider",
    },
    {
      key: "1",
      label: <Link to={constantPaths.landing}>Add new vault</Link>,
    },
    {
      key: "2",
      label: <Link to={`#${constantModals.JOIN_AIRDROP}`}>Join Airdrop</Link>,
    },
  ];

  return (
    <>
      <div className="balance-page">
        <div className="breadcrumb">
          <Dropdown menu={{ items }} className="menu">
            <Input value={vault?.name || ""} readOnly />
          </Dropdown>
          <Button type="link">
            <RefreshOutlined />
          </Button>
        </div>
        <div className="balance">
          <span className="title">Total Balance</span>
          <span className="value">$0</span>
        </div>
        {vault ? (
          vault.coins.length ? (
            vault.coins
              .filter((coin) => coin.isNativeToken)
              .map(({ chain, ...res }) => (
                <BalanceItem key={chain} {...{ ...res, chain }} />
              ))
          ) : (
            <Card className="empty">
              <Empty description="Choose a chain..." />
            </Card>
          )
        ) : (
          <Spin />
        )}
        <Link to={`#${constantModals.CHOOSE_CHAIN}`} className="add">
          <PlusFilled /> Choose Chains
        </Link>
      </div>

      <ChooseChain />
      <JoinAirdrop />
    </>
  );
};

export default Component;
