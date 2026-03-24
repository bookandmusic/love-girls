import api from "./api";
import type { FileInfo } from "./upload";

export interface Photo {
  id: number;
  albumId: number;
  file?: FileInfo;
  alt?: string;
  createdAt: string;
}

export interface Album {
  id: number;
  name: string;
  description: string;
  coverImage?: Photo;
  photoCount: number;
  createdAt: string;
}

interface GetAlbumsResponse {
  code: number;
  data: {
    albums: Album[];
    totalPages: number;
    total?: number;
    totalCount?: number;
    page?: number;
    size?: number;
  };
  message?: string;
}

interface GetPhotosResponse {
  code: number;
  data: {
    photos: Photo[];
    totalPages: number;
    total?: number;
    totalCount?: number;
    page?: number;
    size?: number;
  };
  message?: string;
}

export const albumApi = {
  async getAlbums(page: number, size: number) {
    const response = await api.get<GetAlbumsResponse>("/albums", {
      params: {
        page,
        size,
        sort_by: "created_at",
        order: "desc",
      },
    });
    return response.data;
  },

  async getPhotos(albumId: number, page: number, size: number) {
    const response = await api.get<GetPhotosResponse>(
      `/albums/${albumId}/photos`,
      {
        params: {
          page,
          size,
        },
      },
    );
    return response.data;
  },
};
