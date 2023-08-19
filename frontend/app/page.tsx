import Link from "next/link";

const getPosts = async (): Promise<Post[]> => {
  const res = await fetch("http://127.0.0.1:8080/posts");
  if (!res.ok) {
    throw new Error("fetch failed");
  }

  return res.json();
};

interface User {
  id: number;
  username: string;
}

interface Post {
  id: number;
  title: string;
  user: User;
}

export default async function Home() {
  const posts = await getPosts();

  return (
    <main>
      <h1>yonaha-forum</h1>
      {posts.map((post) => (
        <>
          <p>{post.title}</p>
          <p>Posted by {post.user.username}</p>
        </>
      ))}
      <Link href={"/signup"}>Sign Up</Link>
    </main>
  );
}
