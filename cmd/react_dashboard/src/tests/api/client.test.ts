import axios from "axios";
import { describe, expect, it } from "vitest";
import { getApiErrorMessage } from "api/client";

describe("getApiErrorMessage", () => {
  it("uses API response message fields first", () => {
    const error = new axios.AxiosError(
      "Request failed",
      undefined,
      undefined,
      undefined,
      {
        data: { message: "Not enough balance" },
        status: 400,
        statusText: "Bad Request",
        headers: {},
        config: {} as any,
      }
    );

    expect(getApiErrorMessage(error)).toBe("Not enough balance");
  });

  it("falls back to API error fields and generic Error messages", () => {
    const apiError = new axios.AxiosError(
      "Request failed",
      undefined,
      undefined,
      undefined,
      {
        data: { error: "Address not found" },
        status: 404,
        statusText: "Not Found",
        headers: {},
        config: {} as any,
      }
    );

    expect(getApiErrorMessage(apiError)).toBe("Address not found");
    expect(getApiErrorMessage(new Error("Network down"))).toBe("Network down");
    expect(getApiErrorMessage("unknown")).toBe("Unexpected API error");
  });
});
