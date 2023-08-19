export interface User {
  id: number;
  username: string;
}

export interface Post {
  id: number;
  title: string;
  user: User;
}