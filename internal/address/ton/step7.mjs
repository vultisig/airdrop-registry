import { mnemonicToWalletKey } from "@ton/crypto";
import { WalletContractV4 } from "@ton/ton";

async function main() {
  // open wallet v4 (notice the correct wallet version here)
  const mnemonic = "steel address phone tobacco harsh powder denial differ mix jealous kind immune mobile easily stairs ivory original exercise attitude young luggage exotic fresh cost"
  const key = await mnemonicToWalletKey(mnemonic.split(" "));
  const wallet = WalletContractV4.create({ publicKey: key.publicKey, workchain: 0 });

  console.log(wallet.address);
  console.log(key.publicKey);
  // print wallet address
  console.log(wallet.address.toString({ testOnly: true }));

  // print wallet workchain
  console.log("workchain:", wallet.address.workChain);
}

main();
