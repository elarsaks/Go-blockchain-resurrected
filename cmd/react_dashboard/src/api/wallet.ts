import { apiClient } from "api/client";

function fetchUserWalletDetails(): Promise<WalletDetails> {
  return apiClient
    .post<WalletDetailsResponse>("/user/wallet")
    .then(({ data }) => {
      const camelCaseResponseData: WalletDetails = {
        blockchainAddress: data.blockchainAddress,
        privateKey: data.privateKey,
        publicKey: data.publicKey,
      };

      return camelCaseResponseData;
    });
}

function fetchWalletBalance(blockchainAddress: string): Promise<number> {
  return apiClient
    .get<BalanceResponse>(
      `/wallet/balance?blockchainAddress=${blockchainAddress}`
    )
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
