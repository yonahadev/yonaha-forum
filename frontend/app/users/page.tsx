import { User } from "@/interfaces";
import Link from "next/link";
import React, { useEffect, useState } from "react";

const getUsers = async (): Promise<User[]> => {
  const res = await fetch("http://127.0.0.1:8080/users", { method: "GET" });
  if (!res.ok) {
    throw new Error("fetch failed");
  }
  return res.json();
};

const page = async () => {
  const users = await getUsers();
  // const [users, setUsers] = useState<User[]>([]);

  // useEffect(() => {
  //   const fetchData = async () => {
  //     try {
  //       const userData = await getUsers();
  //       setUsers(userData);
  //     } catch (error) {
  //       console.log("error", error);
  //     }
  //   };
  //   fetchData();
  // }, []);

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
