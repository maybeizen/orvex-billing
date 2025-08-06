import { Request, Response } from "express";
import { IUserDocument } from "../db/models/User";

export interface ApiResponse<T = any> {
  success: boolean;
  message: string;
  data?: T;
  errors?: Array<{
    field: string;
    message: string;
  }>;
}

export interface PaginatedResponse<T> extends ApiResponse<T[]> {
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

export interface ApiError {
  status: number;
  message: string;
  code?: string;
  details?: any;
}

export interface AuthenticatedRequest<T = any> extends Request<any, any, T> {
  user: IUserDocument;
}

export interface TypedRequest<TParams = any, TBody = any, TQuery = any>
  extends Request<TParams, any, TBody, TQuery> {}

export interface AuthenticatedTypedRequest<
  TParams = any,
  TBody = any,
  TQuery = any
> extends Request<TParams, any, TBody, TQuery> {
  user: IUserDocument;
}

export type TypedResponse<T = any> = Response<ApiResponse<T>>;

export interface AuthUserResponse {
  id: string;
  uuid: string;
  firstName: string;
  lastName: string;
  username: string;
  email: string;
  isEmailVerified: boolean;
  createdAt: Date;
  updatedAt: Date;
}

export interface LoginResponse extends AuthUserResponse {
  sessionId?: string;
}

export interface HealthCheckResponse {
  status: "ok" | "error";
  timestamp: string;
  uptime: number;
  environment: string;
  version: string;
  services: {
    database: ServiceStatus;
    redis?: ServiceStatus;
    session: ServiceStatus;
  };
  system: {
    memory: {
      used: number;
      total: number;
      percentage: number;
    };
    cpu: {
      usage: number;
    };
  };
}

export interface ServiceStatus {
  status: "healthy" | "unhealthy" | "degraded";
  responseTime?: number;
  message?: string;
  lastChecked: string;
}

export type RequestHandler<
  TParams = any,
  TBody = any,
  TQuery = any,
  TResponse = any
> = (
  req: TypedRequest<TParams, TBody, TQuery>,
  res: TypedResponse<TResponse>,
  next: any
) => Promise<void> | void;

export type AuthenticatedRequestHandler<
  TParams = any,
  TBody = any,
  TQuery = any,
  TResponse = any
> = (
  req: AuthenticatedTypedRequest<TParams, TBody, TQuery>,
  res: TypedResponse<TResponse>,
  next: any
) => Promise<void> | void;
