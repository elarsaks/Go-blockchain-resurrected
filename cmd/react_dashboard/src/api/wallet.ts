import { apiClient } from "api/client";

function fetchUserWalletDetails(signal?: AbortSignal): Promise<WalletDetails> {
  return apiClient
    .post<WalletDetailsResponse>("/user/wallet", null, { signal })
    .then(({ data }) => {
      const camelCaseResponseData: WalletDetails = {
        blockchainAddress: data.blockchainAddress,
        privateKey: data.privateKey,
        publicKey: data.publicKey,
      };

      return camelCaseResponseData;
    });
}

function fetchWalletBalance(
  blockchainAddress: string,
  signal?: AbortSignal
): Promise<number> {
  return apiClient
    .get<BalanceResponse>("/wallet/balance", {
      params: { blockchainAddress },
      signal,
    })
    .then(({ data }) => {
      if (data.error) {
        throw new Error(data.error);
      }
      return data.balance;
    });
}

function transaction(transaction: Transaction): Promise<any> {
  // Why this string ends up in golang as a number is beyond me
  return apiClient
    .post<any>("/transaction", transaction)
    .then(({ data }) => data);
}

export { fetchUserWalletDetails, fetchWalletBalance, transaction };
