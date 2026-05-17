import { describe, expect, it } from "vitest";
import walletReducer from "store/WalletReducer";

const baseWallet: StoreWallet = {
  amount: "",
  balance: "0.00",
  blockchainAddress: "",
  privateKey: "",
  publicKey: "",
  recipientAddress: "",
  util: {
    isActive: false,
    type: "info",
    message: "",
  },
};

describe("WalletReducer", () => {
  it("merges wallet updates without dropping existing fields", () => {
    const nextState = walletReducer(baseWallet, {
      type: "SET_WALLET",
      payload: {
        amount: "1",
        blockchainAddress: "sender-address",
      },
    });

    expect(nextState).toEqual({
      ...baseWallet,
      amount: "1",
      blockchainAddress: "sender-address",
    });
  });

  it("replaces util state independently from wallet details", () => {
    const nextState = walletReducer(baseWallet, {
      type: "SET_WALLET_UTIL",
      payload: {
        isActive: true,
        type: "error",
        message: "Transaction failed",
      },
    });

    expect(nextState).toEqual({
      ...baseWallet,
      util: {
        isActive: true,
        type: "error",
        message: "Transaction failed",
      },
    });
  });
});
