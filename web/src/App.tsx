import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { ConfigProvider } from "antd";

import translation from "i18n/constant-keys";

import Routes from "routes";

const App = () => {
  const { t } = useTranslation();
  const initialState = { loaded: false };
  const [state, setState] = useState(initialState);
  const { loaded } = state;

  useEffect(() => {
    setTimeout(() => {
      setState((prevState) => ({ ...prevState, loaded: true }));
    }, 1000);
  }, []);

  return (
    <ConfigProvider
      theme={{
        token: {
          colorPrimary: "#33e6bf",
          fontFamily: "Montserrat",
        },
      }}
    >
      {loaded ? (
        <Routes />
      ) : (
        <div className="splash-screen">
          <img src="/images/logo-radiation.svg" className="logo" alt="Logo" />
          <h1 className="heading">{t(translation.VULTISIG)}</h1>
          <p className="text">{t(translation.SECURE_CRYPTO_VAULT)}</p>
        </div>
      )}
    </ConfigProvider>
  );
};

export default App;
