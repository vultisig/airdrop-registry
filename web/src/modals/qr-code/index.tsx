import { FC, useEffect, useState } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import { Modal, QRCode, Typography } from "antd";

import constantModals from "modals/constant-modals";

interface InitialState {
  visible: boolean;
}

const Component: FC<{ address: string }> = ({ address }) => {
  const initialState: InitialState = { visible: false };
  const [state, setState] = useState(initialState);
  const { visible } = state;
  const { hash } = useLocation();
  const navigate = useNavigate();

  const componentDidUpdate = () => {
    switch (hash) {
      case `#${constantModals.QR_CODE}`: {
        if (address) setState((prevState) => ({ ...prevState, visible: true }));

        break;
      }
      default: {
        setState(initialState);

        break;
      }
    }
  };

  useEffect(componentDidUpdate, [hash]);
  return (
    <Modal
      className="modal-qr-code"
      centered={true}
      footer={false}
      onCancel={() => navigate(-1)}
      open={visible}
      title="Address"
      width={360}
    >
      <Typography.Title level={5}>{address}</Typography.Title>
      <QRCode value={address} size={312} />
    </Modal>
  );
};

export default Component;
