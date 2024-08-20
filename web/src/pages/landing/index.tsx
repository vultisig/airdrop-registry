import { FC, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Button, Upload, UploadProps } from "antd";

import { useVaultContext } from "context";
import { errorKey } from "context/constants";
import { FileProps, Vault } from "context/interfaces";
import { CloseOutlined } from "utils/icons";
import constantPaths from "routes/constant-paths";

type Status = "default" | "error" | "success";

interface InitialState {
  file?: FileProps;
  status: Status;
  vault?: Vault.Params;
}

const Component: FC = () => {
  const initialState: InitialState = { status: "default" };
  const [state, setState] = useState(initialState);
  const { file, status, vault } = state;
  const { addVault, qrReader } = useVaultContext();
  const navigate = useNavigate();

  const handleStart = (): void => {
    if (vault && status === "success") {
      addVault(vault)
        .then(() => {
          navigate(constantPaths.balance);
        })
        .catch(() => {
          setState((prevState) => ({ ...prevState, status: "error" }));
        });
    }
  };

  const handleRemove = (): void => {
    setState(initialState);
  };

  const handleUpload = (file: File): false => {
    setState(initialState);

    qrReader(file)
      .then(({ file, vault }) => {
        setState((prevState) => ({
          ...prevState,
          file,
          vault,
          status: "success",
        }));
      })
      .catch(({ error, file }) => {
        setState((prevState) => ({ ...prevState, file, status: "error" }));

        switch (error) {
          case errorKey.INVALID_EXTENSION:
            console.error("Invalid file extension");
            break;
          case errorKey.INVALID_FILE:
            console.error("Invalid file");
            break;
          case errorKey.INVALID_QRCODE:
            console.error("Invalid qr code");
            break;
          case errorKey.INVALID_VAULT:
            console.error("Invalid vault data");
            break;
          default:
            console.error("Someting is wrong");
            break;
        }
      });

    return false;
  };

  const props: UploadProps = {
    multiple: false,
    showUploadList: false,
    beforeUpload: handleUpload,
    fileList: [],
  };

  const componentDidMount = () => {};

  useEffect(componentDidMount, []);

  return (
    <div className="landing-page">
      <img src="/images/logo-type.svg" alt="logo" className="logo" />
      <div className="wrapper">
        <h2 className="heading">Upload your vault share to start</h2>
        <Upload.Dragger {...props} className={status}>
          {file ? (
            <>
              <Button type="link" className="close" onClick={handleRemove}>
                <CloseOutlined />
              </Button>
              <img src={file.data} className="image" alt="image" />
              <h3 className="name">{`${file.name} Uploaded`}</h3>
            </>
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
          If you didn’t save the QR code yet, you can find it in the app in the
          top right on the main screen
        </p>
        <span
          className={`btn${status !== "success" ? " disabled" : ""}`}
          onClick={handleStart}
        >
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
  );
};

export default Component;
