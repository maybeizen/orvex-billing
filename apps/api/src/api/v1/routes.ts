import express, { Router } from "express";
import { requireAuth, requireGuest } from "../../middleware/auth";
import {
  forgotPasswordSchema,
  resetPasswordSchema,
  loginSchema,
  registerSchema,
} from "../../validators/auth";
import { validate } from "../../middleware/validation";
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

const router: Router = Router();
const authRouter: Router = Router();

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

export default router;
