/**
 * ApiClientOptions
 *
 * Configures the API client with a backend base URL.
 */
export type ApiClientOptions = {
  baseUrl?: string;
};

/**
 * ApiRequestOptions
 *
 * Represents supported HTTP request options for frontend API calls.
 */
export type ApiRequestOptions = {
  method?: "GET" | "POST" | "PATCH" | "PUT" | "DELETE";
  body?: unknown;
};

/**
 * ApiClient
 *
 * Centralizes HTTP request behavior for the ApplyBy frontend.
 * Components and pages should call through this boundary instead of using
 * fetch directly.
 */
export class ApiClient {
  private readonly baseUrl: string;

  /**
   * constructor
   *
   * Creates an API client using the provided base URL or a relative URL by default.
   */
  constructor(options: ApiClientOptions = {}) {
    this.baseUrl = options.baseUrl ?? "";
  }

  /**
   * request
   *
   * Sends a JSON request to the backend and returns the typed response body.
   * Non-successful HTTP statuses are converted into thrown errors.
   */
  async request<TResponse>(path: string, options: ApiRequestOptions = {}): Promise<TResponse> {
    const response = await fetch(`${this.baseUrl}${path}`, {
      method: options.method ?? "GET",
      headers: {
        "Content-Type": "application/json"
      },
      body: options.body === undefined ? undefined : JSON.stringify(options.body)
    });

    if (!response.ok) {
      throw new Error(`API request failed with status ${response.status}.`);
    }

    return response.json() as Promise<TResponse>;
  }
}

/**
 * apiClient
 *
 * Provides the default frontend API client configured from Vite environment values.
 */
export const apiClient = new ApiClient({
  baseUrl: import.meta.env.VITE_API_BASE_URL ?? ""
});