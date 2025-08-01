export interface IUser {
  _id: string;
  email: string;
  password: string;
  isEmailVerified: boolean;
  resetPasswordToken?: string;
  resetPasswordExpires?: Date;
  createdAt: Date;
  updatedAt: Date;
  role: string;
}

declare module "express-session" {
  interface SessionData {
    userId?: string;
  }
}
