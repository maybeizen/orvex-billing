import { Request, Response } from "express";
import { User } from "../../../db/models/User";
import { RegisterInput } from "../../../validators/auth";
import { AuthUserResponse } from "../../../types/api";
import { ResponseHelper } from "../../../utils/response";

export const register = async (req: Request, res: Response) => {
  try {
    const { email, password } = req.body;

    const existingUser = await User.findOne({ email });
    if (existingUser) {
      ResponseHelper.conflict(res, "User already exists with this email");
      return;
    }

    const user = new User({ email, password });
    await user.save();

    const userId = user._id.toString();
    req.session.regenerate((err) => {
      if (err) {
        console.error("Session regeneration error:", err);
        return ResponseHelper.error(res, "Internal server error");
      }

      req.session.userId = userId;

      const userData: AuthUserResponse = {
        id: user._id.toString(),
        email: user.email,
        isEmailVerified: user.isEmailVerified,
        createdAt: user.createdAt,
        updatedAt: user.updatedAt,
      };

      ResponseHelper.success(
        res,
        userData,
        "User registered successfully",
        201
      );
    });
    return;
  } catch (error) {
    console.error("Register error:", error);
    ResponseHelper.error(res, "Internal server error");
  }
};
