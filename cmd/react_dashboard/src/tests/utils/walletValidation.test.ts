import { describe, expect, it } from "vitest";
import { isValidTransferAmount } from "utils/walletValidation";

describe("isValidTransferAmount", () => {
  it("accepts a positive amount within the available balance", () => {
    expect(isValidTransferAmount("1", "2")).toBe(true);
    expect(isValidTransferAmount("0.25", "0.25")).toBe(true);
  });

  it("rejects empty, zero, negative, and non-numeric amounts", () => {
    expect(isValidTransferAmount("", "2")).toBe(false);
    expect(isValidTransferAmount("0", "2")).toBe(false);
    expect(isValidTransferAmount("-1", "2")).toBe(false);
    expect(isValidTransferAmount("not-a-number", "2")).toBe(false);
  });

  it("rejects amounts above balance or invalid balances", () => {
    expect(isValidTransferAmount("3", "2")).toBe(false);
    expect(isValidTransferAmount("1", "")).toBe(false);
    expect(isValidTransferAmount("1", "not-a-number")).toBe(false);
  });
});
