import { FC, useEffect } from "react";
import { Link } from "react-router-dom";
import { Button, Card, Dropdown, Empty, Input, MenuProps, Spin } from "antd";

import { useVaultContext } from "context";
import constantModals from "modals/constant-modals";
import constantPaths from "routes/constant-paths";

import { CaretRightOutlined, PlusCircleFilled, RefreshOutlined } from "icons";
import BalanceItem from "components/balance-item";
import ChooseToken from "modals/choose-token";
import JoinAirdrop from "modals/join-airdrop";
import { useTranslation } from "react-i18next";
import translation from "i18n/constant-keys";

const Component: FC = () => {
  const { t } = useTranslation();
  const { useVault, currency, vault, vaults } = useVaultContext();

  const componentDidUpdate = () => {};

  const componentDidMount = () => {};

  useEffect(componentDidMount, []);
  useEffect(componentDidUpdate, [vault?.uid]);

  const items: MenuProps["items"] = [
    ...vaults.map((vault) => ({
      label: vault.name,
      key: vault.uid,
      onClick: () => useVault(vault, currency),
    })),
    {
      type: "divider",
    },
    {
      key: "1",
      label: (
        <>
          <Link to={constantPaths.landing}>+ {t(translation.ADD_NEW_VAULT)}</Link>
          <CaretRightOutlined />
        </>
      ),
      className: "primary",
    },
    {
      key: "2",
      label: (
        <>
          <Link to={`#${constantModals.JOIN_AIRDROP}`}>{t(translation.JOIN_AIRDROP)}</Link>
          <CaretRightOutlined />
        </>
      ),
      className: "primary",
    },
  ];

  return (
    <>
      <div className="balance-page">
        <div className="breadcrumb">
          <Dropdown menu={{ items }} className="menu">
            <Input value={vault?.name || ""} readOnly />
          </Dropdown>
          {vault && (
            <Button type="link" onClick={() => useVault(vault, currency)}>
              <RefreshOutlined />
            </Button>
          )}
        </div>
        <div className="balance">
          <span className="title">{t(translation.TOTAL_BALANCE)}</span>
          <span className="value">
            {vault
              ? `$${vault.coins
                  .reduce((acc, coin) => acc + coin.balance * coin.value, 0)
                  .toFixed(2)}`
              : "$0.00"}
          </span>
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
          <PlusCircleFilled /> {t(translation.CHOOSE_CHAIN)}
        </Link>
      </div>

      <ChooseToken />
      <JoinAirdrop />
    </>
  );
};

export default Component;
