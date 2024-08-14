import { StrictMode } from "react";
import { Provider } from "react-redux";
import ReactDOM from "react-dom/client";

import store from "utils/store";

import App from "App.tsx";

import "i18n/config";
import "styles/index.scss";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <Provider store={store}>
      <App />
    </Provider>
  </StrictMode>
);
