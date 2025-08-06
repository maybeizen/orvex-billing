import { User } from "../../../../db/models/User";
import { ResponseHelper } from "../../../../utils/response";
import {
  AuthenticatedRequestHandler,
} from "../../../../types/api";
import { CreateUserInput } from "../../../../validators/users";

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

export const createUser: AuthenticatedRequestHandler<
  {},
  CreateUserInput,
  {},
  UserResponse
> = async (req, res) => {
  try {
    const userData = req.body;

    const existingUser = await User.findOne({
      $or: [{ email: userData.email }, { username: userData.username }]
    });

    if (existingUser) {
      if (existingUser.email === userData.email) {
        return ResponseHelper.conflict(res, 'Email already exists');
      }
      if (existingUser.username === userData.username) {
        return ResponseHelper.conflict(res, 'Username already exists');
      }
    }

    const newUser = new User(userData);
    await newUser.save();

    const userResponse = sanitizeUser(newUser.toObject());

    return ResponseHelper.success(
      res,
      userResponse,
      'User created successfully',
      201
    );
  } catch (error) {
    console.error('Create user error:', error);
    return ResponseHelper.error(res, 'Failed to create user');
  }
};