export interface User {
  id: number;
  username: string;
  password: string; 
}

export interface Post {
  id: number;
  title: string;
  user: User;
}

export interface ApiResponse {
  error?: string;
  message?: string;
}