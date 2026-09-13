/**
 * Base API Client with RFC 7807 Problem Details handling
 */
export class ApiError extends Error {
  constructor(problem) {
    super(problem.detail || problem.title || 'API Error');
    this.name = 'ApiError';
    this.status = problem.status;
    this.type = problem.type;
    this.title = problem.title;
    this.detail = problem.detail;
    this.code = problem.code;
    this.invalidParams = problem.invalid_params;
  }
}

export class CampusApiClient {
  constructor(baseUrl = '/api/v1') {
    this.baseUrl = baseUrl;
    this.token = null;
  }

  setToken(token) {
    this.token = token;
  }

  async request(path, options = {}) {
    const headers = {
      'Content-Type': 'application/json',
      Accept: 'application/json, application/problem+json',
      ...options.headers,
    };

    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`;
    }

    const response = await fetch(`${this.baseUrl}${path}`, {
      ...options,
      headers,
    });

    if (response.status === 204) {
      return null;
    }

    const contentType = response.headers.get('content-type') || '';
    const isJson = contentType.includes('application/json') || contentType.includes('application/problem+json');
    const data = isJson ? await response.json() : null;

    if (!response.ok) {
      if (data && (data.type || data.status || data.title)) {
        throw new ApiError(data);
      }
      throw new ApiError({
        status: response.status,
        title: response.statusText,
        detail: 'An unexpected HTTP error occurred.',
      });
    }

    return data;
  }

  get(path, options) {
    return this.request(path, { ...options, method: 'GET' });
  }

  post(path, body, options) {
    return this.request(path, {
      ...options,
      method: 'POST',
      body: body ? JSON.stringify(body) : undefined,
    });
  }

  patch(path, body, options) {
    return this.request(path, {
      ...options,
      method: 'PATCH',
      body: body ? JSON.stringify(body) : undefined,
    });
  }

  delete(path, options) {
    return this.request(path, { ...options, method: 'DELETE' });
  }
}

export const apiClient = new CampusApiClient();
