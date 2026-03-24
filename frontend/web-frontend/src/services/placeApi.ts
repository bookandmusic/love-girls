import api from "./api";
import type { FileInfo } from "./upload";

export interface Photo {
  id: number;
  placeId: number;
  file?: FileInfo;
}

export interface Place {
  id: number;
  name: string;
  latitude: number;
  longitude: number;
  image?: Photo;
  description: string;
  date: string;
}

interface GetPlacesResponse {
  code: number;
  data: {
    places: Place[];
    totalPages: number;
    total?: number;
    totalCount?: number;
    page?: number;
    size?: number;
  };
  msg?: string;
}

export const placeApi = {
  async getPlaces(page: number, size: number) {
    const response = await api.get<GetPlacesResponse>("/places", {
      params: {
        page,
        size,
      },
    });
    return response.data;
  },
};
