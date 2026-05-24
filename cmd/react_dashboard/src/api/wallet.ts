import { apiClient } from "api/client";

function fetchUserWalletDetails(
  minerIdOrSignal: string | AbortSignal = "1",
  signal?: AbortSignal,
): Promise<WalletDetails> {
  const minerId = typeof minerIdOrSignal === "string" ? minerIdOrSignal : "1";
  const requestSignal = typeof minerIdOrSignal === "string" ? signal : minerIdOrSignal;

  return apiClient
    .post<WalletDetailsResponse>("/user/wallet", null, {
      params: { miner_id: minerId },
      signal: requestSignal,
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

function fetchWalletBalance(
  blockchainAddress: string,
  minerIdOrSignal: string | AbortSignal = "1",
  signal?: AbortSignal,
): Promise<number> {
  const minerId = typeof minerIdOrSignal === "string" ? minerIdOrSignal : "1";
  const requestSignal = typeof minerIdOrSignal === "string" ? signal : minerIdOrSignal;

  return apiClient
    .get<BalanceResponse>("/wallet/balance", {
      params: { blockchainAddress, miner_id: minerId },
      signal: requestSignal,
    })
    .then(({ data }) => {
      if (data.error) {
        throw new Error(data.error);
      }
      return data.balance;
    });
}

function transaction(
  transaction: Transaction,
  minerId: string = "1",
): Promise<TransactionResponse> {
  // Why this string ends up in golang as a number is beyond me
  return apiClient
    .post<TransactionResponse>("/transaction", transaction, {
      params: { miner_id: minerId },
    })
    .then(({ data }) => data);
}

export { fetchUserWalletDetails, fetchWalletBalance, transaction };
