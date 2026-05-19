/**
 * @vitest-environment jsdom
 */
import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import WalletHead from "components/wallet/WalletHead";

const walletDetails: WalletState = {
  amount: "1",
  balance: "12.50",
  blockchainAddress: "wallet-address",
  privateKey: "private-key",
  publicKey: "public-key",
  recipientAddress: "recipient-address",
};

describe("WalletHead", () => {
  it("renders the user wallet title and balance", () => {
    render(
      <WalletHead
        type="User"
        walletDetails={walletDetails}
        selectedMinerId="1"
        onMinerChange={vi.fn()}
      />
    );

    expect(screen.getByRole("heading", { name: "User Wallet" })).toBeVisible();
    expect(screen.getByText("12.50₿")).toBeVisible();
  });

  it("renders miner selection and reports miner changes", () => {
    const onMinerChange = vi.fn();

    render(
      <WalletHead
        type="Miner"
        walletDetails={walletDetails}
        selectedMinerId="1"
        onMinerChange={onMinerChange}
      />
    );

    fireEvent.change(screen.getByRole("combobox"), {
      target: { value: "3" },
    });

    expect(screen.getByRole("option", { name: "Miner 3" })).toBeVisible();
    expect(onMinerChange).toHaveBeenCalledWith("3");
  });
});
