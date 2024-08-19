import { ConfigProvider } from "antd";

import VaultContext from "context";
import Routes from "routes";

const App = () => {
  return (
    <ConfigProvider
      theme={{
        token: {
          colorPrimary: "#33e6bf",
          fontFamily: "Montserrat",
        },
      }}
    >
      <VaultContext>
        <Routes />
      </VaultContext>
    </ConfigProvider>
  );
};

export default App;
