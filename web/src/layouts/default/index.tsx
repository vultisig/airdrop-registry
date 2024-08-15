import { FC } from "react";
import { Outlet, Link } from "react-router-dom";
import { Button, Dropdown } from "antd";

import type { MenuProps } from "antd";

import {
  ChainOutlined,
  CurrencyOutlined,
  GearOutlined,
  GlobeOutlined,
  QuestionOutlined,
  ShareOutlined,
  UserOutlined,
} from "utils/icons";
import paths from "routes/constant-paths";

const Component: FC = () => {
  const items: MenuProps["items"] = [
    {
      key: "1",
      label: "Vault Settings",
      icon: <GearOutlined />,
    },
    {
      key: "2",
      label: "Language",
      icon: <GlobeOutlined />,
    },
    {
      key: "3",
      label: "Currency",
      icon: <CurrencyOutlined />,
    },
    {
      key: "4",
      label: "Default Chains",
      icon: <ChainOutlined />,
    },
    {
      key: "5",
      label: "FAQ",
      icon: <QuestionOutlined />,
    },
    {
      key: "6",
      type: "group",
      label: "Other",
      children: [
        {
          key: "6-1",
          label: "The $VULT Token",
        },
        {
          key: "6-2",
          label: "Share The App",
          icon: <ShareOutlined />,
        },
      ],
    },
  ];

  return (
    <div className="default-layout">
      <div className="header">
        <Link to={paths.root} className="logo">
          <img src="/images/logo-type.svg" alt="logo" />
        </Link>
        <Dropdown menu={{ items }} className="menu" >
          <Button type="link">
            <UserOutlined />
          </Button>
        </Dropdown>
      </div>
      <Outlet />
    </div>
  );
};

export default Component;
