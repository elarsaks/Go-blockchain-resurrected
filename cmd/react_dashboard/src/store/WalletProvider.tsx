import { fetchUserWalletDetails, fetchWalletBalance } from "api/wallet";
import { fetchMinerWalletDetails } from "api/miner";
import { isApiRequestCanceled } from "api/client";
import WalletReducer from "store/WalletReducer";
import React, {
  createContext,
  useCallback,
  useEffect,
  useReducer,
  useRef,
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
  const walletRequestIdRef = useRef(0);
  const minerBalanceRequestIdRef = useRef(0);
  const userBalanceRequestIdRef = useRef(0);
  const walletRequestAbortRef = useRef<AbortController | null>(null);
  const minerBalanceAbortRef = useRef<AbortController | null>(null);
  const userBalanceAbortRef = useRef<AbortController | null>(null);

  const startRequest = useCallback((
    abortRef: React.MutableRefObject<AbortController | null>
  ): AbortSignal => {
    abortRef.current?.abort();
    const controller = new AbortController();
    abortRef.current = controller;
    return controller.signal;
  }, []);

  const abortActiveRequests = useCallback(() => {
    walletRequestAbortRef.current?.abort();
    minerBalanceAbortRef.current?.abort();
    userBalanceAbortRef.current?.abort();
  }, []);

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
    const requestId = walletRequestIdRef.current + 1;
    walletRequestIdRef.current = requestId;
    minerBalanceRequestIdRef.current += 1;
    userBalanceRequestIdRef.current += 1;
    const signal = startRequest(walletRequestAbortRef);
    minerBalanceAbortRef.current?.abort();
    userBalanceAbortRef.current?.abort();

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

    Promise.all([
      fetchMinerWalletDetails(minerId, signal),
      fetchUserWalletDetails(signal),
    ])
      .then(([minerDetails, userDetails]) => {
        if (requestId !== walletRequestIdRef.current) return;

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
      .catch((error) => {
        if (isApiRequestCanceled(error)) return;
        if (requestId !== walletRequestIdRef.current) return;

        setWalletSetupError("Failed to initialize wallets");
      });
  }, [setWalletSetupError, startRequest]);

  const selectMiner = useCallback(
    (minerId: string) => {
      const requestId = walletRequestIdRef.current + 1;
      walletRequestIdRef.current = requestId;
      minerBalanceRequestIdRef.current += 1;
      const signal = startRequest(walletRequestAbortRef);
      minerBalanceAbortRef.current?.abort();

      setSelectedMinerId(minerId);

      dispatchMinerWallet({
        type: "SET_WALLET_UTIL",
        payload: {
          isActive: true,
          type: "info",
          message: "Fetching miner wallet details",
        },
      });

      fetchMinerWalletDetails(minerId, signal)
        .then((minerDetails) => {
          if (requestId !== walletRequestIdRef.current) return;

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
        .catch((error) => {
          if (isApiRequestCanceled(error)) return;
          if (requestId !== walletRequestIdRef.current) return;

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
    [startRequest, userWallet.blockchainAddress]
  );

  const getUserWalletWalletBalance = useCallback(() => {
    const requestId = userBalanceRequestIdRef.current + 1;
    userBalanceRequestIdRef.current = requestId;
    const signal = startRequest(userBalanceAbortRef);

    fetchWalletBalance(userWallet.blockchainAddress, signal)
      .then((userBalance) => {
        if (requestId !== userBalanceRequestIdRef.current) return;

        dispatchUserWallet({
          type: "SET_WALLET",
          payload: { balance: String(userBalance) },
        });

        clearLoader("User");
      })
      .catch((error) => {
        if (isApiRequestCanceled(error)) return;

        if (requestId === userBalanceRequestIdRef.current) {
          dispatchUserWallet({
            type: "SET_WALLET_UTIL",
            payload: {
              isActive: true,
              type: "error",
              message: "Failed to fetch user wallet details",
            },
          });
        }
      });
  }, [startRequest, userWallet.blockchainAddress]);

  const getMinerWalletWalletBalance = useCallback(() => {
    const requestId = minerBalanceRequestIdRef.current + 1;
    minerBalanceRequestIdRef.current = requestId;
    const signal = startRequest(minerBalanceAbortRef);

    fetchWalletBalance(minerWallet.blockchainAddress, signal)
      .then((minerBalance) => {
        if (requestId !== minerBalanceRequestIdRef.current) return;

        dispatchMinerWallet({
          type: "SET_WALLET",
          payload: { balance: String(minerBalance) },
        });
        clearLoader("Miner");
      })
      .catch((error) => {
        if (isApiRequestCanceled(error)) return;

        if (requestId === minerBalanceRequestIdRef.current) {
          dispatchMinerWallet({
            type: "SET_WALLET_UTIL",
            payload: {
              isActive: true,
              type: "error",
              message: "Failed to fetch miner wallet details",
            },
          });
        }
      });
  }, [minerWallet.blockchainAddress, startRequest]);

  // Fetch wallet details
  useEffect(() => {
    loadWallets("1");
  }, [loadWallets]);

  useEffect(() => {
    return () => {
      walletRequestIdRef.current += 1;
      minerBalanceRequestIdRef.current += 1;
      userBalanceRequestIdRef.current += 1;
      abortActiveRequests();
    };
  }, [abortActiveRequests]);

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
