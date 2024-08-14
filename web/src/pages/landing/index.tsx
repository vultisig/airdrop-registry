import { FC, useEffect, useState } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import { Button, Form, Input, Modal, Space, Upload, message } from "antd";
import { initWasm, WalletCore } from "@trustwallet/wallet-core";

import type { UploadProps } from "antd";

import { modals } from "utils/constants";
import VaultDecryptor from "utils/vault-decryptor";

type Status = "default" | "error" | "success";

interface InitialState {
  core?: WalletCore;
  file?: VaultDecryptor.FileProps;
  status: Status;
}

interface PasswdModalProps {
  onCancel: () => void;
  onConfirm: (passwd: string) => void;
}

const PasswdModal: FC<PasswdModalProps> = ({ onCancel, onConfirm }) => {
  const initialState = { visible: false };
  const [state, setState] = useState(initialState);
  const { visible } = state;
  const { hash } = useLocation();
  const [form] = Form.useForm();

  const handleConfirm = () => {
    form
      .validateFields()
      .then(({ password }) => {
        onConfirm(password);
      })
      .catch(() => {});
  };

  const componentDidUpdate = () => {
    switch (hash) {
      case `#${modals.PASSWORD}`: {
        setState((prevState) => ({ ...prevState, visible: true }));

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
      title="Enter Password"
      open={visible}
      onOk={handleConfirm}
      onCancel={() => onCancel()}
      footer={
        <Space>
          <Button onClick={() => onCancel()}>Cancel</Button>
          <Button type="primary" onClick={handleConfirm}>
            Confirm
          </Button>
        </Space>
      }
    >
      <Form form={form} layout="vertical" onFinish={handleConfirm}>
        <Form.Item
          name="password"
          label="Password"
          rules={[{ required: true }]}
        >
          <Input.Password />
        </Form.Item>

        <Button htmlType="submit" style={{ display: "none" }} />
      </Form>
    </Modal>
  );
};

const Component: FC = () => {
  const initialState: InitialState = { status: "default" };
  const [state, setState] = useState(initialState);
  const { core, file, status } = state;
  const [messageApi, contextHolder] = message.useMessage();
  const navigate = useNavigate();

  const handleAddress = (ecdsa: string, eddsa: string): void => {
    if (core) {
      const ecdsaBytes = core.HexCoding.decode(ecdsa);
      const eddsaBytes = core.HexCoding.decode(eddsa);

      const ecdsaKey = core.PublicKey.createWithData(
        ecdsaBytes,
        core.PublicKeyType.secp256k1
      );

      const eddsaKey = core.PublicKey.createWithData(
        eddsaBytes,
        core.PublicKeyType.ed25519
      );

      const bitcoin = core.AnyAddress.createWithPublicKey(
        ecdsaKey,
        core.CoinType.bitcoin
      ).description();

      const ethereum = core.AnyAddress.createWithPublicKey(
        ecdsaKey,
        core.CoinType.ethereum
      ).description();

      const solana = core.AnyAddress.createWithPublicKey(
        eddsaKey,
        core.CoinType.solana
      ).description();

      const thorchain = core.AnyAddress.createWithPublicKey(
        eddsaKey,
        core.CoinType.thorchain
      ).description();

      console.log(`bitcoin: ${bitcoin}`);
      console.log(`ethereum: ${ethereum}`);
      console.log(`solana: ${solana}`);
      console.log(`thorchain: ${thorchain}`);
    }
  };

  const handleData = (file: VaultDecryptor.FileProps) => {
    setState((prevState) => ({ ...prevState, file }));
  };

  const handleCancel = () => {
    navigate(-1);
  };

  const handleConfirm = (passwd: string) => {
    console.log(passwd);
  };

  const handlePassword = (): Promise<string> => {
    return new Promise((resolve, reject) => {
      if (true) {
        resolve("123");
      } else {
        reject("");
      }
    });
  };

  const handleRemove = (): void => {
    setState(initialState);
  };

  const props: UploadProps = {
    multiple: false,
    showUploadList: false,
    beforeUpload: (file) => {
      setState(initialState);

      VaultDecryptor.decryptor(file, handleData, handlePassword)
        .then((result) => {
          setState((prevState) => ({ ...prevState, status: "success" }));

          handleAddress(
            result.publicKeyEcdsa || result.public_key_ecdsa,
            result.publicKeyEddsa || result.public_key_eddsa
          );

          console.log(result);
        })
        .catch((error) => {
          setState((prevState) => ({ ...prevState, status: "error" }));

          switch (error) {
            case VaultDecryptor.Error.INVALID_CONTAINER:
              messageApi.error("Invalid vault container data");
              break;
            case VaultDecryptor.Error.INVALID_ENCODING:
              messageApi.error("Invalid file encode");
              break;
            case VaultDecryptor.Error.INVALID_EXTENSION:
              messageApi.error("Invalid file extension");
              break;
            case VaultDecryptor.Error.INVALID_FILE:
              messageApi.error("Invalid file");
              break;
            case VaultDecryptor.Error.PASSWD_REQUIRED:
              messageApi.error("Password is required");
              break;
            case VaultDecryptor.Error.INVALID_PASSWD:
              messageApi.error("Invalid vault data");
              break;
            case VaultDecryptor.Error.INVALID_QRCODE:
              messageApi.error("Invalid qr code");
              break;
            case VaultDecryptor.Error.INVALID_VAULT:
              messageApi.error("Invalid vault data");
              break;
            default:
              messageApi.error("Someting is wrong");
              break;
          }
        });

      return false;
    },
    fileList: [],
  };

  const componentDidMount = () => {
    initWasm()
      .then((core) => {
        setState((prevState) => ({ ...prevState, core }));
      })
      .catch((error) => {
        console.log(error);
      });
  };

  useEffect(componentDidMount, []);

  return (
    <>
      <div className="landing-page">
        <img src="/images/logo-type.svg" alt="logo" className="logo" />
        <div className="wrapper">
          <h2 className="heading">Upload your vault share to start</h2>
          <Upload.Dragger {...props} className={status}>
            {file?.type === VaultDecryptor.Type.IMAGE ? (
              <>
                <span className="close" onClick={handleRemove} />
                <img src={file.data} className="image" alt="image" />
                <h3 className="name">{`${file.name} Uploaded`}</h3>
              </>
            ) : file?.type === VaultDecryptor.Type.DATA ? (
              <></>
            ) : (
              <>
                <img src="/images/qr-code.svg" className="icon" alt="qr" />
                <h3 className="title">Upload your QR code here</h3>
                <span className="text">
                  Drop your file here or <u>upload it</u>
                </span>
              </>
            )}
          </Upload.Dragger>
          <p className="hint">
            If you didn’t save the QR code yet, you can find it in the app in
            the top right on the main screen
          </p>
          <span className={`btn${status !== "success" ? " disabled" : ""}`}>
            Start
          </span>
        </div>
        <p className="hint">Don’t have a vault yet? Download Vault now</p>
        <ul className="download">
          <li>
            <a
              href="https://testflight.apple.com/join/kpVufItl"
              target="_blank"
              rel="noopener noreferrer"
              className="image"
            >
              <img src="/images/app-store.png" alt="iPhone" />
            </a>
            <a
              href="https://testflight.apple.com/join/kpVufItl"
              target="_blank"
              rel="noopener noreferrer"
              className="text"
            >
              iPhone
            </a>
          </li>
          <li>
            <a
              href="https://play.google.com/store/apps/details?id=com.vultisig.wallet"
              target="_blank"
              rel="noopener noreferrer"
              className="image"
            >
              <img src="/images/google-play.png" alt="Android" />
            </a>
            <a
              href="https://play.google.com/store/apps/details?id=com.vultisig.wallet"
              target="_blank"
              rel="noopener noreferrer"
              className="text"
            >
              Android
            </a>
          </li>
          <li>
            <a
              href="https://github.com/vultisig/vultisig-ios/releases"
              target="_blank"
              rel="noopener noreferrer"
              className="image"
            >
              <img src="/images/github.png" alt="Mac" />
            </a>
            <a
              href="https://github.com/vultisig/vultisig-ios/releases"
              target="_blank"
              rel="noopener noreferrer"
              className="text"
            >
              Mac
            </a>
          </li>
        </ul>
      </div>

      <PasswdModal onCancel={handleCancel} onConfirm={handleConfirm} />

      {contextHolder}
    </>
  );
};

export default Component;
