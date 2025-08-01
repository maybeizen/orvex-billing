import mongoose from 'mongoose';
import { HealthCheckResponse, ServiceStatus } from '../types/health';

export class HealthService {
  private static startTime = Date.now();

  static async getHealthStatus(): Promise<HealthCheckResponse> {
    const timestamp = new Date().toISOString();
    const uptime = Math.floor((Date.now() - this.startTime) / 1000);
    
    // Check services
    const database = await this.checkDatabase();
    const session = await this.checkSession();
    
    // Get system metrics
    const memory = this.getMemoryUsage();
    const cpu = { usage: 0 }; // Basic implementation
    
    // Determine overall status
    const servicesHealthy = [database, session].every(
      service => service.status === 'healthy'
    );
    
    return {
      status: servicesHealthy ? 'ok' : 'error',
      timestamp,
      uptime,
      environment: process.env.NODE_ENV || 'development',
      version: process.env.API_VERSION || '1.0.0',
      services: {
        database,
        session,
      },
      system: {
        memory,
        cpu,
      },
    };
  }

  private static async checkDatabase(): Promise<ServiceStatus> {
    const startTime = Date.now();
    
    try {
      // Simple ping to check database connectivity
      await mongoose.connection.db.admin().ping();
      
      const responseTime = Date.now() - startTime;
      const readyState = mongoose.connection.readyState;
      
      let status: ServiceStatus['status'] = 'healthy';
      let message = 'Database connection is healthy';
      
      if (readyState !== 1) { // 1 = connected
        status = 'unhealthy';
        message = `Database connection state: ${this.getConnectionState(readyState)}`;
      } else if (responseTime > 1000) {
        status = 'degraded';
        message = 'Database response time is slow';
      }
      
      return {
        status,
        responseTime,
        message,
        lastChecked: new Date().toISOString(),
      };
    } catch (error) {
      return {
        status: 'unhealthy',
        responseTime: Date.now() - startTime,
        message: `Database error: ${error instanceof Error ? error.message : 'Unknown error'}`,
        lastChecked: new Date().toISOString(),
      };
    }
  }

  private static async checkSession(): Promise<ServiceStatus> {
    const startTime = Date.now();
    
    try {
      // Check if session store is accessible
      // This is a basic check - in production you might want to test actual session operations
      const isConnected = mongoose.connection.readyState === 1;
      
      if (isConnected) {
        return {
          status: 'healthy',
          responseTime: Date.now() - startTime,
          message: 'Session store is accessible',
          lastChecked: new Date().toISOString(),
        };
      } else {
        return {
          status: 'unhealthy',
          responseTime: Date.now() - startTime,
          message: 'Session store is not accessible',
          lastChecked: new Date().toISOString(),
        };
      }
    } catch (error) {
      return {
        status: 'unhealthy',
        responseTime: Date.now() - startTime,
        message: `Session store error: ${error instanceof Error ? error.message : 'Unknown error'}`,
        lastChecked: new Date().toISOString(),
      };
    }
  }

  private static getMemoryUsage() {
    const memUsage = process.memoryUsage();
    const totalMemory = memUsage.heapTotal;
    const usedMemory = memUsage.heapUsed;
    
    return {
      used: Math.round(usedMemory / 1024 / 1024), // MB
      total: Math.round(totalMemory / 1024 / 1024), // MB
      percentage: Math.round((usedMemory / totalMemory) * 100),
    };
  }

  private static getConnectionState(state: number): string {
    const states = {
      0: 'disconnected',
      1: 'connected',
      2: 'connecting',
      3: 'disconnecting',
    };
    return states[state as keyof typeof states] || 'unknown';
  }
}