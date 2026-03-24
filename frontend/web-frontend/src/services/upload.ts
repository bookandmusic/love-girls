export interface FileInfo {
  id: number;
  url: string;
  thumbnail: string;
  name?: string;
  size?: number;
  mime_type?: string;
}

export interface UploadResponse {
  file: FileInfo;
}
