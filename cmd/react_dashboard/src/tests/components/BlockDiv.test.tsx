/**
 * @vitest-environment jsdom
 */
import { render, screen, within } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import BlockDiv from "components/BlockDiv";

const blockWithTransaction: Block = {
  timestamp: 1710000000,
  nonce: 42,
  previousHash: "previous-hash-123",
  transactions: [
    {
      message: "USER TRANSACTION",
      recipientBlockchainAddress: "recipient-blockchain-address-abc123456789",
      senderBlockchainAddress: "sender-blockchain-address-xyz987654321",
      senderPrivateKey: "sender-private-key",
      senderPublicKey: "sender-public-key",
      value: "1.25",
    },
  ],
};

describe("BlockDiv", () => {
  it("renders block metadata and transaction details", () => {
    render(<BlockDiv block={blockWithTransaction} />);

    expect(screen.getByRole("heading", { name: "Block" })).toBeInTheDocument();
    expect(screen.getByText("1710000000")).toBeInTheDocument();
    expect(screen.getByText("42")).toBeInTheDocument();
    expect(screen.getByText("previous-hash-123")).toBeInTheDocument();

    const table = screen.getByRole("table");
    expect(within(table).getByText("Sender")).toBeInTheDocument();
    expect(within(table).getByText("Recipient")).toBeInTheDocument();
    expect(within(table).getByText("USER TRANSACTION")).toBeInTheDocument();
    expect(within(table).getByText("1.25")).toBeInTheDocument();
    expect(within(table).getByText("...ss-xyz987654321")).toBeInTheDocument();
    expect(within(table).getByText("...ss-abc123456789")).toBeInTheDocument();
  });

  it("shows a genesis block message when there are no transactions", () => {
    render(<BlockDiv block={{ ...blockWithTransaction, transactions: [] }} />);

    expect(screen.getByText("No transactions (genesis block).")).toBeInTheDocument();
  });
});
