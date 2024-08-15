import jsQR from "jsqr";
import CaseConverter from "utils/case-converter";

namespace QRReader {
  export enum Error {
    INVALID_EXTENSION,
    INVALID_QRCODE,
    INVALID_VAULT,
    INVALID_FILE,
  }

  export interface FileProps {
    data: string;
    name: string;
  }

  export interface VaultProps {
    uid: string;
    name: string;
    publicKeyEcdsa: string;
    publicKeyEddsa: string;
    hexChainCode: string;
  }

  export interface ResultProps {
    file: FileProps;
    vault: VaultProps;
  }

  const imageFormats: string[] = [
    "image/jpg",
    "image/jpeg",
    "image/png",
    "image/bmp",
  ];

  const readFile = (file: File): Promise<FileProps> => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();

      reader.onload = () => {
        resolve({
          data: (reader.result || "").toString(),
          name: file.name,
        });
      };

      reader.onerror = () => {
        reject(Error.INVALID_FILE);
      };

      if (imageFormats.indexOf(file.type) >= 0) {
        reader.readAsDataURL(file);
      } else {
        reject(Error.INVALID_EXTENSION);
      }
    });
  };

  const readImage = (data: string): Promise<VaultProps> => {
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
            const vaultData: VaultProps = CaseConverter.toCamel(
              JSON.parse(qrData.data)
            );

            resolve(vaultData);
          } else {
            reject();
          }
        } else {
          reject(Error.INVALID_QRCODE);
        }
      };

      image.onerror = () => {
        reject(Error.INVALID_FILE);
      };
    });
  };

  export const read = (file: File): Promise<ResultProps> => {
    return new Promise((resolve, reject) => {
      readFile(file)
        .then((file) => {
          readImage(file.data)
            .then((vault) => {
              resolve({file, vault});
            })
            .catch((error) => {
              reject({file, error});
            });
        })
        .catch((error) => {
          reject(error);
        });
    });
  };
}

export default QRReader;
