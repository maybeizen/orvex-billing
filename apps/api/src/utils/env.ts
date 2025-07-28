import { z } from "zod";
import "dotenv/config";

const envSchema = z.object({
  PORT: z.string().default("3001"),
  NODE_ENV: z.enum(["development", "production"]).default("development"),

  MONGODB_URI: z.url(),

  SESSION_SECRET: z.string().min(32),

  CORS_ORIGIN: z.url(),

  REDIS_HOST: z.string(),
  REDIS_PORT: z.string(),
  REDIS_PASSWORD: z.string().optional(),

  RATE_LIMIT_WINDOW: z.string().default("10"),
  RATE_LIMIT_MAX: z.string().default("100"),

  TWOFA_ISSUER: z.string(),
  TWOFA_SECRET_KEY: z.string(),

  EMAIL_HOST: z.string(),
  EMAIL_PORT: z.string(),
  EMAIL_USER: z.string(),
  EMAIL_PASS: z.string(),
  EMAIL_FROM: z.string(),

  AVATAR_UPLOAD_DIR: z.string(),
  CDN_BASE_URL: z.string(),

  PTERO_API_URL: z.string().url(),
  PTERO_API_KEY: z.string(),

  STRIPE_SECRET_KEY: z.string(),
  STRIPE_WEBHOOK_SECRET: z.string(),

  APP_NAME: z.string(),
  APP_URL: z.string().url(),
});

const _env = envSchema.safeParse(process.env);

if (!_env.success) {
  console.error(`❌ Invalid environment variables: ${_env.error}`);

  process.exit(1);
}

export const env = _env.data;
