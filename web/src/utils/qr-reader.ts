import jsQR from "jsqr";

import { toCamelCase } from "utils/case-converter";
import { FileProps, QRCodeProps, VaultProps } from "utils/interfaces";
import { errorKey } from "utils/constants";

const readQRCode = (data: string): Promise<VaultProps> => {
  return new Promise((resolve, reject) => {
    const canvas = document.createElement("canvas");
    const ctx = canvas.getContext("2d");
    const image = new Image();

    image.src = data;

    image.onload = () => {
      canvas.width = image.width;
      canvas.height = image.height;

      ctx?.drawImage(image, 0, 0, image.width, image.height);

      const imageData = ctx?.getImageData(
        0,
        0,
        image.width,
        image.height
      )?.data;

      if (imageData) {
        const qrData = jsQR(imageData, image.width, image.height);

        if (qrData) {
          const vaultData: VaultProps = toCamelCase(JSON.parse(qrData.data));

          resolve(vaultData);
        } else {
          reject();
        }
      } else {
        reject(errorKey.INVALID_QRCODE);
      }
    };

    image.onerror = () => {
      reject(errorKey.INVALID_FILE);
    };
  });
};

const readImage = (file: File): Promise<FileProps> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    const imageFormats: string[] = [
      "image/jpg",
      "image/jpeg",
      "image/png",
      "image/bmp",
    ];

    reader.onload = () => {
      resolve({
        data: (reader.result || "").toString(),
        name: file.name,
      });
    };

    reader.onerror = () => {
      reject(errorKey.INVALID_FILE);
    };

    if (imageFormats.indexOf(file.type) >= 0) {
      reader.readAsDataURL(file);
    } else {
      reject(errorKey.INVALID_EXTENSION);
    }
  });
};

const qrReader = (file: File): Promise<QRCodeProps> => {
  return new Promise((resolve, reject) => {
    readImage(file)
      .then((file) => {
        readQRCode(file.data)
          .then((vault) => {
            resolve({ file, vault });
          })
          .catch((error) => {
            reject({ file, error });
          });
      })
      .catch((error) => {
        reject(error);
      });
  });
};

export default qrReader;
