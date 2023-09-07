import { Post } from "@/interfaces";
import TextPost from "./components/TextPost";
import useAuth from "./hooks/useAuth";

const getPosts = async (): Promise<Post[]> => {
  const res = await fetch("http://127.0.0.1:8080/posts", {
    cache: "no-store",
  });
  if (!res.ok) {
    throw new Error("fetch failed");
  }

  return res.json();
};

export default async function Home() {
  const posts = await getPosts();

  return (
    <main className="w-full h-[calc(100%-4rem)] flex justify-center">
      {posts.map((post) => (
        <TextPost
          title={post.title}
          content={post.text_content}
          username={post.user.username}
          postID={post.id}
        />
      ))}
    </main>
  );
}
