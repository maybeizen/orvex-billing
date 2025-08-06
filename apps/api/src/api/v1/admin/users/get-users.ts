import { User } from "../../../../db/models/User";
import { ResponseHelper } from "../../../../utils/response";
import {
  AuthenticatedRequestHandler,
} from "../../../../types/api";
import { GetUsersQuery } from "../../../../validators/users";

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

export const getUsers: AuthenticatedRequestHandler<
  {},
  {},
  GetUsersQuery,
  UserResponse[]
> = async (req, res) => {
  try {
    const { page = 1, limit = 10, search, role, isEmailVerified } = req.query;

    const filter: any = {};

    if (search) {
      filter.$or = [
        { firstName: { $regex: search, $options: 'i' } },
        { lastName: { $regex: search, $options: 'i' } },
        { username: { $regex: search, $options: 'i' } },
        { email: { $regex: search, $options: 'i' } }
      ];
    }

    if (role) {
      filter.role = role;
    }

    if (typeof isEmailVerified === 'boolean') {
      filter.isEmailVerified = isEmailVerified;
    }

    const skip = (page - 1) * limit;
    
    const [users, total] = await Promise.all([
      User.find(filter)
        .select('-password -resetPasswordToken -resetPasswordExpires')
        .sort({ createdAt: -1 })
        .skip(skip)
        .limit(limit)
        .lean(),
      User.countDocuments(filter)
    ]);

    const sanitizedUsers = users.map(sanitizeUser);

    return ResponseHelper.paginated(
      res,
      sanitizedUsers,
      {
        page,
        limit,
        total,
        totalPages: Math.ceil(total / limit)
      },
      'Users retrieved successfully'
    );
  } catch (error) {
    console.error('Get users error:', error);
    return ResponseHelper.error(res, 'Failed to retrieve users');
  }
};