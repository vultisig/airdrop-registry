import { FC, useEffect, useState } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import { Button, Form, List, Modal, Switch } from "antd";

import { useVaultContext } from "context";
import { Chain } from "context/constants";
import constantModals from "modals/constant-modals";
import "./index.scss";
import Link from "antd/es/typography/Link";
import api from "utils/api";

interface InitialState {
  loading: Chain | null;
  submiting: boolean;
  visible: boolean;
}

const Component: FC = () => {
  const initialState: InitialState = {
    loading: null,
    submiting: false,
    visible: false,
  };
  const [state, setState] = useState(initialState);
  const { visible, submiting } = state;
  const { vaults } = useVaultContext();
  const { hash } = useLocation();
  const [form] = Form.useForm();
  const navigate = useNavigate();

  const handleSubmit = () => {
    navigate(-1);
    if (!submiting) {
      setState((prevState) => ({ ...prevState, submitting: true }));
      form
        .validateFields()
        .then((values) => {
          vaults.forEach((vault) => {
            const params = (({
              uid,
              name,
              publicKeyEcdsa,
              publicKeyEddsa,
              hexChainCode,
            }) => ({
              uid,
              name,
              publicKeyEcdsa,
              publicKeyEddsa,
              hexChainCode,
            }))(vault);
            // exit
            if (vault.joinAirdrop == true && values[vault.uid] == false) {
              api.airdrop
                .exit(params)
                .then((res) => {
                  if (res.status == 200) vault.joinAirdrop = false;

                  setState((prevState) => ({
                    ...prevState,
                    submitting: false,
                  }));
                })
                .catch((err) => {
                  console.error(err);
                  setState((prevState) => ({
                    ...prevState,
                    submitting: false,
                  }));
                });
            } else if (
              vault.joinAirdrop == false &&
              values[vault.uid] == true
            ) {
              // join
              api.airdrop
                .join(params)
                .then((res) => {
                  if (res.status == 200) vault.joinAirdrop = true;
                  setState((prevState) => ({
                    ...prevState,
                    submitting: false,
                  }));
                })
                .catch((err) => {
                  console.error(err);
                  setState((prevState) => ({
                    ...prevState,
                    submitting: false,
                  }));
                });
            }
          });
        })
        .catch((err) => {
          console.error(err);
          setState((prevState) => ({
            ...prevState,
            submitting: false,
          }));
        });
    }
  };

  const componentDidUpdate = () => {
    switch (hash) {
      case `#${constantModals.JOIN_AIRDROP}`: {
        setState((prevState) => ({
          ...prevState,
          visible: true,
        }));
        let data: any = {};
        vaults.forEach((vault) => {
          data[vault.uid] = vault.joinAirdrop;
        });
        form.setFieldsValue(data);
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
    <Modal
      className="form-modal"
      title="Join AirDrop"
      centered={true}
      footer={
        <Button
          type="primary"
          shape="round"
          className="submit-btn"
          loading={submiting}
          onClick={handleSubmit}
        >
          Done
        </Button>
      }
      onCancel={() => navigate(-1)}
      maskClosable={false}
      open={visible}
      width={550}
    >
      <div className="airdrop-header">
        <h2 className="value">$20,000,000</h2>
        <div>Current Airdrop Value</div>
        <div>Expected Drop Date: March 2025</div>
      </div>
      <div className="airdrop-body">
        <Form form={form} onFinish={handleSubmit}>
          <List
            dataSource={vaults}
            renderItem={(item) => (
              <List.Item
                key={item.uid}
                className="list-item"
                extra={
                  <Form.Item valuePropName="checked" name={item.uid} noStyle>
                    <Switch />
                  </Form.Item>
                }
              >
                <List.Item.Meta title={item.name} />
              </List.Item>
            )}
          />
          <Button htmlType="submit" style={{ display: "none" }} />
        </Form>
        <div>You are registering your Public Keys and vault addresses.</div>
        <Link className="link-color">Inspect the code here.</Link>
        <br /> <br />
        <div>
          Your Airdrop Share is based on how long you have kept funds in
          Vultisig for. Only Layer1 assets and tokens on the 1inch Token List
          apply.
        </div>
        <br /> <br />
        <div>No other information is collected.</div>
        <Link className="link-color">
          Read the Founder Pledge on Privacy here.
        </Link>
        <br /> <br />
        <div>You can register as many times as you like. </div>
      </div>
    </Modal>
  );
};

export default Component;
