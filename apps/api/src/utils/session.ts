import session, { Store } from "express-session";
import MongoStore from "connect-mongo";
import mongoose from "mongoose";
import { env } from "./env";
import type { RequestHandler } from "express";

export const sessionMiddleware: RequestHandler = session({
  name: "sid",
  secret: env.SESSION_SECRET,
  resave: false,
  saveUninitialized: false,
  rolling: true,
  store: MongoStore.create({
    mongoUrl: env.MONGODB_URI,
    collectionName: "sessions",
    ttl: 60 * 60 * 24 * 7, // 7 days
    touchAfter: 24 * 3600,
  }) as unknown as Store,
  cookie: {
    httpOnly: true,
    secure: env.NODE_ENV === "production",
    sameSite: env.NODE_ENV === "production" ? "strict" : "lax",
    maxAge: 1000 * 60 * 60 * 24 * 7, // 7 days
  },
});
