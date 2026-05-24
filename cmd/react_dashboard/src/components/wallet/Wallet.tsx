import { getApiErrorMessage } from "api/client";
import { transaction } from "api/wallet";
import Notification from "components/shared/Notification";
import React, { useEffect, useState, useContext } from "react";
import styled from "styled-components";
import WalletHead from "./WalletHead";
import { WalletContext } from "store/WalletProvider";
import { isValidTransferAmount } from "utils/walletValidation";

interface WalletContainerProps {
  isMiner: boolean;
}

const WalletContainer = styled.div<WalletContainerProps>`
  background-color: #f2f2f2;
  padding: 1rem;
  margin: 1rem;
  border: 1px solid #00acd7;
  border-radius: 8px;
  width: 350px;

  @media (min-width: 850px) {
    margin-left: ${(props) => (props.isMiner ? "0" : "2rem")};
    margin-right: ${(props) => (props.isMiner ? "2rem" : "0")};
  }

  @media (max-width: 850px) {
    width: 80vw;
  }
`;

const Form = styled.div`
  margin-top: 1rem;
`;

const Field = styled.div`
  margin-bottom: 1rem;
`;

const Label = styled.label`
  display: block;
  margin-bottom: 0.5rem;
  font-weight: bold;
  text-align: left;
`;

const TextArea = styled.textarea`
  width: 95%;
  padding: 0.5rem;
  text-align: left;

  &:read-only {
    background-color: #e9ecef;
    color: #495057;
    cursor: default;
  }
`;

const Input = styled.input`
  width: 95%;
  padding: 0.5rem;
  text-align: left;
`;

interface ButtonProps {
  disabled: boolean;
}
const SendButton = styled.button<ButtonProps>`
  margin-top: 1rem;
  padding: 0.75rem 1.5rem;
  background-color: ${(props) => (props.disabled ? "#ccc" : "#00acd7")};
  color: white;
  border: none;
  border-radius: 3px;
  font-weight: bold;
  cursor: ${(props) => (props.disabled ? "not-allowed" : "pointer")};
  float: right;
  opacity: ${(props) => (props.disabled ? "0.6" : "1")};
`;

type WalletType = "Miner" | "User";

type WalletProps = {
  type: WalletType;
};

const Wallet: React.FC<WalletProps> = ({ type }) => {
  const [isSendDisabled, setIsSendDisabled] = useState(false);

  const walletContext = useContext(WalletContext);

  const wallet =
    type === "Miner"
      ? {
          details: walletContext.minerWallet,
          setDetails: walletContext?.setMinerWallet,
          setUtil: walletContext?.setMinerWalletUtil,
        }
      : {
          details: walletContext.userWallet,
          setDetails: walletContext?.setUserWallet,
          setUtil: walletContext?.setUserWalletUtil,
        };

  useEffect(() => {
    const isMissingRequiredField =
      wallet.details.blockchainAddress === "" ||
      wallet.details.privateKey === "" ||
      wallet.details.publicKey === "" ||
      wallet.details.recipientAddress === "";

    setIsSendDisabled(
      isMissingRequiredField ||
        !isValidTransferAmount(wallet.details.amount, wallet.details.balance) ||
        wallet.details.util.isActive,
    );
  }, [wallet.details]);

  // Event Handlers
  const handleInputChange = (
    event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
  ) => {
    const { name, value } = event.target;

    wallet.setDetails({
      ...wallet.details,
      [name]: value,
    });
  };

  const sendCrypto = () => {
    wallet.setUtil({
      isActive: true,
      type: "info",
      message: "Transaction request sent",
    });

    transaction(
      {
        message: "USER TRANSACTION",
        recipientBlockchainAddress: wallet.details.recipientAddress,
        senderBlockchainAddress: wallet.details.blockchainAddress,
        senderPrivateKey: wallet.details.privateKey,
        senderPublicKey: wallet.details.publicKey,
        value: wallet.details.amount,
      },
      walletContext.selectedMinerId,
    )
      .then(() => {
        walletContext.setMinerWalletUtil({
          isActive: true,
          type: "success",
          message:
            "The balance will be updated once the next block is mined. This process can take up to 28 seconds.",
        });

        walletContext.setUserWalletUtil({
          isActive: true,
          type: "success",
          message:
            "The balance will be updated once the next block is mined. This process can take up to 28 seconds.",
        });
      })
      .catch((error) => {
        const message = getApiErrorMessage(error);

        wallet.setUtil({
          isActive: true,
          type: "error",
          message,
        });
      });
  };

  return (
    <WalletContainer isMiner={type === "Miner"}>
      <WalletHead
        type={type}
        walletDetails={wallet.details}
        selectedMinerId={walletContext.selectedMinerId}
        onMinerChange={walletContext.selectMiner}
      />

      <Form>
        <Field>
          <Label>Public Key</Label>
          <TextArea
            rows={4}
            name="publicKey"
            value={wallet.details.publicKey}
            readOnly
          />
        </Field>

        <Field>
          <Label>Private Key</Label>
          <TextArea
            rows={2}
            name="privateKey"
            value={wallet.details.privateKey}
            readOnly
          />
        </Field>

        <Field>
          <Label>{type === "Miner" ? "Miner " : "User"} Blockchain Address </Label>
          <TextArea
            rows={2}
            name="blockchainAddress"
            value={wallet.details.blockchainAddress}
            readOnly
          />
        </Field>

        <Field>
          <Label>Recipient Blockchain Address</Label>
          <TextArea
            rows={2}
            name="recipientAddress"
            placeholder={
              type === "Miner" ? "User Blockchain Address" : "Miner Blockchain Address"
            }
            value={wallet.details.recipientAddress}
            onChange={handleInputChange}
          />
        </Field>

        <Field>
          <Label>Amount: </Label>
          <Input
            type="number"
            name="amount"
            placeholder="0.00₿"
            value={wallet.details.amount.toString()}
            onChange={handleInputChange}
            max={wallet.details.balance}
            min="0.00000001"
            step="any"
          />
        </Field>

        {wallet.details.util.isActive && (
          <Notification
            type={wallet.details.util.type}
            message={wallet.details.util.message}
            insideContainer={true}
          />
        )}

        <SendButton type="submit" disabled={isSendDisabled} onClick={sendCrypto}>
          Send crypto
        </SendButton>
      </Form>
    </WalletContainer>
  );
};

export default Wallet;
