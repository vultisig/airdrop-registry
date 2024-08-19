import { useTranslation } from "react-i18next";

import translation from "i18n/constant-keys";

const Component = () => {
  const { t } = useTranslation();

  return (
    <div className="splash-screen">
      <img src="/images/logo-radiation.svg" className="logo" alt="Logo" />
      <h1 className="heading">{t(translation.VULTISIG)}</h1>
      <p className="text">{t(translation.SECURE_CRYPTO_VAULT)}</p>
    </div>
  );
};

export default Component;
