import { StrictMode } from "react";
import ReactDOM from "react-dom/client";
import "core-js";

import App from "App.tsx";

import "i18n/config";
import "utils/prototypes";
import "styles/index.scss";

import.meta.env.DEV
  ? ReactDOM.createRoot(document.getElementById("root")!).render(<App />)
  : ReactDOM.createRoot(document.getElementById("root")!).render(
      <StrictMode>
        <App />
      </StrictMode>
    );
