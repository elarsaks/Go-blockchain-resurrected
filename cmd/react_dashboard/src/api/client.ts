import axios from "axios";

const gatewayApiUrl =
  import.meta.env.VITE_GATEWAY_API_URL || "http://localhost:5000";

const apiClient = axios.create({
  baseURL: gatewayApiUrl,
});

function getApiErrorMessage(error: unknown): string {
  if (axios.isAxiosError(error)) {
    const message = error.response?.data?.message || error.response?.data?.error;
    return message || error.message;
  }

  return error instanceof Error ? error.message : "Unexpected API error";
}

function isApiRequestCanceled(error: unknown): boolean {
  return axios.isCancel(error);
}

export { apiClient, getApiErrorMessage, isApiRequestCanceled };
