"use client";

import React, {
  createContext,
  useContext,
  useState,
  useCallback,
  ReactNode,
} from "react";
import { ToastType, ToastPosition } from "@/components/ui/toast";

export interface NotificationData {
  id: string;
  type: ToastType;
  title?: string;
  message: string;
  duration?: number;
  position?: ToastPosition;
  disabled?: boolean;
}

interface NotificationContextType {
  notifications: NotificationData[];
  addNotification: (notification: Omit<NotificationData, "id">) => string;
  removeNotification: (id: string) => void;
  clearAll: () => void;
  success: (
    message: string,
    options?: Partial<Omit<NotificationData, "id" | "type" | "message">>
  ) => string;
  error: (
    message: string,
    options?: Partial<Omit<NotificationData, "id" | "type" | "message">>
  ) => string;
  warning: (
    message: string,
    options?: Partial<Omit<NotificationData, "id" | "type" | "message">>
  ) => string;
  info: (
    message: string,
    options?: Partial<Omit<NotificationData, "id" | "type" | "message">>
  ) => string;
}

const NotificationContext = createContext<NotificationContextType | undefined>(
  undefined
);

export const useNotifications = (): NotificationContextType => {
  const context = useContext(NotificationContext);
  if (!context) {
    throw new Error(
      "useNotifications must be used within a NotificationProvider"
    );
  }
  return context;
};

interface NotificationProviderProps {
  children: ReactNode;
  defaultPosition?: ToastPosition;
  defaultDuration?: number;
}

export const NotificationProvider: React.FC<NotificationProviderProps> = ({
  children,
  defaultPosition = "top-right",
  defaultDuration = 5000,
}) => {
  const [notifications, setNotifications] = useState<NotificationData[]>([]);

  const generateId = useCallback((): string => {
    return `toast-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
  }, []);

  const addNotification = useCallback(
    (notification: Omit<NotificationData, "id">): string => {
      const id = generateId();
      const newNotification: NotificationData = {
        id,
        position: defaultPosition,
        duration: defaultDuration,
        ...notification,
      };

      setNotifications((prev) => [...prev, newNotification]);
      return id;
    },
    [generateId, defaultPosition, defaultDuration]
  );

  const removeNotification = useCallback((id: string): void => {
    setNotifications((prev) =>
      prev.filter((notification) => notification.id !== id)
    );
  }, []);

  const clearAll = useCallback((): void => {
    setNotifications([]);
  }, []);

  const success = useCallback(
    (
      message: string,
      options?: Partial<Omit<NotificationData, "id" | "type" | "message">>
    ): string => {
      return addNotification({ type: "success", message, ...options });
    },
    [addNotification]
  );

  const error = useCallback(
    (
      message: string,
      options?: Partial<Omit<NotificationData, "id" | "type" | "message">>
    ): string => {
      return addNotification({ type: "error", message, ...options });
    },
    [addNotification]
  );

  const warning = useCallback(
    (
      message: string,
      options?: Partial<Omit<NotificationData, "id" | "type" | "message">>
    ): string => {
      return addNotification({ type: "warning", message, ...options });
    },
    [addNotification]
  );

  const info = useCallback(
    (
      message: string,
      options?: Partial<Omit<NotificationData, "id" | "type" | "message">>
    ): string => {
      return addNotification({ type: "info", message, ...options });
    },
    [addNotification]
  );

  const value: NotificationContextType = {
    notifications,
    addNotification,
    removeNotification,
    clearAll,
    success,
    error,
    warning,
    info,
  };

  return (
    <NotificationContext.Provider value={value}>
      {children}
    </NotificationContext.Provider>
  );
};
