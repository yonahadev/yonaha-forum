"use client";
import { ApiResponse } from "@/interfaces";
import Cookies from "js-cookie";
import { useState } from "react";
import { FieldValues, useForm } from "react-hook-form";

const page = () => {
  const [response, setResponse] = useState<ApiResponse | null>();
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
          Cookies.set("jwtToken", responseData.token, { httpOnly: false });
          console.log("Signed in");
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
    </div>
  );
};

export default page;
