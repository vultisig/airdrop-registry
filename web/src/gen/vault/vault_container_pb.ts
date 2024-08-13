// @generated by protoc-gen-es v2.0.0 with parameter "target=ts"
// @generated from file vault_container.proto (package vultisig.vault.v1, syntax proto3)
/* eslint-disable */

import type { GenFile, GenMessage } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, messageDesc } from "@bufbuild/protobuf/codegenv1";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file vault_container.proto.
 */
export const file_vault_container: GenFile = /*@__PURE__*/
  fileDesc("ChV2YXVsdF9jb250YWluZXIucHJvdG8SEXZ1bHRpc2lnLnZhdWx0LnYxIkYKDlZhdWx0Q29udGFpbmVyEg8KB3ZlcnNpb24YASABKAQSDQoFdmF1bHQYAiABKAkSFAoMaXNfZW5jcnlwdGVkGAMgASgIQisKEXZ1bHRpc2lnLnZhdWx0LnYxWhF2dWx0aXNpZy92YXVsdC92MboCAlZTYgZwcm90bzM");

/**
 * @generated from message vultisig.vault.v1.VaultContainer
 */
export type VaultContainer = Message<"vultisig.vault.v1.VaultContainer"> & {
  /**
   * version of data format
   *
   * @generated from field: uint64 version = 1;
   */
  version: bigint;

  /**
   * vault contained the container
   *
   * @generated from field: string vault = 2;
   */
  vault: string;

  /**
   * is vault encrypted with password
   *
   * @generated from field: bool is_encrypted = 3;
   */
  isEncrypted: boolean;
};

/**
 * Describes the message vultisig.vault.v1.VaultContainer.
 * Use `create(VaultContainerSchema)` to create a new message.
 */
export const VaultContainerSchema: GenMessage<VaultContainer> = /*@__PURE__*/
  messageDesc(file_vault_container, 0);
