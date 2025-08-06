import { User } from "../../../../db/models/User";
import { ResponseHelper } from "../../../../utils/response";
import {
  AuthenticatedRequestHandler,
} from "../../../../types/api";
import { UpdateUserPasswordInput, UserParams } from "../../../../validators/users";

export const updateUserPassword: AuthenticatedRequestHandler<
  UserParams,
  UpdateUserPasswordInput,
  {},
  {}
> = async (req, res) => {
  try {
    const { id } = req.params;
    const { password } = req.body;

    const user = await User.findById(id);
    if (!user) {
      return ResponseHelper.notFound(res, 'User not found');
    }

    user.password = password;
    await user.save();

    return ResponseHelper.success(
      res,
      undefined,
      'Password updated successfully'
    );
  } catch (error) {
    console.error('Update user password error:', error);
    return ResponseHelper.error(res, 'Failed to update password');
  }
};