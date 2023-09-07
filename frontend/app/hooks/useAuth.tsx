"use client";
import { ApiResponse, User } from "@/interfaces";
import { error } from "console";
import Cookies from "js-cookie";
import React, { useEffect, useState } from "react";

const useAuth = () => {
  const [authenticated, setAuthenticated] = useState<boolean>();
  const [user, setUser] = useState<User>();

  useEffect(() => {
    const checkAuth = async () => {
      const res = await fetch("http://127.0.0.1:8080/auth/getUser", {
        credentials: "include",
        cache: "no-store",
      });
      const responseData: ApiResponse = await res.json();
      if (responseData.user) {
        setUser(responseData.user);
      } else {
        console.error(responseData.error);
      }
      responseData.error ? setAuthenticated(false) : setAuthenticated(true);
    };

    checkAuth();
  }, []);

  const signOut = async () => {
    const res = await fetch("http://127.0.0.1:8080/users/signout", {
      credentials: "include",
    });
    const responseData: ApiResponse = await res.json();
    if (res.ok) {
      Cookies.set("jwtToken", "", {
        httpOnly: false,
        expires: 0,
      });
      return { error: "", message: responseData.message! };
    } else {
      return { error: responseData.error!, message: "" };
    }
  };
  return {
    authenticated,
    signOut,
    user,
  };
};

export default useAuth;
