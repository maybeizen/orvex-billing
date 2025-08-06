import { Response } from "express";
import { AuthenticatedRequest } from "../../../types/api";

export const getMe = (req: AuthenticatedRequest, res: Response) => {
  res.json({
    success: true,
    user: {
      id: req.user._id,
      uuid: req.user.uuid,
      firstName: req.user.firstName,
      lastName: req.user.lastName,
      username: req.user.username,
      email: req.user.email,
      isEmailVerified: req.user.isEmailVerified,
      createdAt: req.user.createdAt,
      updatedAt: req.user.updatedAt,
    },
  });
};
