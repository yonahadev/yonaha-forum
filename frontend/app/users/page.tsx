import { User } from "@/interfaces";
import Link from "next/link";
import React from "react";

const getUsers = async (): Promise<User[]> => {
  const res = await fetch("http://127.0.0.1:8080/users");
  if (!res.ok) {
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
          <p>{user.username}</p>
        </>
      ))}
      <Link href={"/"}>To Homepage</Link>
    </main>
  );
};

export default page;
