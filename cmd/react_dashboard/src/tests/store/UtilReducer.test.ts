import { describe, expect, it } from "vitest";
import utilReducer from "store/UtilReducer";

const inactiveState: UtilState = {
  isActive: false,
  type: "info",
  message: "",
};

describe("UtilReducer", () => {
  it("turns notifications on with the requested message", () => {
    expect(
      utilReducer(inactiveState, {
        type: "ON",
        payload: {
          type: "error",
          message: "Failed to fetch blockchain data",
        },
      })
    ).toEqual({
      isActive: true,
      type: "error",
      message: "Failed to fetch blockchain data",
    });
  });

  it("turns notifications off without losing the last message", () => {
    const activeState: UtilState = {
      isActive: true,
      type: "success",
      message: "Transaction sent",
    };

    expect(
      utilReducer(activeState, {
        type: "OFF",
        payload: null,
      })
    ).toEqual({
      isActive: false,
      type: "success",
      message: "Transaction sent",
    });
  });
});
