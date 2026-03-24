import api from "./api";
import type { FileInfo } from "./upload";

export interface Photo {
  id: number;
  momentId: number;
  file?: FileInfo;
}

export interface Moment {
  id: number;
  content: string;
  images?: Photo[];
  imageIds?: number[];
  likes: number;
  createdAt: string;
  author: {
    name: string;
    avatar?: FileInfo;
  };
  isPublic: boolean;
  userId?: number;
}

interface GetMomentsResponse {
  code: number;
  data: {
    moments: Moment[];
    totalPages: number;
    total?: number;
    totalCount?: number;
    page?: number;
    size?: number;
  };
  msg?: string;
}

interface LikeMomentResponse {
  code: number;
  data: {
    likes: number;
  };
  msg?: string;
}

export const momentApi = {
  async getMoments(page: number, size: number) {
    const response = await api.get<GetMomentsResponse>("/moments", {
      params: {
        page,
        size,
        sort_by: "created_at",
        order: "desc",
      },
    });
    return response.data;
  },

  async likeMoment(momentId: number) {
    const response = await api.post<LikeMomentResponse>(
      `/moments/${momentId}/like`,
    );
    return response.data;
  },
};
