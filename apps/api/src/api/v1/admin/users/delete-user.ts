import { User } from "../../../../db/models/User";
import { ResponseHelper } from "../../../../utils/response";
import {
  AuthenticatedRequestHandler,
} from "../../../../types/api";
import { UserParams } from "../../../../validators/users";

export const deleteUser: AuthenticatedRequestHandler<
  UserParams,
  {},
  {},
  {}
> = async (req, res) => {
  try {
    const { id } = req.params;

    if (req.user._id.toString() === id) {
      return ResponseHelper.error(res, 'Cannot delete your own account', 400);
    }

    const deletedUser = await User.findByIdAndDelete(id);

    if (!deletedUser) {
      return ResponseHelper.notFound(res, 'User not found');
    }

    return ResponseHelper.success(
      res,
      undefined,
      'User deleted successfully'
    );
  } catch (error) {
    console.error('Delete user error:', error);
    return ResponseHelper.error(res, 'Failed to delete user');
  }
};