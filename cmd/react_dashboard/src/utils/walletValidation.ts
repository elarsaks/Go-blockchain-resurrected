function isValidTransferAmount(amount: string, balance: string): boolean {
  const parsedAmount = Number(amount);
  const parsedBalance = Number(balance);

  return (
    Number.isFinite(parsedAmount) &&
    Number.isFinite(parsedBalance) &&
    parsedAmount > 0 &&
    parsedAmount <= parsedBalance
  );
}

export { isValidTransferAmount };
