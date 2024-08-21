import { ConfigProvider, theme } from "antd";

import VaultContext from "context";
import Routes from "routes";

const App = () => {
  return (
    <ConfigProvider
      theme={{
        algorithm: theme.darkAlgorithm,
        components: {
          Drawer: {
            colorBgElevated: "#02122b",
            paddingLG: 8,
          },
          Modal: {
            colorBgElevated: "#02122b",
            paddingLG: 8,
          },
          Dropdown: {
            paddingBlock: 8,
          },
          List: {
            itemPadding: "12px 0",
          },
        },
        token: {
          colorPrimary: "#33e6bf",
          colorPrimaryText: "#000000",
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
