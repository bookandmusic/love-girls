import api from "./api";
import type { FileInfo } from "./upload";

export interface NotificationSender {
  id: number;
  name: string;
  avatar?: FileInfo;
}

export interface Notification {
  id: number;
  type: "comment" | "reply";
  momentId: number;
  commentId: number;
  sender: NotificationSender;
  content: string;
  isRead: boolean;
  createdAt: string;
}

interface GetNotificationsResponse {
  code: number;
  data: {
    notifications: Notification[];
    total: number;
    page: number;
    size: number;
  };
  msg?: string;
}

interface GetUnreadCountResponse {
  code: number;
  data: {
    count: number;
  };
  msg?: string;
}

interface MarkReadResponse {
  code: number;
  msg?: string;
}

export const notificationApi = {
  async getUnreadNotifications(page = 1, size = 20) {
    const response = await api.get<GetNotificationsResponse>(
      "/notifications/unread",
      {
        params: { page, size },
      },
    );
    return response.data;
  },

  async markAsRead(notificationId: number) {
    const response = await api.post<MarkReadResponse>(
      `/notifications/${notificationId}/read`,
    );
    return response.data;
  },

  async getUnreadCount() {
    const response = await api.get<GetUnreadCountResponse>(
      "/notifications/count",
    );
    return response.data;
  },

  async markAllAsRead() {
    const response = await api.post<MarkReadResponse>(
      "/notifications/read-all",
    );
    return response.data;
  },
};
