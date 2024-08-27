import { FC } from "react";
import { Outlet, Link } from "react-router-dom";
import { Button, Dropdown, MenuProps } from "antd";
import { useTranslation } from "react-i18next";

import { useVaultContext } from "context";
import { Language, languageName } from "utils/constants";
import i18n from "i18n/config";
import translation from "i18n/constant-keys";
import constantModals from "modals/constant-modals";
import constantPaths from "routes/constant-paths";

import {
  ChainOutlined,
  CurrencyOutlined,
  GearOutlined,
  GlobeOutlined,
  QuestionOutlined,
  ShareOutlined,
  UserOutlined,
} from "icons";
import ChangeCurrency from "modals/change-currency";
import ChangeLanguage from "modals/change-language";

const Component: FC = () => {
  const { t } = useTranslation();
  const { currency } = useVaultContext();
  let language: Language;

  switch (i18n.language) {
    case Language.CROATIA:
      language = Language.CROATIA;
      break;
    case Language.DUTCH:
      language = Language.DUTCH;
      break;
    case Language.GERMAN:
      language = Language.GERMAN;
      break;
    case Language.ITALIAN:
      language = Language.ITALIAN;
      break;
    case Language.PORTUGUESE:
      language = Language.PORTUGUESE;
      break;
    case Language.RUSSIAN:
      language = Language.RUSSIAN;
      break;
    case Language.SPANISH:
      language = Language.SPANISH;
      break;
    default:
      language = Language.ENGLISH;
      break;
  }

  const items: MenuProps["items"] = [
    {
      key: "1",
      label: t(translation.VAULT_SETTINGS),
      icon: <GearOutlined />,
    },
    {
      key: "2",
      label: (
        <>
          <Link to={`#${constantModals.CHANGE_LANG}`}>
            {t(translation.LANGUAGE)}
          </Link>
          <span>{languageName[language]}</span>
        </>
      ),
      icon: <GlobeOutlined />,
    },
    {
      key: "3",
      label: (
        <>
          <Link to={`#${constantModals.CHANGE_CURRENCY}`}>
            {t(translation.CURRENCY)}
          </Link>
          <span>{currency}</span>
        </>
      ),
      icon: <CurrencyOutlined />,
    },
    {
      key: "4",
      label: t(translation.DEFAULT_CHAINS),
      icon: <ChainOutlined />,
    },
    {
      key: "5",
      label: (
        <a
          href="https://vultisig.com/faq"
          rel="noopener noreferrer"
          target="_blank"
        >
          {t(translation.FAQ)}
        </a>
      ),
      icon: <QuestionOutlined />,
    },
    {
      key: "6",
      type: "group",
      label: t(translation.OTHER),
      children: [
        {
          key: "6-1",
          label: (
            <a
              href="https://vultisig.com/vult"
              rel="noopener noreferrer"
              target="_blank"
            >
              {t(translation.VULT_TOKEN)}
            </a>
          ),
          icon: <img src="/images/logo.svg" alt="logo" />,
        },
        {
          key: "6-2",
          label: (
            <a
              href="https://vultisig.com/#store-section"
              rel="noopener noreferrer"
              target="_blank"
            >
              {t(translation.SHARE_APP)}
            </a>
          ),
          icon: <ShareOutlined />,
        },
      ],
    },
  ];

  return (
    <>
      <div className="default-layout">
        <div className="header">
          <Link to={constantPaths.root} className="logo">
            <img src="/images/logo-type.svg" alt={t(translation.LOGO)} />
          </Link>
          <Dropdown menu={{ items }} className="menu">
            <Button type="link">
              <UserOutlined />
            </Button>
          </Dropdown>
        </div>
        <Outlet />
      </div>

      <ChangeCurrency />
      <ChangeLanguage />
    </>
  );
};

export default Component;
