import { FC, useEffect, useState } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import { List, Modal } from "antd";

import { Language, languageName, storageKey } from "utils/constants";
import i18n from "i18n/config";
import constantModals from "modals/constant-modals";

interface InitialState {
  visible: boolean;
}

const Component: FC = () => {
  const initialState: InitialState = { visible: false };
  const [state, setState] = useState(initialState);
  const { visible } = state;
  const { hash } = useLocation();
  const navigate = useNavigate();

  const handleSelect = (key: string): void => {
    i18n.changeLanguage(key);

    localStorage.setItem(storageKey.LANGUAGE, key);

    navigate(-1);
  };

  const componentDidUpdate = (): void => {
    switch (hash) {
      case `#${constantModals.CHANGE_LANG}`: {
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
      key: Language.ENGLISH,
      title: languageName[Language.ENGLISH],
    },
    {
      key: Language.GERMAN,
      title: languageName[Language.GERMAN],
    },
    {
      key: Language.SPANISH,
      title: languageName[Language.SPANISH],
    },
    {
      key: Language.ITALIAN,
      title: languageName[Language.ITALIAN],
    },
    {
      key: Language.CROATIA,
      title: languageName[Language.CROATIA],
    },
    {
      key: Language.RUSSIAN,
      title: languageName[Language.RUSSIAN],
    },
    {
      key: Language.DUTCH,
      title: languageName[Language.DUTCH],
    },
    {
      key: Language.PORTUGUESE,
      title: languageName[Language.PORTUGUESE],
    },
  ];

  return (
    <Modal
      className="modal-language"
      title="Change Language"
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
            className={i18n.language === key ? "active" : ""}
          >
            {title}
          </List.Item>
        )}
      />
    </Modal>
  );
};

export default Component;
