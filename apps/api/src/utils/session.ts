import session, { Store } from "express-session";
import MongoStore from "connect-mongo";
import mongoose from "mongoose";
import { env } from "./env";

export const sessionMiddleware = session({
  name: "sid",
  secret: env.SESSION_SECRET,
  resave: false,
  saveUninitialized: false,
  store: MongoStore.create({
    client: mongoose.connection.getClient(),
    collectionName: "sessions",
    ttl: 60 * 60 * 24 * 5, // 5 days
  }) as unknown as Store,
  cookie: {
    httpOnly: true,
    secure: env.NODE_ENV === "production",
    sameSite: "lax",
    maxAge: 1000 * 60 * 60 * 24 * 5, // 5 days
  },
});
