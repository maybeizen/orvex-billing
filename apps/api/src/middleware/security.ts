import { Request, Response, NextFunction } from "express";
import { env } from "../utils/env";

export const securityHeaders = (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  res.setHeader("X-Frame-Options", "DENY");
  res.setHeader("X-Content-Type-Options", "nosniff");
  res.setHeader("X-XSS-Protection", "1; mode=block");
  res.setHeader("Referrer-Policy", "strict-origin-when-cross-origin");
  res.setHeader("X-Permitted-Cross-Domain-Policies", "none");

  if (env.NODE_ENV === "production") {
    res.setHeader(
      "Strict-Transport-Security",
      "max-age=31536000; includeSubDomains; preload"
    );
  }

  res.setHeader(
    "Content-Security-Policy",
    "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self'; frame-ancestors 'none';"
  );

  next();
};

export const csrfProtection = (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  if (["POST", "PUT", "PATCH", "DELETE"].includes(req.method)) {
    const origin = req.get("Origin");
    const referer = req.get("Referer");

    if (origin && origin === env.CORS_ORIGIN) {
      return next();
    }

    if (referer && referer.startsWith(env.CORS_ORIGIN)) {
      return next();
    }

    if (
      req.is("application/json") &&
      req.headers["content-type"]?.includes("application/json")
    ) {
      return next();
    }

    return res.status(403).json({
      success: false,
      message: "Forbidden: Invalid request origin",
    });
  }

  next();
};
