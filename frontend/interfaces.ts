export interface User {
  id: number;
  username: string;
  password: string; 
}

export interface Post {
  id: number;
  title: string;
  text_content: string;
  user: User;
}

export interface ApiResponse {
  error?: string;
  message?: string;
  token?: string;
  user?: User;
  tokenExpiry?: number
}
