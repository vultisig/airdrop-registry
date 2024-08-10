import { Buffer } from "buffer";
import { fromBinary } from "@bufbuild/protobuf";
import crypto from "crypto-js";
import jsQR from "jsqr";

import { VaultContainerSchema } from "gen/vault/vault_container_pb";
import { VaultSchema } from "gen/vault/vault_pb";

import type { VaultContainer } from "gen/vault/vault_container_pb";
import type { Vault } from "gen/vault/vault_pb";

namespace VaultDecryptor {
  export enum Error {
    ENCODING,
    CONTAINER,
    PASSWD_REQUIRED,
    INVALID_EXTENSION,
    INVALID_FILE,
    INVALID_PASSWD,
    VAULT,
  }

  export enum Type {
    IMAGE,
    DATA,
  }

  export interface FileProps {
    data: string;
    name: string;
    type: number;
  }

  export interface VaultProps {
    container: VaultContainer;
    vault: Vault;
  }

  const imageFormats: string[] = [
    "image/jpg",
    "image/jpeg",
    "image/png",
    "image/bmp",
    "application/pdf",
  ];

  const fileFormats: string[] = [];

  export const validateBase64 = (str: string): boolean => {
    const regex =
      /^(?:[A-Za-z0-9+/]{4})*?(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$/;

    return regex.test(str);
  };

  export const decodeData = (data: string): Buffer => {
    return Buffer.from(data, "base64");
  };

  export const decodeContainer = (bytes: Buffer): VaultContainer => {
    return fromBinary(VaultContainerSchema, bytes);
  };

  export const decodeVault = (bytes: Buffer): Vault => {
    return fromBinary(VaultSchema, bytes);
  };

  export const decryptVault = (bytes: Buffer, passwd: string): Buffer => {
    const key = crypto.SHA256(passwd);

    const decrypted = crypto.AES.decrypt(bytes.toString(), key.toString());

    return Buffer.from(decrypted.toString(crypto.enc.Utf8));
  };

  export const readFile = (file: File): Promise<FileProps> => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      let type: number;

      reader.onload = () => {
        resolve({
          data: (reader.result || "").toString(),
          name: file.name,
          type,
        });
      };

      reader.onerror = () => {
        reject(Error.INVALID_FILE);
      };

      if (imageFormats.indexOf(file.type) >= 0) {
        type = Type.IMAGE;

        reader.readAsDataURL(file);
      } else if (fileFormats.indexOf(file.type) >= 0) {
        type = Type.DATA;

        reader.readAsText(file);
      } else {
        reject(Error.INVALID_EXTENSION);
      }
    });
  };

  export const readImage = (data: string) => {
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

        if (qrData) console.table(JSON.parse(qrData.data));
      }
    };
  };

  export const decryptor = (
    data: string,
    getPasswd: () => Promise<string>
  ): Promise<VaultProps> => {
    return new Promise((resolve, reject) => {
      if (validateBase64(data)) {
        const _decodedContainer = decodeData(data);
        const decodedContainer = decodeContainer(_decodedContainer);

        if (decodedContainer.vault) {
          const vaultData = decodeData(decodedContainer.vault);

          if (decodedContainer.isEncrypted) {
            getPasswd()
              .then((passwd) => {
                const decryptedVault = decryptVault(vaultData, passwd);

                const decodedVault = decodeVault(decryptedVault);

                if (decodedVault.hexChainCode) {
                  resolve({
                    container: decodedContainer,
                    vault: decodedVault,
                  });
                } else {
                  reject(Error.VAULT);
                }
              })
              .catch(() => {
                reject(Error.PASSWD_REQUIRED);
              });
          } else {
            const decodedVault = decodeVault(vaultData);

            if (decodedVault.hexChainCode) {
              resolve({
                container: decodedContainer,
                vault: decodedVault,
              });
            } else {
              reject(Error.VAULT);
            }
          }
        } else {
          reject(Error.CONTAINER);
        }
      } else {
        reject(Error.ENCODING);
      }
    });
  };
}

//console.log(jsQR);

export default VaultDecryptor;
