import { Request, Response } from 'express';

export const logout = (req: Request, res: Response) => {
  req.session.destroy((err) => {
    if (err) {
      console.error('Logout error:', err);
      return res.status(500).json({
        success: false,
        message: 'Failed to logout'
      });
    }

    res.clearCookie('sid');
    res.json({
      success: true,
      message: 'Logout successful'
    });
  });
};