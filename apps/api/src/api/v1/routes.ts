import express, { Router } from "express";
import { requireAuth, requireGuest, requireAdmin } from "../../middleware/auth";
import {
  forgotPasswordSchema,
  resetPasswordSchema,
  loginSchema,
  registerSchema,
} from "../../validators/auth";
import {
  createUserSchema,
  updateUserSchema,
  updateUserPasswordSchema,
  getUsersQuerySchema,
  userParamsSchema
} from "../../validators/users";
import { validate } from "../../middleware/validation";
import { ResponseHelper } from "../../utils/response";
import {
  loginRateLimiter,
  registerRateLimiter,
  passwordResetRateLimiter,
} from "../../middleware/rate-limit";
import { forgotPassword, resetPassword } from "./auth/forgot-password";
import { getMe } from "./auth/me";
import { logout } from "./auth/logout";
import { register } from "./auth/register";
import { login } from "./auth/login";
import { healthCheck, liveness, readiness } from "./health";
import { getUsers } from "./admin/users/get-users";
import { getUserById } from "./admin/users/get-user";
import { createUser } from "./admin/users/create-user";
import { updateUser } from "./admin/users/update-user";
import { updateUserPassword } from "./admin/users/update-user-password";
import { deleteUser } from "./admin/users/delete-user";

const router: Router = Router();
const authRouter: Router = Router();

const validateQuery = (schema: any) => {
  return (req: any, res: any, next: any) => {
    try {
      schema.parse(req.query || {});
      next();
    } catch (error) {
      if (error instanceof Error && 'issues' in error) {
        const zodError = error as any;
        const errors = zodError.issues.map((err: any) => ({
          field: err.path.join("."),
          message: err.message,
        }));
        return ResponseHelper.validationError(res, errors);
      }
      next(error);
    }
  };
};

const validateParams = (schema: any) => {
  return (req: any, res: any, next: any) => {
    try {
      schema.parse(req.params || {});
      next();
    } catch (error) {
      if (error instanceof Error && 'issues' in error) {
        const zodError = error as any;
        const errors = zodError.issues.map((err: any) => ({
          field: err.path.join("."),
          message: err.message,
        }));
        return ResponseHelper.validationError(res, errors);
      }
      next(error);
    }
  };
};

router.get("/health", healthCheck);
router.get("/health/live", liveness);
router.get("/health/ready", readiness);

authRouter.post(
  "/register",
  registerRateLimiter,
  requireGuest,
  validate(registerSchema),
  register
);
authRouter.post(
  "/login",
  loginRateLimiter,
  requireGuest,
  validate(loginSchema),
  login
);
authRouter.post("/logout", requireAuth, logout);
authRouter.post(
  "/forgot-password",
  passwordResetRateLimiter,
  requireGuest,
  validate(forgotPasswordSchema),
  forgotPassword
);
authRouter.post(
  "/reset-password",
  passwordResetRateLimiter,
  requireGuest,
  validate(resetPasswordSchema),
  resetPassword
);
authRouter.get("/me", requireAuth, getMe);

router.use("/auth", authRouter);

router.get(
  '/admin/users',
  requireAuth,
  requireAdmin,
  validateQuery(getUsersQuerySchema),
  getUsers
);

router.get(
  '/admin/users/:id',
  requireAuth,
  requireAdmin,
  validateParams(userParamsSchema),
  getUserById
);

router.post(
  '/admin/users',
  requireAuth,
  requireAdmin,
  validate(createUserSchema),
  createUser
);

router.put(
  '/admin/users/:id',
  requireAuth,
  requireAdmin,
  validateParams(userParamsSchema),
  validate(updateUserSchema),
  updateUser
);

router.patch(
  '/admin/users/:id/password',
  requireAuth,
  requireAdmin,
  validateParams(userParamsSchema),
  validate(updateUserPasswordSchema),
  updateUserPassword
);

router.delete(
  '/admin/users/:id',
  requireAuth,
  requireAdmin,
  validateParams(userParamsSchema),
  deleteUser
);

export default router;
