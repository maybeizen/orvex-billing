import { Request, Response, NextFunction } from 'express';
import { User } from '../db/models/User';
import { AuthenticatedRequest } from '../types/api';
import { ResponseHelper } from '../utils/response';

export const requireAuth = async (req: AuthenticatedRequest, res: Response, next: NextFunction) => {
  try {
    if (!req.session.userId) {
      return ResponseHelper.unauthorized(res, 'Authentication required');
    }

    const user = await User.findById(req.session.userId).select('-password');
    if (!user) {
      req.session.destroy((err) => {
        if (err) console.error('Session destroy error:', err);
      });
      return ResponseHelper.unauthorized(res, 'User not found');
    }

    req.user = user;
    next();
  } catch (error) {
    console.error('Auth middleware error:', error);
    return ResponseHelper.error(res, 'Internal server error');
  }
};

export const requireGuest = (req: Request, res: Response, next: NextFunction) => {
  if (req.session.userId) {
    return ResponseHelper.error(res, 'Already authenticated', 400);
  }
  next();
};

export const requireAdmin = (req: AuthenticatedRequest, res: Response, next: NextFunction) => {
  if (!req.user) {
    return ResponseHelper.unauthorized(res, 'Authentication required');
  }
  
  if (req.user.role !== 'admin') {
    return ResponseHelper.forbidden(res, 'Admin access required');
  }
  
  next();
};