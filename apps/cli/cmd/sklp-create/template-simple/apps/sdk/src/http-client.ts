import axios, {
  type AxiosInstance,
  type AxiosRequestConfig,
  type AxiosResponse,
} from "axios";

export interface CoreClientOptions {
  /**
   * Base URL of the Core API. Defaults to http://localhost:4100.
   * Override for staging or production.
   */
  baseURL?: string;

  /**
   * Bearer token sent as `Authorization: Bearer <token>`.
   */
  apiKey?: string;

  /**
   * Optional custom axios instance. Use this to inject your own
   * interceptors, retry policy or telemetry. When provided, `baseURL` and
   * `apiKey` are ignored — wire those into your axios instance directly.
   */
  axios?: AxiosInstance;
}

let sharedClient: AxiosInstance | null = null;

/**
 * Configure the shared axios instance used by every generated SDK call.
 * Call this once at boot; safe to re-call (it replaces the previous
 * instance).
 */
export function configureCoreClient(opts: CoreClientOptions): AxiosInstance {
  if (opts.axios) {
    sharedClient = opts.axios;
    return sharedClient;
  }
  const instance = axios.create({
    baseURL: opts.baseURL ?? "http://localhost:4100",
    headers: opts.apiKey
      ? { Authorization: `Bearer ${opts.apiKey}` }
      : undefined,
  });
  sharedClient = instance;
  return instance;
}

function getClient(): AxiosInstance {
  if (!sharedClient) {
    sharedClient = axios.create({ baseURL: "http://localhost:4100" });
  }
  return sharedClient;
}

/**
 * Mutator used by the generated orval client. Every generated function
 * routes through this so we can swap the underlying axios instance at
 * runtime via `configureCoreClient`.
 */
export const coreHttp = async <T>(config: AxiosRequestConfig): Promise<T> => {
  const response: AxiosResponse<T> = await getClient().request<T>(config);
  return response.data;
};

export default coreHttp;
