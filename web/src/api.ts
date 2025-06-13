
interface RequestOptions extends RequestInit {
  headers?: Record<string, string>;
}

class ApiClient {
  private baseURL: string;
  private timeout: number;

  constructor() {
    this.baseURL = 'http://localhost:8088';
    this.timeout = 5000;
  }

  private async request<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;
    
    const config: RequestOptions = {
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    };

    // timeout
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), this.timeout);

    try {
      const response = await fetch(url, {
        ...config,
        signal: controller.signal,
      });
      
      clearTimeout(timeoutId);
      
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
      }

      if (response.status === 204) {
        return null as T; 
      }
      
      return await response.json();
    } catch (error) {
      clearTimeout(timeoutId);
      console.error('API request failed:', error);
      throw error;
    }
  }

  async get<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint);
  }

  async post<T, U = any>(endpoint: string, data: U): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async put<T, U = any>(endpoint: string, data: U): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async delete<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'DELETE',
    });
  }
}

const apiClient = new ApiClient()

export interface SignUpRequest {
  /** The username to sign up with. */
  username: string;
  /** The password to sign up with. */
  password: string;
}

export interface User {
    id: number;
    username: string;
}

export interface ApiResponse<T> {
  data: T;
  message: string;
  success: boolean;
}

export const authService = {
  signUp: (request: SignUpRequest): Promise<ApiResponse<User>> => 
    apiClient.post<ApiResponse<User>, SignUpRequest>('/v1/user/signup', request),
  login: (request: SignUpRequest): Promise<ApiResponse<User>> => 
    apiClient.post<ApiResponse<User>, SignUpRequest>('/v1/user/login', request),
  logout: (): Promise<ApiResponse<null>> =>
    apiClient.post<ApiResponse<null>>('/v1/user/logout', ""),
}
