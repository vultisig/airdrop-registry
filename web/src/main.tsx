import { StrictMode } from "react";
import ReactDOM from "react-dom/client";

import App from "App.tsx";

import "i18n/config";
import "styles/index.scss";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <App />
  </StrictMode>
);
