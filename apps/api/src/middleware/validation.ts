import { Request, Response, NextFunction } from 'express';
import { z } from 'zod';
import { ResponseHelper } from '../utils/response';

export const validate = (schema: z.ZodSchema) => {
  return (req: Request, res: Response, next: NextFunction) => {
    try {
      const validatedData = schema.parse(req.body);
      req.body = validatedData;
      next();
    } catch (error) {
      if (error instanceof z.ZodError) {
        const errors = error.errors.map(err => ({
          field: err.path.join('.'),
          message: err.message
        }));
        return ResponseHelper.validationError(res, errors);
      }
      next(error);
    }
  };
};