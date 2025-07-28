import rateLimit from "express-rate-limit";
import { env } from "../utils/env";

export const apiRateLimiter = rateLimit({
  windowMs: parseInt(env.RATE_LIMIT_WINDOW) * 60 * 1000,
  max: parseInt(env.RATE_LIMIT_MAX),
  standardHeaders: true,
  legacyHeaders: false,
  message: {
    success: false,
    message: "Too many requests, please try again later.",
  },
});
