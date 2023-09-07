"use client";
import { ApiResponse } from "@/interfaces";
import Link from "next/link";
import React, { useState } from "react";
import useAuth from "../hooks/useAuth";

const LoginButton = () => {
  const { authenticated, signOut, user } = useAuth();

  const logout = async () => {
    const data = await signOut();
    if (data.error != "") {
      console.log(data.error);
    } else {
      console.log("redirecting");
      location.replace("/");
    }
  };
  return (
    <div>
      {authenticated ? (
        <p
          className="p-2 bg-blue-500 text-lg rounded-md"
          onClick={() => logout()}
        >
          Sign Out {user?.username}
        </p>
      ) : (
        <Link className="p-2 bg-blue-500 text-lg rounded-md" href={"/signin"}>
          Sign in
        </Link>
      )}
    </div>
  );
};

export default LoginButton;
