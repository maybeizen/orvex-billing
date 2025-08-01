import { Request, Response } from 'express';
import type { HealthCheckResponse } from '../../types/api';
import { HealthService } from '../../utils/health';
import { ResponseHelper } from '../../utils/response';

export const healthCheck = async (req: Request, res: Response) => {
  try {
    const healthStatus = await HealthService.getHealthStatus();

    const statusCode = healthStatus.status === 'ok' ? 200 : 503;

    res.status(statusCode).json({
      success: healthStatus.status === 'ok',
      message: healthStatus.status === 'ok' ? 'System is healthy' : 'System has issues',
      data: healthStatus,
    });
  } catch (error) {
    console.error('Health check error:', error);
    ResponseHelper.error(res, 'Health check failed', 503);
  }
};

export const liveness = async (req: Request, res: Response) => {
  ResponseHelper.success(res, { alive: true }, 'Service is alive');
};

export const readiness = async (req: Request, res: Response) => {
  try {
    const healthStatus = await HealthService.getHealthStatus();
    const isReady = healthStatus.services.database.status === 'healthy';

    if (isReady) {
      ResponseHelper.success(res, { ready: true }, 'Service is ready');
    } else {
      ResponseHelper.error(res, 'Service is not ready', 503);
    }
  } catch (error) {
    console.error('Readiness check error:', error);
    ResponseHelper.error(res, 'Readiness check failed', 503);
  }
};
