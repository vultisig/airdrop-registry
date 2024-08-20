import { FC, useEffect } from "react";
import { Link } from "react-router-dom";
import { Button, Card, Dropdown, Empty, Input, MenuProps, Spin } from "antd";

import { useVaultContext } from "context";
import { PlusFilled, RefreshOutlined } from "utils/icons";
import constantModals from "modals/constant-modals";
import constantPaths from "routes/constant-paths";

import ChainItem from "components/chain-item";
import ChooseChain from "modals/choose-coin";

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
      label: "Join Airdrop",
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
            vault.coins.map(({ chain, ...res }) => (
              <ChainItem key={chain} {...{ ...res, chain }} />
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
    </>
  );
};

export default Component;
