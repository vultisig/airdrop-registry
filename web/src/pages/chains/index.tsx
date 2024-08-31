import { FC } from "react";
import { Link } from "react-router-dom";
import { useTranslation } from "react-i18next";
import {
  Button,
  Card,
  Dropdown,
  Empty,
  Input,
  MenuProps,
  Spin,
  Tooltip,
} from "antd";

import { useVaultContext } from "context";
import translation from "i18n/constant-keys";
import constantModals from "modals/constant-modals";
import constantPaths from "routes/constant-paths";

import { CaretRightOutlined, PlusCircleFilled, RefreshOutlined } from "icons";
import BalanceItem from "components/balance-item";
import ChooseChain from "modals/choose-chain";
import JoinAirdrop from "modals/join-airdrop";

const Component: FC = () => {
  const { t } = useTranslation();
  const { useVault, currency, vault, vaults } = useVaultContext();

  const items: MenuProps["items"] = [
    ...vaults.map((vault) => ({
      label: vault.name,
      key: vault.uid,
      onClick: () => useVault(vault),
    })),
    {
      type: "divider",
    },
    {
      key: "1",
      label: (
        <>
          <Link to={constantPaths.import}>
            + {t(translation.ADD_NEW_VAULT)}
          </Link>
          <CaretRightOutlined />
        </>
      ),
      className: "primary",
    },
    {
      key: "2",
      label: (
        <>
          <Link to={`#${constantModals.JOIN_AIRDROP}`}>
            {t(translation.JOIN_AIRDROP)}
          </Link>
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
            <Tooltip title="Refresh">
              <Button type="link" onClick={() => useVault(vault)}>
                <RefreshOutlined />
              </Button>
            </Tooltip>
          )}
        </div>
        <div className="balance">
          <span className="title">{t(translation.TOTAL_BALANCE)}</span>
          <span className="value">
            {vault
              ? vault.coins
                  .reduce((acc, coin) => acc + coin.balance * coin.value, 0)
                  .toValueFormat(currency)
              : (0).toValueFormat(currency)}
          </span>
        </div>
        {vault ? (
          vault.coins.length ? (
            vault.coins
              .filter((coin) => coin.isNativeToken)
              .slice()
              .sort((a, b) => b.totalValue - a.totalValue)
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

      <ChooseChain />
      <JoinAirdrop />
    </>
  );
};

export default Component;
