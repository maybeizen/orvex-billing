import { Redis } from "ioredis";
import { env } from "./env";

const redis = new Redis({
  host: env.REDIS_HOST,
  port: parseInt(env.REDIS_PORT, 10),
  password: env.REDIS_PASSWORD,
});

export default redis;
