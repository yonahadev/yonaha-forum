import { Post } from "@/interfaces";

const getPosts = async (): Promise<Post[]> => {
  const res = await fetch("http://127.0.0.1:8080/posts", {
    next: { revalidate: 10 },
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
        <div className="w-1/2 bg-white p-4 h-fit mt-2">
          <p className="text-2xl">{post.title}</p>
          <p className="opacity-50">{post.user.username}</p>
          <p>{post.text_content}</p>
        </div>
      ))}
    </main>
  );
}
