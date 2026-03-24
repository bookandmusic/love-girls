import api from "./api";
import type { FileInfo } from "./upload";

export interface SystemInfo {
  site: {
    name: string;
    description: string;
    startDate: string;
  };
  couple: {
    boy: {
      name: string;
      avatar?: FileInfo;
    };
    girl: {
      name: string;
      avatar?: FileInfo;
    };
  };
}

interface GetSystemInfoResponse {
  code: number;
  data: SystemInfo;
}

export const systemApi = {
  getSystemInfo() {
    return api.get<GetSystemInfoResponse>("/system/info");
  },
};
