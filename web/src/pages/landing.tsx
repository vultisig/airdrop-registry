"use client";

import { FC, useState } from "react";
import { Upload, message } from "antd";

import type { UploadFile, UploadProps } from "antd";

import VaultDecryptor from "modules/VaultDecryptor";

type Status = "default" | "error" | "success";

const Component: FC = () => {
  const initialState = {
    fileList: [] as UploadFile[],
    status: "default" as Status,
  };
  const [state, setState] = useState(initialState);
  const { fileList, status } = state;
  const [messageApi, contextHolder] = message.useMessage();
  //const [file] = fileList;

  const props: UploadProps = {
    multiple: false,
    showUploadList: false,
    onRemove: () => {
      setState((prevState) => ({ ...prevState, fileList: [] }));
    },
    beforeUpload: (file) => {
      if (VaultDecryptor.validateExtension(file.name)) {
        console.log(file.name);

        VaultDecryptor.readFile(file)
          .then((data) => {
            setState((prevState) => ({
              ...prevState,
              status: "default",
            }));

            if (VaultDecryptor.validateBase64(data)) {
              const _decodedContainer = VaultDecryptor.decodeData(data);
              const decodedContainer =
                VaultDecryptor.decodeContainer(_decodedContainer);

              console.table(decodedContainer);

              if (decodedContainer.vault) {
                if (decodedContainer.isEncrypted) {
                  messageApi.warning("isEncrypted");
                } else {
                  messageApi.success("isDecrypted");

                  const _decodedVault = VaultDecryptor.decodeData(
                    decodedContainer.vault
                  );

                  const decodedVault =
                    VaultDecryptor.decodeVault(_decodedVault);

                  console.table(decodedVault);

                  if (decodedVault.hexChainCode) {
                    setState((prevState) => ({
                      ...prevState,
                      status: "success",
                    }));
                  } else {
                    messageApi.error("invalid vault data");

                    setState((prevState) => ({
                      ...prevState,
                      status: "error",
                    }));
                  }
                }
              } else {
                messageApi.error("invalid file data");
              }
            }
          })
          .catch(() => {});
      } else {
        messageApi.error("invalid file");
      }

      return false;
    },
    fileList,
  };

  return (
    <>
      <div className="vault-upload">
        <img src="/images/logo-type.svg" alt="logo" className="logo" />
        <div className="wrapper">
          <h2 className="heading">Upload your vault share to start</h2>
          <Upload.Dragger {...props} className={status}>
            <img src="/images/qr-code.svg" className="qr-code" alt="qr" />
            <h3 className="title">Upload your QR code here</h3>
            <span className="text">
              Drop your file here or <u>upload it</u>
            </span>
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
            <img src="/images/app-store.png" alt="iPhone" />
            <span>iPhone</span>
          </li>
          <li>
            <img src="/images/google-play.png" alt="Android" />
            <span>Android</span>
          </li>
          <li>
            <img src="/images/github.png" alt="Mac" />
            <span>Mac</span>
          </li>
        </ul>
      </div>

      {contextHolder}
    </>
  );
};

export default Component;
