import { FC, useEffect, useState } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import { List, Modal } from "antd";

import { useVaultContext } from "context";
import { Currency, currencyName } from "utils/constants";
import constantModals from "modals/constant-modals";

interface InitialState {
  visible: boolean;
}

const Component: FC = () => {
  const initialState: InitialState = { visible: false };
  const [state, setState] = useState(initialState);
  const { visible } = state;
  const { changeCurrency, currency, vault } = useVaultContext();
  const { hash } = useLocation();
  const navigate = useNavigate();

  const handleSelect = (key: Currency): void => {
    if (vault) {
      changeCurrency(vault, key);

      navigate(-1);
    }
  };

  const componentDidUpdate = (): void => {
    switch (hash) {
      case `#${constantModals.CHANGE_CURRENCY}`: {
        setState((prevState) => ({ ...prevState, visible: true }));

        break;
      }
      default: {
        setState(initialState);

        break;
      }
    }
  };

  useEffect(componentDidUpdate, [hash]);

  const data = [
    {
      key: Currency.USD,
      title: currencyName[Currency.USD],
    },
    {
      key: Currency.AUD,
      title: currencyName[Currency.AUD],
    },
    {
      key: Currency.CAD,
      title: currencyName[Currency.CAD],
    },
    {
      key: Currency.SGD,
      title: currencyName[Currency.SGD],
    },
    {
      key: Currency.EUR,
      title: currencyName[Currency.EUR],
    },
    {
      key: Currency.RUB,
      title: currencyName[Currency.RUB],
    },
    {
      key: Currency.GPB,
      title: currencyName[Currency.GPB],
    },
    {
      key: Currency.JPY,
      title: currencyName[Currency.JPY],
    },
    {
      key: Currency.CNY,
      title: currencyName[Currency.CNY],
    },
    {
      key: Currency.SEK,
      title: currencyName[Currency.SEK],
    },
  ];

  return (
    <Modal
      className="modal-currency"
      title="Change Currency"
      centered={true}
      footer={false}
      onCancel={() => navigate(-1)}
      maskClosable={false}
      open={visible}
      width={360}
    >
      <List
        dataSource={data}
        renderItem={({ key, title }) => (
          <List.Item
            onClick={() => handleSelect(key)}
            className={currency === key ? "active" : ""}
          >
            <span>{title}</span>
            <span>{key}</span>
          </List.Item>
        )}
      />
    </Modal>
  );
};

export default Component;
