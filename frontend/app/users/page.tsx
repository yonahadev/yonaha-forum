import { User } from "@/interfaces";
import Link from "next/link";

const getUsers = async (): Promise<User[]> => {
  const res = await fetch("http://127.0.0.1:8080/users", {
    method: "GET",
  });
  if (!res.ok) {
    console.error("Fetch failed:", res.statusText);
    throw new Error("fetch failed");
  }
  return res.json();
};

const page = async () => {
  const users = await getUsers();

  return (
    <main>
      <h1>Current Users</h1>
      {users.map((user) => (
        <>
          <p key={user.id}>{user.username}</p>
        </>
      ))}
      <Link href={"/"}>To Homepage</Link>
    </main>
  );
};

export default page;
