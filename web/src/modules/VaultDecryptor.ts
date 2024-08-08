import { Buffer } from "buffer";
import { fromBinary } from "@bufbuild/protobuf";
import crypto from "crypto-js";

import { VaultContainerSchema } from "gen/vault/vault_container_pb";
import { VaultSchema } from "gen/vault/vault_pb";

import type { VaultContainer } from "gen/vault/vault_container_pb";
import type { Vault } from "gen/vault/vault_pb";

class VaultDecryptor {
  static extension = ".bak";
  static encoding = "base64";

  static validateExtension = (str: string): boolean => {
    return str.endsWith(VaultDecryptor.extension);
  };

  static validateBase64 = (str: string): boolean => {
    const regex =
      /^(?:[A-Za-z0-9+/]{4})*?(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$/;

    return regex.test(str);
  };

  static decodeData = (data: string): Buffer => {
    return Buffer.from(data, VaultDecryptor.encoding);
  };

  static decodeContainer = (bytes: Buffer): VaultContainer => {
    return fromBinary(VaultContainerSchema, bytes);
  };

  static decodeVault = (bytes: Buffer): Vault => {
    return fromBinary(VaultSchema, bytes);
  };

  static decryptVault = (data: Buffer, passwd: string): Buffer => {
    const key = crypto.SHA256(passwd);

    const decrypted = crypto.AES.decrypt(data.toString(), key.toString());

    return Buffer.from(decrypted.toString(crypto.enc.Utf8));
  };

  static readFile = (file: File): Promise<string> => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();

      reader.onload = () => {
        resolve((reader.result || "").toString());
      };

      reader.onerror = (error) => {
        reject(error);
      };

      reader.readAsText(file);
    });
  };
}

export default VaultDecryptor;
