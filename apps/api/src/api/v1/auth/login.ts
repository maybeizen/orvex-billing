import { Request, Response } from 'express';
import { User } from "../../../db/models/User";
import { LoginInput } from "../../../validators/auth";
import { LoginResponse } from "../../../types/api";
import { ResponseHelper } from "../../../utils/response";

export const login = async (req: Request, res: Response) => {
  try {
    const { email, password } = req.body;

    const user = await User.findOne({ email });
    if (!user) {
      ResponseHelper.unauthorized(res, "Invalid credentials");
      return;
    }

    const isValidPassword = await user.comparePassword(password);
    if (!isValidPassword) {
      ResponseHelper.unauthorized(res, "Invalid credentials");
      return;
    }

    req.session.userId = user._id.toString();

    const userData: LoginResponse = {
      id: user._id.toString(),
      email: user.email,
      isEmailVerified: user.isEmailVerified,
      createdAt: user.createdAt,
      updatedAt: user.updatedAt,
      sessionId: req.sessionID,
    };

    ResponseHelper.success(res, userData, "Login successful");
  } catch (error) {
    console.error("Login error:", error);
    ResponseHelper.error(res, "Internal server error");
  }
};
