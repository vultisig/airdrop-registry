import { FC, useEffect, useState } from "react";
import { Link, useNavigate, useLocation } from "react-router-dom";
import { Button, Form, List, Modal, Switch } from "antd";

import { useVaultContext } from "context";
import constantModals from "modals/constant-modals";
import api from "utils/api";

interface InitialState {
  submitting: boolean;
  visible: boolean;
}

const Component: FC = () => {
  const initialState: InitialState = { submitting: false, visible: false };
  const [state, setState] = useState(initialState);
  const { visible, submitting } = state;
  const { setVault, vaults } = useVaultContext();
  const { hash } = useLocation();
  const [form] = Form.useForm();
  const navigate = useNavigate();

  const handleSubmit = () => {
    if (!submitting) {
      setState((prevState) => ({ ...prevState, submitting: true }));

      form
        .validateFields()
        .then((values) => {
          vaults.forEach((vault) => {
            if (vault.joinAirdrop !== values[vault.uid]) {
              const params = { ...vault, joinAirdrop: values[vault.uid] };

              const action = params.joinAirdrop
                ? api.airdrop.join(params)
                : api.airdrop.exit(params);

              action
                .then(() => {
                  setVault(params);
                })
                .catch(() => {});
            } else {
              setState((prevState) => ({ ...prevState, submitting: false }));
            }
          });

          setTimeout(() => {
            navigate(-1);
          }, 1000);
        })
        .catch(() => {
          setState((prevState) => ({ ...prevState, submitting: false }));
        });
    }
  };

  const componentDidUpdate = () => {
    switch (hash) {
      case `#${constantModals.JOIN_AIRDROP}`: {
        setState((prevState) => ({ ...prevState, visible: true }));

        let data: any = {};

        vaults.forEach((vault) => (data[vault.uid] = vault.joinAirdrop));

        form.setFieldsValue(data);

        break;
      }
      default: {
        if (visible) form.resetFields();

        setState(initialState);

        break;
      }
    }
  };

  useEffect(componentDidUpdate, [hash]);

  return (
    <Modal
      className="join-airdrop"
      title="Join AirDrop"
      centered={true}
      footer={
        <Button
          loading={submitting}
          onClick={handleSubmit}
          shape="round"
          size="large"
          type="primary"
          block
        >
          Done
        </Button>
      }
      onCancel={() => navigate(-1)}
      maskClosable={false}
      open={visible}
      width={550}
    >
      <span className="value">$20,000,000</span>
      <span className="text">Current Airdrop Value</span>
      <span className="text">Expected Drop Date: March 2025</span>

      <Form form={form} onFinish={handleSubmit}>
        <List
          dataSource={vaults}
          renderItem={(item, index) => (
            <List.Item
              key={index}
              className="list-item"
              extra={
                <Form.Item valuePropName="checked" name={item.uid} noStyle>
                  <Switch loading={submitting} />
                </Form.Item>
              }
            >
              <List.Item.Meta title={item.name} />
            </List.Item>
          )}
        />
        <Button htmlType="submit" style={{ display: "none" }} />
      </Form>

      <div className="hint">
        <span className="desc">
          You are registering your Public Keys and vault addresses.
        </span>
        <Link
          to="https://github.com/vultisig/airdrop-registry"
          className="link"
        >
          Inspect the code here.
        </Link>
      </div>

      <div className="hint">
        <span className="desc">
          Your Airdrop Share is based on how long you have kept funds in
          Vultisig for. Only Layer1 assets and tokens on the 1inch Token List
          apply.
        </span>
      </div>

      <div className="hint">
        <span className="desc">No other information is collected.</span>
        <Link
          to="https://github.com/vultisig/docs/blob/main/other/privacy.md"
          className="link"
        >
          Read the Privacy Policy here
        </Link>
      </div>

      <div className="hint">
        <span className="desc">
          You can register as many times as you like.
        </span>
      </div>
    </Modal>
  );
};

export default Component;
