export interface ProblemDetails {
  type?: string;
  title?: string;
  status: number;
  detail?: string;
  instance?: string;
  code?: string;
  invalid_params?: { name: string; reason: string }[];
  trace_id?: string;
}

export declare class ApiError extends Error {
  status: number;
  type?: string;
  title?: string;
  detail?: string;
  code?: string;
  invalidParams?: { name: string; reason: string }[];
  constructor(problem: ProblemDetails);
}

export declare class CampusApiClient {
  baseUrl: string;
  token: string | null;
  constructor(baseUrl?: string);
  setToken(token: string | null): void;
  request<T = any>(path: string, options?: RequestInit): Promise<T>;
  get<T = any>(path: string, options?: RequestInit): Promise<T>;
  post<T = any>(path: string, body?: any, options?: RequestInit): Promise<T>;
  patch<T = any>(path: string, body?: any, options?: RequestInit): Promise<T>;
  delete<T = any>(path: string, options?: RequestInit): Promise<T>;
}

export declare const apiClient: CampusApiClient;
