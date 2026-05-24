import { afterEach, describe, expect, it, vi } from "vitest";
import { apiClient } from "api/client";
import { fetchUserWalletDetails, fetchWalletBalance, transaction } from "api/wallet";

describe("wallet api", () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("registers and maps a user wallet", async () => {
    const post = vi.spyOn(apiClient, "post").mockResolvedValue({
      data: {
        blockchainAddress: "user-address",
        privateKey: "user-private-key",
        publicKey: "user-public-key",
      },
    });

    await expect(fetchUserWalletDetails()).resolves.toEqual({
      blockchainAddress: "user-address",
      privateKey: "user-private-key",
      publicKey: "user-public-key",
    });
    expect(post).toHaveBeenCalledWith("/user/wallet", null, {
      params: { miner_id: "1" },
      signal: undefined,
    });
  });

  it("fetches balances and raises API balance errors", async () => {
    const get = vi
      .spyOn(apiClient, "get")
      .mockResolvedValueOnce({ data: { balance: 3, error: "" } })
      .mockResolvedValueOnce({ data: { balance: 0, error: "Address missing" } });

    await expect(fetchWalletBalance("wallet-address")).resolves.toBe(3);
    await expect(fetchWalletBalance("missing-address")).rejects.toThrow(
      "Address missing",
    );
    expect(get).toHaveBeenCalledWith("/wallet/balance", {
      params: { blockchainAddress: "wallet-address", miner_id: "1" },
      signal: undefined,
    });
  });

  it("posts transactions unchanged", async () => {
    const payload: Transaction = {
      message: "USER TRANSACTION",
      recipientBlockchainAddress: "recipient-address",
      senderBlockchainAddress: "sender-address",
      senderPrivateKey: "sender-private-key",
      senderPublicKey: "sender-public-key",
      value: "1",
    };
    const post = vi
      .spyOn(apiClient, "post")
      .mockResolvedValue({ data: { accepted: true } });

    await expect(transaction(payload)).resolves.toEqual({ accepted: true });
    expect(post).toHaveBeenCalledWith("/transaction", payload, {
      params: { miner_id: "1" },
    });
  });
});
