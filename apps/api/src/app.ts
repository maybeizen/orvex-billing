import express, { type Express } from "express";
import cors from "cors";
import { env } from "./utils/env";
import { apiRateLimiter } from "./middleware/rate-limit";
import { sessionMiddleware } from "./utils/session";
import { errorHandler, notFoundHandler } from "./middleware/error-handler";
import { securityHeaders, csrfProtection } from "./middleware/security";
import v1Routes from "./api/v1/routes";

const app: Express = express();

app.use(securityHeaders);
app.use(apiRateLimiter);
app.use(
  cors({
    origin: env.CORS_ORIGIN,
    credentials: true,
  })
);

app.use(express.json({ limit: '10mb' }));
app.use(express.urlencoded({ extended: true, limit: '10mb' }));
app.use(sessionMiddleware);
app.use(csrfProtection);

app.use("/api/v1", v1Routes);

app.use(notFoundHandler);
app.use(errorHandler);

export default app;
