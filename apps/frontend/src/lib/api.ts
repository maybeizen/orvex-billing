import axios, { AxiosInstance, AxiosResponse } from "axios";

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:3000/api";

export interface ApiResponse<T = any> {
  success?: boolean;
  error?: string;
  message?: string;
  data?: T;
  user?: User;
}

export interface User {
  id: number;
  email: string;
  username: string;
  first_name: string;
  last_name: string;
  bio?: string;
  avatar_url?: string;
  email_verified: boolean;
  two_factor_enabled: boolean;
  two_factor_verified?: boolean;
  role?: string;
  created_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
  totp_code?: string;
}

export interface RegisterRequest {
  email: string;
  username: string;
  password: string;
  first_name: string;
  last_name: string;
}

class ApiClient {
  private client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: API_BASE_URL,
      headers: {
        "Content-Type": "application/json",
      },
      withCredentials: true,
    });

    // Add response interceptor to handle errors consistently
    this.client.interceptors.response.use(
      (response: AxiosResponse) => response,
      (error) => {
        const message =
          error.response?.data?.error || 
          error.response?.data?.message || 
          error.message || 
          "Network error";
        throw new Error(message);
      }
    );
  }

  async login(data: LoginRequest): Promise<ApiResponse> {
    const response = await this.client.post("/v1/auth/login", data);
    return response.data;
  }

  async register(data: RegisterRequest): Promise<ApiResponse> {
    const response = await this.client.post("/v1/auth/register", data);
    return response.data;
  }

  async logout(): Promise<ApiResponse> {
    const response = await this.client.post("/v1/auth/logout");
    return response.data;
  }

  async getProfile(): Promise<ApiResponse<User>> {
    const response = await this.client.get("/v1/user/profile");
    return response.data;
  }

  async getHealth(): Promise<ApiResponse> {
    const response = await this.client.get("/health");
    return response.data;
  }

  async getServices(): Promise<ApiResponse> {
    const response = await this.client.get("/services");
    return response.data;
  }

  async getInvoices(): Promise<ApiResponse> {
    const response = await this.client.get("/services/invoices");
    return response.data;
  }

  // Profile Management
  async updateProfile(data: Partial<User>): Promise<ApiResponse> {
    const response = await this.client.put("/v1/user/profile", data);
    return response.data;
  }

  async changePassword(data: { current_password: string; new_password: string }): Promise<ApiResponse> {
    const response = await this.client.put("/v1/user/password", data);
    return response.data;
  }

  async uploadAvatar(formData: FormData): Promise<ApiResponse> {
    const response = await this.client.post("/v1/user/avatar", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    });
    return response.data;
  }

  async removeAvatar(): Promise<ApiResponse> {
    const response = await this.client.delete("/v1/user/avatar");
    return response.data;
  }

  // Two-Factor Authentication
  async setup2FA(): Promise<ApiResponse> {
    const response = await this.client.post("/v1/user/2fa/setup");
    return response.data;
  }

  async enable2FA(data: { secret: string; token: string; password: string }): Promise<ApiResponse> {
    const response = await this.client.post("/v1/user/2fa/enable", data);
    return response.data;
  }

  async disable2FA(token: string): Promise<ApiResponse> {
    const response = await this.client.post("/v1/user/2fa/disable", { token });
    return response.data;
  }

  async generateBackupCodes(token: string): Promise<ApiResponse> {
    const response = await this.client.post("/v1/user/2fa/backup-codes", { token });
    return response.data;
  }

  async verify2FA(token: string | { token: string }): Promise<ApiResponse> {
    const data = typeof token === 'string' ? { token } : token;
    const response = await this.client.post("/v1/user/2fa/verify", data);
    return response.data;
  }

  async verifyBackupCode(code: string): Promise<ApiResponse> {
    const response = await this.client.post("/v1/user/2fa/backup-verify", { code });
    return response.data;
  }

  // Session Management
  async getUserSessions(): Promise<ApiResponse> {
    const response = await this.client.get("/v1/user/sessions");
    return response.data;
  }

  async revokeSession(sessionId: string): Promise<ApiResponse> {
    const response = await this.client.delete(`/v1/user/sessions/${sessionId}`);
    return response.data;
  }

  async revokeAllOtherSessions(): Promise<ApiResponse> {
    const response = await this.client.delete("/v1/user/sessions/others");
    return response.data;
  }
}

export const api = new ApiClient();
