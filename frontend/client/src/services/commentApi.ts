import api from "./api";
import type { FileInfo } from "./upload";

export interface CommentAuthor {
  id: number;
  name: string;
  avatar?: FileInfo;
}

export interface Comment {
  id: number;
  content: string;
  momentId: number;
  parentId?: number;
  replyToId?: number;
  replyTo?: CommentAuthor;
  userId: number;
  author: CommentAuthor;
  depth: number;
  createdAt: string;
  children?: Comment[];
}

interface GetCommentsResponse {
  code: number;
  data: {
    comments: Comment[];
    total: number;
    page: number;
    size: number;
  };
  msg?: string;
}

interface CreateCommentResponse {
  code: number;
  data: Comment;
  msg?: string;
}

interface DeleteCommentResponse {
  code: number;
  msg?: string;
}

export const commentApi = {
  async getComments(momentId: number, page = 1, size = 20) {
    const response = await api.get<GetCommentsResponse>(
      `/moments/${momentId}/comments`,
      {
        params: { page, size },
      },
    );
    return response.data;
  },

  async createComment(
    momentId: number,
    data: {
      content: string;
      parentId?: number;
      replyToId?: number;
    },
  ) {
    const response = await api.post<CreateCommentResponse>(
      `/moments/${momentId}/comments`,
      data,
    );
    return response.data;
  },

  async deleteComment(commentId: number) {
    const response = await api.delete<DeleteCommentResponse>(
      `/comments/${commentId}`,
    );
    return response.data;
  },
};
