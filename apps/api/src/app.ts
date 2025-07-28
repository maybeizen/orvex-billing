import express, { type Express } from "express";
import cors from "cors";
import { env } from "./utils/env";
import { apiRateLimiter } from "./middleware/rate-limit";

const app: Express = express();

app.use(apiRateLimiter);
app.use(
  cors({
    origin: env.CORS_ORIGIN,
    credentials: true,
  })
);

app.use(express.json());

export default app;
