"use client";
import { ApiResponse } from "@/interfaces";
import Cookies from "js-cookie";
import Link from "next/link";
import { useState } from "react";
import { FieldValues, useForm } from "react-hook-form";

const page = () => {
  const [response, setResponse] = useState({ error: "", message: "" });
  const { register, handleSubmit } = useForm();
  const onSubmit = async (data: FieldValues) => {
    try {
      const response = await fetch("http://127.0.0.1:8080/users/signin", {
        method: "POST",
        body: JSON.stringify(data),
        headers: {
          "Content-Type": "application/json",
        },
      });
      const responseData: ApiResponse = await response.json();
      if (response.ok) {
        setResponse({ error: "", message: responseData.message! });
        if (responseData.token) {
          Cookies.set("jwtToken", responseData.token, {
            httpOnly: false,
            expires: new Date(responseData.tokenExpiry!),
          });

          console.log("Signed in");
          location.replace("/");
        }
      } else {
        setResponse({ error: responseData.error!, message: "" });
      }
    } catch (error) {
      setResponse({ error: "Network Error", message: "" });
    }
  };

  return (
    <div>
      Sign in
      <form onSubmit={handleSubmit(onSubmit)}>
        <input {...register("username")} placeholder="enter username"></input>
        <input {...register("password")} placeholder="enter password"></input>
        <input type="submit" />
      </form>
      <p>{response?.error}</p>
      <p>{response?.message}</p>
      Don't have an account?{" "}
      <Link className="text-blue-600 underline" href={"/signup"}>
        Sign Up
      </Link>
    </div>
  );
};

export default page;
