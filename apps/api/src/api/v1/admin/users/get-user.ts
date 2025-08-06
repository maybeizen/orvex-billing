import { User } from "../../../../db/models/User";
import { ResponseHelper } from "../../../../utils/response";
import {
  AuthenticatedRequestHandler,
} from "../../../../types/api";
import { UserParams } from "../../../../validators/users";

interface UserResponse {
  id: string;
  uuid: string;
  firstName: string;
  lastName: string;
  username: string;
  email: string;
  role: string;
  isEmailVerified: boolean;
  createdAt: Date;
  updatedAt: Date;
}

const sanitizeUser = (user: any): UserResponse => ({
  id: user._id.toString(),
  uuid: user.uuid,
  firstName: user.firstName,
  lastName: user.lastName,
  username: user.username,
  email: user.email,
  role: user.role,
  isEmailVerified: user.isEmailVerified,
  createdAt: user.createdAt,
  updatedAt: user.updatedAt
});

export const getUserById: AuthenticatedRequestHandler<
  UserParams,
  {},
  {},
  UserResponse
> = async (req, res) => {
  try {
    const { id } = req.params;

    const user = await User.findById(id)
      .select('-password -resetPasswordToken -resetPasswordExpires')
      .lean();

    if (!user) {
      return ResponseHelper.notFound(res, 'User not found');
    }

    return ResponseHelper.success(
      res,
      sanitizeUser(user),
      'User retrieved successfully'
    );
  } catch (error) {
    console.error('Get user by ID error:', error);
    return ResponseHelper.error(res, 'Failed to retrieve user');
  }
};