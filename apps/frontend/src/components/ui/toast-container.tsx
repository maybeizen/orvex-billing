"use client";

import React from "react";
import { createPortal } from "react-dom";
import { Toast, ToastPosition } from "./toast";
import { useNotifications } from "@/contexts/notification-context";

const positionStyles: Record<ToastPosition, string> = {
  "top-left": "top-4 left-4",
  "top-right": "top-4 right-4",
  "top-center": "top-4 left-1/2 transform -translate-x-1/2",
  "bottom-left": "bottom-4 left-4",
  "bottom-right": "bottom-4 right-4",
  "bottom-center": "bottom-4 left-1/2 transform -translate-x-1/2",
};

interface ToastContainerProps {
  position?: ToastPosition;
  maxToasts?: number;
}

export const ToastContainer: React.FC<ToastContainerProps> = ({
  position = "top-right",
  maxToasts = 5,
}) => {
  const { notifications, removeNotification } = useNotifications();

  const notificationsByPosition = notifications.reduce((acc, notification) => {
    const pos = notification.position || position;
    if (!acc[pos]) acc[pos] = [];
    acc[pos].push(notification);
    return acc;
  }, {} as Record<ToastPosition, typeof notifications>);

  if (notifications.length === 0) {
    return null;
  }

  return createPortal(
    <>
      {Object.entries(notificationsByPosition).map(
        ([pos, positionNotifications]) => {
          const positionKey = pos as ToastPosition;
          const limitedNotifications = positionNotifications.slice(-maxToasts);

          return (
            <div
              key={positionKey}
              className={`
              fixed z-50 pointer-events-none
              ${positionStyles[positionKey]}
              flex flex-col
              ${positionKey.includes("top") ? "items-start" : "items-start"}
              ${positionKey.includes("center") ? "items-center" : ""}
            `}
              style={{
                maxWidth: "calc(100vw - 2rem)",
                width: "auto",
              }}
            >
              {limitedNotifications.map((notification) => (
                <div key={notification.id} className="pointer-events-auto">
                  <Toast {...notification} onClose={removeNotification} />
                </div>
              ))}
            </div>
          );
        }
      )}
    </>,
    document.body
  );
};
