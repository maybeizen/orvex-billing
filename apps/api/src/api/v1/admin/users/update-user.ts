import { User } from "../../../../db/models/User";
import { ResponseHelper } from "../../../../utils/response";
import {
  AuthenticatedRequestHandler,
} from "../../../../types/api";
import { UpdateUserInput, UserParams } from "../../../../validators/users";

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

export const updateUser: AuthenticatedRequestHandler<
  UserParams,
  UpdateUserInput,
  {},
  UserResponse
> = async (req, res) => {
  try {
    const { id } = req.params;
    const updateData = req.body;

    if (Object.keys(updateData).length === 0) {
      return ResponseHelper.error(res, 'No update data provided', 400);
    }

    if (updateData.email || updateData.username) {
      const filter: any = { _id: { $ne: id } };
      const orConditions: any[] = [];

      if (updateData.email) {
        orConditions.push({ email: updateData.email });
      }
      if (updateData.username) {
        orConditions.push({ username: updateData.username });
      }

      filter.$or = orConditions;

      const existingUser = await User.findOne(filter);
      if (existingUser) {
        if (existingUser.email === updateData.email) {
          return ResponseHelper.conflict(res, 'Email already exists');
        }
        if (existingUser.username === updateData.username) {
          return ResponseHelper.conflict(res, 'Username already exists');
        }
      }
    }

    const updatedUser = await User.findByIdAndUpdate(
      id,
      updateData,
      { new: true, runValidators: true }
    ).select('-password -resetPasswordToken -resetPasswordExpires');

    if (!updatedUser) {
      return ResponseHelper.notFound(res, 'User not found');
    }

    return ResponseHelper.success(
      res,
      sanitizeUser(updatedUser.toObject()),
      'User updated successfully'
    );
  } catch (error) {
    console.error('Update user error:', error);
    return ResponseHelper.error(res, 'Failed to update user');
  }
};