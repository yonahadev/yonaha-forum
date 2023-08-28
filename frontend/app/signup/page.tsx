"use client";
import { ApiResponse } from "@/interfaces";
import { useState } from "react";
import { FieldValues, useForm } from "react-hook-form";

const page = () => {
  const [response, setResponse] = useState({ error: "", message: "" });
  const { register, handleSubmit } = useForm();
  const onSubmit = async (data: FieldValues) => {
    try {
      const response = await fetch("http://127.0.0.1:8080/users", {
        method: "POST",
        body: JSON.stringify(data),
        headers: {
          "Content-Type": "application/json",
        },
      });
      const responseData: ApiResponse = await response.json();
      if (response.ok) {
        setResponse({ error: "", message: responseData.message! });
      } else {
        setResponse({ error: responseData.error!, message: "" });
      }
      console.log(responseData);
    } catch (error) {
      setResponse({ error: "Network Error", message: "" });
    }
  };

  return (
    <div>
      Create an account!
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
