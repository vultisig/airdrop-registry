import { FC, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Button, Upload, UploadProps } from "antd";
import { useTranslation } from "react-i18next";
import { ReaderOptions, readBarcodesFromImageFile } from "zxing-wasm/reader";

import { useVaultContext } from "context";
import { toCamelCase } from "utils/case-converter";
import { errorKey } from "utils/constants";
import { FileProps, VaultProps } from "utils/interfaces";
import translation from "i18n/constant-keys";
import constantPaths from "routes/constant-paths";

import { CloseOutlined } from "icons";

interface InitialState {
  file?: FileProps;
  loading: boolean;
  status: "default" | "error" | "success";
  vault?: VaultProps;
}

const Component: FC = () => {
  const { t } = useTranslation();
  const initialState: InitialState = { loading: false, status: "default" };
  const [state, setState] = useState(initialState);
  const { file, loading, status, vault } = state;
  const { addVault } = useVaultContext();
  const navigate = useNavigate();

  const handleStart = (): void => {
    if (!loading && vault && status === "success") {
      setState((prevState) => ({ ...prevState, loading: true }));

      addVault(vault)
        .then(() => {
          setState((prevState) => ({ ...prevState, loading: false }));

          navigate(constantPaths.chains);
        })
        .catch(() => {
          setState((prevState) => ({
            ...prevState,
            loading: false,
            status: "error",
          }));
        });
    }
  };

  const handleRemove = (): void => {
    setState(initialState);
  };

  const handleError = (error: string) => {
    setState((prevState) => ({ ...prevState, status: "error" }));

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
  };

  const handleUpload = (file: File): false => {
    setState(initialState);

    const reader = new FileReader();

    const imageFormats: string[] = [
      "image/jpg",
      "image/jpeg",
      "image/png",
      "image/bmp",
    ];

    reader.onload = () => {
      const readerOptions: ReaderOptions = {
        tryHarder: true,
        formats: ["QRCode"],
        maxNumberOfSymbols: 1,
      };

      setState((prevState) => ({
        ...prevState,
        file: {
          data: reader.result?.toString() || "",
          name: file.name,
        },
      }));

      readBarcodesFromImageFile(file, readerOptions)
        .then(([result]) => {
          if (result) {
            try {
              const vault: VaultProps = JSON.parse(result.text);

              setState((prevState) => ({
                ...prevState,
                vault: toCamelCase(vault),
                status: "success",
              }));
            } catch {
              handleError(errorKey.INVALID_VAULT);
            }
          }
        })
        .catch(() => {
          handleError(errorKey.INVALID_QRCODE);
        });
    };

    reader.onerror = () => {
      handleError(errorKey.INVALID_FILE);
    };

    if (imageFormats.indexOf(file.type) >= 0) {
      reader.readAsDataURL(file);
    } else {
      handleError(errorKey.INVALID_EXTENSION);
    }

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
        <h2 className="heading">{t(translation.UPLOAD_VAULT_SHARE)}</h2>
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
              <h3 className="title">{t(translation.UPLOAD_QR_CODE)}</h3>
              <span className="text">
                {t(translation.DROP_FILE_HERE)}
                <u>{t(translation.UPLOAD_IT)}</u>
              </span>
            </>
          )}
        </Upload.Dragger>
        <p className="hint">{t(translation.HINT)}</p>
        <Button
          disabled={status !== "success"}
          loading={loading}
          onClick={handleStart}
          type={status === "success" ? "primary" : "default"}
          block
        >
          {t(translation.START)}
        </Button>
      </div>
      <p className="hint">{t(translation.DOWNLOAD_APP)}</p>
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
