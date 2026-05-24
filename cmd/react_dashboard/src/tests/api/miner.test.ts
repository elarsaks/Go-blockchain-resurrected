import { afterEach, describe, expect, it, vi } from "vitest";
import { apiClient } from "api/client";
import { fetchBlockchainData, fetchMinerWalletDetails } from "api/miner";

describe("miner api", () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("fetches the latest ten blocks", async () => {
    const blocks: Block[] = [
      {
        timestamp: 1,
        nonce: 2,
        previousHash: "hash",
        transactions: [],
      },
    ];
    const get = vi.spyOn(apiClient, "get").mockResolvedValue({ data: blocks });

    await expect(fetchBlockchainData()).resolves.toEqual(blocks);
    expect(get).toHaveBeenCalledWith("/miner/blocks", {
      params: { amount: 10, miner_id: "1" },
      signal: undefined,
    });
  });

  it("maps miner wallet responses to wallet details", async () => {
    const post = vi.spyOn(apiClient, "post").mockResolvedValue({
      data: {
        blockchainAddress: "miner-address",
        privateKey: "miner-private-key",
        publicKey: "miner-public-key",
      },
    });

    await expect(fetchMinerWalletDetails("2")).resolves.toEqual({
      blockchainAddress: "miner-address",
      privateKey: "miner-private-key",
      publicKey: "miner-public-key",
    });
    expect(post).toHaveBeenCalledWith("/miner/wallet", null, {
      params: { miner_id: "2" },
      signal: undefined,
    });
  });
});
