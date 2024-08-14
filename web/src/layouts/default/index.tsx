import { FC } from "react";
import { Outlet, Link } from "react-router-dom";
import { Button, Dropdown } from "antd";

import type { MenuProps } from "antd";

import { UserOutlined } from "utils/icons";
import paths from "routes/constant-paths";

const Component: FC = () => {
  const items: MenuProps["items"] = [
    {
      key: "1",
      label: "Vault Settings",
    },
    {
      key: "2",
      label: "Language",
    },
    {
      key: "3",
      label: "Currency",
    },
    {
      key: "4",
      label: "Default Chains",
    },
    {
      key: "5",
      label: "FAQ",
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
        <Dropdown menu={{ items }} className="menu">
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
