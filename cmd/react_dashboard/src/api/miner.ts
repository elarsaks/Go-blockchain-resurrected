import { apiClient } from "api/client";

// Fetch latest blocks
function fetchBlockchainData(signal?: AbortSignal): Promise<Block[]> {
  return apiClient
    .get<Block[]>("/miner/blocks", {
      params: { amount: 10 },
      signal,
    })
    .then((response) => response.data);
}

// Fetch miner wallet details
function fetchMinerWalletDetails(
  minerId: string,
  signal?: AbortSignal
): Promise<WalletDetails> {
  return apiClient
    .post<WalletDetailsResponse>("/miner/wallet", null, {
      params: { miner_id: minerId },
      signal,
    })
    .then(({ data }) => {
      const camelCaseResponseData: WalletDetails = {
        blockchainAddress: data.blockchainAddress,
        privateKey: data.privateKey,
        publicKey: data.publicKey,
      };

      return camelCaseResponseData;
    });
}

export { fetchBlockchainData, fetchMinerWalletDetails };
