import { Response } from "express";
import { AuthenticatedRequest } from "../../../types/api";

export const getMe = (req: AuthenticatedRequest, res: Response) => {
  res.json({
    success: true,
    user: {
      id: req.user._id,
      email: req.user.email,
      isEmailVerified: req.user.isEmailVerified,
      createdAt: req.user.createdAt,
    },
  });
};
