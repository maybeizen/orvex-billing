import { Response } from 'express';
import { ApiResponse, PaginatedResponse } from '../types/api';

export class ResponseHelper {
  static success<T>(
    res: Response,
    data?: T,
    message: string = 'Success',
    statusCode: number = 200
  ): Response<ApiResponse<T>> {
    return res.status(statusCode).json({
      success: true,
      message,
      data,
    });
  }

  static error(
    res: Response,
    message: string = 'An error occurred',
    statusCode: number = 500,
    errors?: Array<{ field: string; message: string }>
  ): Response<ApiResponse> {
    return res.status(statusCode).json({
      success: false,
      message,
      errors,
    });
  }

  static validationError(
    res: Response,
    errors: Array<{ field: string; message: string }>,
    message: string = 'Validation failed'
  ): Response<ApiResponse> {
    return res.status(400).json({
      success: false,
      message,
      errors,
    });
  }

  static unauthorized(
    res: Response,
    message: string = 'Authentication required'
  ): Response<ApiResponse> {
    return res.status(401).json({
      success: false,
      message,
    });
  }

  static forbidden(
    res: Response,
    message: string = 'Access denied'
  ): Response<ApiResponse> {
    return res.status(403).json({
      success: false,
      message,
    });
  }

  static notFound(
    res: Response,
    message: string = 'Resource not found'
  ): Response<ApiResponse> {
    return res.status(404).json({
      success: false,
      message,
    });
  }

  static conflict(
    res: Response,
    message: string = 'Resource already exists'
  ): Response<ApiResponse> {
    return res.status(409).json({
      success: false,
      message,
    });
  }

  static paginated<T>(
    res: Response,
    data: T[],
    pagination: {
      page: number;
      limit: number;
      total: number;
      totalPages: number;
    },
    message: string = 'Success'
  ): Response<PaginatedResponse<T>> {
    return res.status(200).json({
      success: true,
      message,
      data,
      pagination,
    });
  }
}