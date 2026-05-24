/**
 * @vitest-environment jsdom
 */
import { render, screen, waitFor } from "@testing-library/react";
import React, { StrictMode } from "react";
import { describe, expect, it, vi } from "vitest";
import { WalletContext, WalletProvider } from "store/WalletProvider";

const apiMocks = vi.hoisted(() => ({
  fetchMinerWalletDetails: vi.fn(),
  fetchUserWalletDetails: vi.fn(),
  fetchWalletBalance: vi.fn(),
}));

vi.mock("api/miner", () => ({
  fetchMinerWalletDetails: apiMocks.fetchMinerWalletDetails,
}));

vi.mock("api/wallet", () => ({
  fetchUserWalletDetails: apiMocks.fetchUserWalletDetails,
  fetchWalletBalance: apiMocks.fetchWalletBalance,
}));

vi.mock("api/client", () => ({
  isApiRequestCanceled: (error: unknown) =>
    error instanceof Error && error.name === "CanceledError",
}));

function abortableWallet(wallet: WalletDetails, signal?: AbortSignal) {
  return new Promise<WalletDetails>((resolve, reject) => {
    const timeoutId = window.setTimeout(() => resolve(wallet), 0);

    signal?.addEventListener("abort", () => {
      window.clearTimeout(timeoutId);
      const error = new Error("canceled");
      error.name = "CanceledError";
      reject(error);
    });
  });
}

function WalletStateProbe() {
  return (
    <WalletContext.Consumer>
      {({ minerWallet, userWallet }) => (
        <>
          <span data-testid="miner-address">{minerWallet.blockchainAddress}</span>
          <span data-testid="miner-loader">{String(minerWallet.util.isActive)}</span>
          <span data-testid="user-address">{userWallet.blockchainAddress}</span>
          <span data-testid="user-loader">{String(userWallet.util.isActive)}</span>
        </>
      )}
    </WalletContext.Consumer>
  );
}

describe("WalletProvider", () => {
  it("reloads wallets after StrictMode aborts the first effect pass", async () => {
    apiMocks.fetchMinerWalletDetails.mockImplementation(
      (_minerId: string, signal?: AbortSignal) =>
        abortableWallet(
          {
            blockchainAddress: "miner-address",
            privateKey: "miner-private-key",
            publicKey: "miner-public-key",
          },
          signal,
        ),
    );
    apiMocks.fetchUserWalletDetails.mockImplementation(
      (_minerId: string, signal?: AbortSignal) =>
        abortableWallet(
          {
            blockchainAddress: "user-address",
            privateKey: "user-private-key",
            publicKey: "user-public-key",
          },
          signal,
        ),
    );
    apiMocks.fetchWalletBalance.mockResolvedValue(0);

    render(
      <StrictMode>
        <WalletProvider selectedMinerId="1" onMinerSelect={vi.fn()} previousHash="hash">
          <WalletStateProbe />
        </WalletProvider>
      </StrictMode>,
    );

    await waitFor(() => {
      expect(screen.getByTestId("miner-address")).toHaveTextContent("miner-address");
      expect(screen.getByTestId("user-address")).toHaveTextContent("user-address");
    });
    expect(screen.getByTestId("miner-loader")).toHaveTextContent("false");
    expect(screen.getByTestId("user-loader")).toHaveTextContent("false");
    expect(apiMocks.fetchMinerWalletDetails).toHaveBeenCalledTimes(2);
    expect(apiMocks.fetchUserWalletDetails).toHaveBeenCalledTimes(2);
  });
});
