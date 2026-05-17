import { apiClient } from "api/client";

// Fetch latest blocks
function fetchBlockchainData(): Promise<Block[]> {
  return apiClient
    .get<Block[]>("/miner/blocks", {
      params: { amount: 10 },
    })
    .then((response) => response.data);
}

// Fetch miner wallet details
function fetchMinerWalletDetails(minerId: string): Promise<WalletDetails> {
  return apiClient
    .post<WalletDetailsResponse>("/miner/wallet", null, {
      params: { miner_id: minerId },
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
