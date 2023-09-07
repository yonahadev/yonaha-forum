"use client";
import { ApiResponse } from "@/interfaces";
import { useState } from "react";
import { FieldValues, useForm } from "react-hook-form";
import useAuth from "../hooks/useAuth";

const page = () => {
  const [response, setResponse] = useState({ error: "", message: "" });
  const { register, handleSubmit } = useForm();
  const onSubmit = async (data: FieldValues) => {
    try {
      const response = await fetch("http://127.0.0.1:8080/posts", {
        method: "POST",
        credentials: "include",
        body: JSON.stringify(data),
        headers: {
          "Content-Type": "application/json",
        },
      });
      const responseData: ApiResponse = await response.json();
      if (response.ok) {
        setResponse({ error: "", message: responseData.message! });
        location.replace("/");
      } else {
        setResponse({ error: responseData.error!, message: "" });
      }
      console.log(responseData);
    } catch (error) {
      setResponse({ error: "Network Error", message: "" });
    }
  };

  const { authenticated } = useAuth();

  return (
    <>
      {authenticated ? (
        <div>
          Create a new post
          <form onSubmit={handleSubmit(onSubmit)}>
            <input {...register("title")} placeholder="enter title"></input>
            <input
              {...register("text_content")}
              placeholder="enter text"
            ></input>
            <input type="submit" />
          </form>
          <p>{response?.error}</p>
          <p>{response?.message}</p>
        </div>
      ) : (
        <p>Not authenticated - please sign in or create an account </p>
      )}
    </>
  );
};

export default page;
