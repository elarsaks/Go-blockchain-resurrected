import React from "react";
import { render, screen, waitFor } from "@testing-library/react";
import App from "./App";

jest.mock("components/layout/Background", () => () => (
  <div data-testid="background" />
));

jest.mock("components/layout/Cube", () => () => <div data-testid="cube" />);

jest.mock("api/client", () => ({
  getApiErrorMessage: () => "API error",
}));

jest.mock("api/miner", () => ({
  fetchBlockchainData: () => Promise.resolve([]),
  fetchMinerWalletDetails: () =>
    Promise.resolve({
      blockchainAddress: "miner-address",
      privateKey: "miner-private-key",
      publicKey: "miner-public-key",
    }),
}));

jest.mock("api/wallet", () => ({
  fetchUserWalletDetails: () =>
    Promise.resolve({
      blockchainAddress: "user-address",
      privateKey: "user-private-key",
      publicKey: "user-public-key",
    }),
  fetchWalletBalance: () => Promise.resolve(1),
  transaction: () => Promise.resolve({ status: "success" }),
}));

test("renders initialized wallets", async () => {
  render(<App />);
  expect(screen.getByText(/go blockchain/i)).toBeInTheDocument();

  await waitFor(() => {
    expect(screen.getAllByDisplayValue("miner-address").length).toBeGreaterThan(0);
    expect(screen.getAllByDisplayValue("user-address").length).toBeGreaterThan(0);
  });
});
