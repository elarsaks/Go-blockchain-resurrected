import { fetchUserWalletDetails, fetchWalletBalance } from "api/wallet";
import { fetchMinerWalletDetails } from "api/miner";
import WalletReducer from "store/WalletReducer";
import React, {
  createContext,
  useCallback,
  useEffect,
  useReducer,
  useState,
} from "react";

const initialState: StoreWallet = {
  amount: "",
  balance: "0.00",
  blockchainAddress: "",
  privateKey: "",
  publicKey: "",
  recipientAddress: "",
  util: {
    isActive: false,
    type: "info",
    message: "",
  },
};

export const WalletContext = createContext({
  minerWallet: initialState,
  userWallet: initialState,
  selectedMinerId: "1",
  selectMiner: (minerId: string) => {},
  setUserWallet: (wallet: Partial<StoreWallet>) => {},
  setMinerWallet: (wallet: Partial<StoreWallet>) => {},
  setUserWalletUtil: (util: UtilState) => {},
  setMinerWalletUtil: (util: UtilState) => {},
});

interface WalletProviderProps {
  children: React.ReactNode;
  previousHash?: string;
}

export const WalletProvider: React.FC<WalletProviderProps> = ({
  children,
  previousHash,
}) => {
  const [minerWallet, dispatchMinerWallet] = useReducer(
    WalletReducer,
    initialState
  );
  const [userWallet, dispatchUserWallet] = useReducer(
    WalletReducer,
    initialState
  );
  const [selectedMinerId, setSelectedMinerId] = useState("1");

  function clearLoader(type: string) {
    if (type === "User") {
      dispatchUserWallet({
        type: "SET_WALLET_UTIL",
        payload: {
          isActive: false,
          type: "info",
          message: "",
        },
      });
    } else {
      dispatchMinerWallet({
        type: "SET_WALLET_UTIL",
        payload: {
          isActive: false,
          type: "info",
          message: "",
        },
      });
    }
  }

  const setWalletSetupError = useCallback((message: string) => {
    dispatchUserWallet({
      type: "SET_WALLET_UTIL",
      payload: {
        isActive: true,
        type: "error",
        message,
      },
    });

    dispatchMinerWallet({
      type: "SET_WALLET_UTIL",
      payload: {
        isActive: true,
        type: "error",
        message,
      },
    });
  }, []);

  const loadWallets = useCallback((minerId: string) => {
    dispatchUserWallet({
      type: "SET_WALLET_UTIL",
      payload: {
        isActive: true,
        type: "info",
        message:
          "Registering the user wallet on the blockchain. This process can take up to 28 seconds.",
      },
    });

    dispatchMinerWallet({
      type: "SET_WALLET_UTIL",
      payload: {
        isActive: true,
        type: "info",
        message: "Fetching miner wallet details",
      },
    });

    Promise.all([fetchMinerWalletDetails(minerId), fetchUserWalletDetails()])
      .then(([minerDetails, userDetails]) => {
        dispatchMinerWallet({
          type: "SET_WALLET",
          payload: {
            ...minerDetails,
            amount: "1",
            recipientAddress: userDetails.blockchainAddress,
          },
        });

        dispatchUserWallet({
          type: "SET_WALLET",
          payload: {
            ...userDetails,
            recipientAddress: minerDetails.blockchainAddress,
          },
        });

        clearLoader("Miner");
        clearLoader("User");
      })
      .catch(() => {
        setWalletSetupError("Failed to initialize wallets");
      });
  }, [setWalletSetupError]);

  const selectMiner = useCallback(
    (minerId: string) => {
      setSelectedMinerId(minerId);

      dispatchMinerWallet({
        type: "SET_WALLET_UTIL",
        payload: {
          isActive: true,
          type: "info",
          message: "Fetching miner wallet details",
        },
      });

      fetchMinerWalletDetails(minerId)
        .then((minerDetails) => {
          dispatchMinerWallet({
            type: "SET_WALLET",
            payload: {
              ...minerDetails,
              amount: "1",
              recipientAddress: userWallet.blockchainAddress,
            },
          });

          dispatchUserWallet({
            type: "SET_WALLET",
            payload: {
              recipientAddress: minerDetails.blockchainAddress,
            },
          });

          clearLoader("Miner");
        })
        .catch(() => {
          dispatchMinerWallet({
            type: "SET_WALLET_UTIL",
            payload: {
              isActive: true,
              type: "error",
              message: "Failed to fetch miner wallet details",
            },
          });
        });
    },
    [userWallet.blockchainAddress]
  );

  const getUserWalletWalletBalance = useCallback(() => {
    fetchWalletBalance(userWallet.blockchainAddress)
      .then((userBalance) => {
        dispatchUserWallet({
          type: "SET_WALLET",
          payload: { balance: String(userBalance) },
        });

        clearLoader("User");
      })
      .catch((error) =>
        dispatchUserWallet({
          type: "SET_WALLET_UTIL",
          payload: {
            isActive: true,
            type: "error",
            message: "Failed to fetch user wallet details",
          },
        })
      );
  }, [userWallet.blockchainAddress]);

  const getMinerWalletWalletBalance = useCallback(() => {
    fetchWalletBalance(minerWallet.blockchainAddress)
      .then((minerBalance) => {
        dispatchMinerWallet({
          type: "SET_WALLET",
          payload: { balance: String(minerBalance) },
        });
        clearLoader("Miner");
      })
      .catch((error) =>
        dispatchMinerWallet({
          type: "SET_WALLET_UTIL",
          payload: {
            isActive: true,
            type: "error",
            message: "Failed to fetch miner wallet details", // Fixed the message to refer to miner instead of user
          },
        })
      );
  }, [minerWallet.blockchainAddress]);

  // Fetch wallet details
  useEffect(() => {
    loadWallets("1");
  }, [loadWallets]);

  // Fetch wallet balance
  useEffect(() => {
    if (minerWallet.blockchainAddress) getMinerWalletWalletBalance();
    if (userWallet.blockchainAddress) getUserWalletWalletBalance();
  }, [
    minerWallet.blockchainAddress,
    userWallet.blockchainAddress,
    getMinerWalletWalletBalance,
    getUserWalletWalletBalance,
    previousHash,
  ]);

  return (
    <WalletContext.Provider
      value={{
        minerWallet,
        userWallet,
        selectedMinerId,
        selectMiner,
        setUserWallet: (wallet: Partial<StoreWallet>) =>
          dispatchUserWallet({ type: "SET_WALLET", payload: wallet }),
        setMinerWallet: (wallet: Partial<StoreWallet>) =>
          dispatchMinerWallet({ type: "SET_WALLET", payload: wallet }),
        setUserWalletUtil: (util: UtilState) =>
          dispatchUserWallet({ type: "SET_WALLET_UTIL", payload: util }),
        setMinerWalletUtil: (util: UtilState) =>
          dispatchMinerWallet({ type: "SET_WALLET_UTIL", payload: util }),
      }}
    >
      {children}
    </WalletContext.Provider>
  );
};

export default WalletProvider;
