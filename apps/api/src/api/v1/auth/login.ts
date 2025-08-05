import { Request, Response } from "express";
import { User } from "../../../db/models/User";
import { LoginInput } from "../../../validators/auth";
import { LoginResponse } from "../../../types/api";
import { ResponseHelper } from "../../../utils/response";

export const login = async (req: Request, res: Response) => {
  try {
    const { email, password } = req.body;

    const user = await User.findOne({ email });
    if (!user) {
      ResponseHelper.unauthorized(res, "Invalid email or password");
      return;
    }

    const isValidPassword = await user.comparePassword(password);
    if (!isValidPassword) {
      ResponseHelper.unauthorized(res, "Invalid email or password");
      return;
    }

    const userId = user._id.toString();
    req.session.regenerate((err) => {
      if (err) {
        console.error("Session regeneration error:", err);
        return ResponseHelper.error(res, "Internal server error");
      }

      req.session.userId = userId;

      const userData: LoginResponse = {
        id: user._id.toString(),
        email: user.email,
        isEmailVerified: user.isEmailVerified,
        createdAt: user.createdAt,
        updatedAt: user.updatedAt,
      };

      ResponseHelper.success(res, userData, "Login successful");
    });
  } catch (error) {
    console.error("Login error:", error);
    ResponseHelper.error(res, "Internal server error");
  }
};
