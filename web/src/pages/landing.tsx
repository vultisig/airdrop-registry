"use client";

import { FC, useEffect, useState } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import { Button, Form, Input, Modal, Space, Upload, message } from "antd";

import type { UploadFile, UploadProps } from "antd";

import { modals } from "utils/constants";
import VaultDecryptor from "utils/vault-decryptor";

type Status = "default" | "error" | "success";

interface InitialState {
  fileName?: string;
  filePath?: string;
  fileType?: number;
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
  const initialState: InitialState = {
    status: "default",
  };
  const [state, setState] = useState(initialState);
  const { fileName, filePath, fileType, status } = state;
  const [messageApi, contextHolder] = message.useMessage();
  const navigate = useNavigate();
  //const [file] = fileList;

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

      VaultDecryptor.readFile(file)
        .then(({ data, name, type }) => {
          setState((prevState) => ({
            ...prevState,
            fileName: name,
            filePath: data,
            fileType: type,
          }));

          switch (type) {
            case VaultDecryptor.Type.IMAGE: {
              VaultDecryptor.readImage(data);

              break;
            }
            case VaultDecryptor.Type.DATA: {
              // VaultDecryptor.decryptor(data, handlePassword)
              //   .then((resolve) => {
              //     console.log(resolve);
              //   })
              //   .catch((error) => {
              //     switch (error) {
              //       case VaultDecryptor.Error.ENCODING:
              //         messageApi.error("Invalid file encoding");
              //         break;
              //       case VaultDecryptor.Error.CONTAINER:
              //         messageApi.error("Invalid vault container data");
              //         break;
              //       case VaultDecryptor.Error.PASSWD_REQUIRED:
              //         messageApi.error("Password is required");
              //         break;
              //       case VaultDecryptor.Error.INVALID_PASSWD:
              //         messageApi.error("Invalid vault data");
              //         break;
              //       case VaultDecryptor.Error.VAULT:
              //         messageApi.error("Invalid vault data");
              //         break;
              //       default:
              //         messageApi.error("Someting is wrong");
              //         break;
              //     }
              //   });

              break;
            }
            default:
              break;
          }
        })
        .catch((error) => {
          switch (error) {
            case VaultDecryptor.Error.INVALID_EXTENSION:
              messageApi.error("Invalid file extension");
              break;
            case VaultDecryptor.Error.INVALID_FILE:
              messageApi.error("Invalid file");
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

  return (
    <>
      <div className="vault-upload">
        <img src="/images/logo-type.svg" alt="logo" className="logo" />
        <div className="wrapper">
          <h2 className="heading">Upload your vault share to start</h2>
          <Upload.Dragger {...props} className={status}>
            {fileType === VaultDecryptor.Type.IMAGE ? (
              <>
                <span className="close" onClick={handleRemove} />
                <img src={filePath} className="image" alt="image" />
                <h3 className="name">{`${fileName} Uploaded`}</h3>
              </>
            ) : fileType === VaultDecryptor.Type.DATA ? (
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
