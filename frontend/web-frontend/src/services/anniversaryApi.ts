import api from "./api";

export interface Anniversary {
  id: number;
  title: string;
  date: string;
  description: string;
  calendar: "solar" | "lunar";
}

interface GetAnniversariesResponse {
  code: number;
  data: {
    anniversaries: Anniversary[];
    totalPages: number;
    total?: number;
    totalCount?: number;
    page?: number;
    size?: number;
  };
  msg?: string;
}

export const anniversaryApi = {
  async getAnniversaries(page: number, size: number) {
    const response = await api.get<GetAnniversariesResponse>("/anniversaries", {
      params: {
        page,
        size,
        sort_by: "created_at",
        order: "desc",
      },
    });
    return response.data;
  },
};
