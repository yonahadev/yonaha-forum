"use client";
import { ApiResponse } from "@/interfaces";
import { useState } from "react";
import { FieldValues, useForm } from "react-hook-form";

const page = () => {
  const [response, setResponse] = useState<ApiResponse | null>();
  const { register, handleSubmit } = useForm();
  const onSubmit = async (data: FieldValues) => {
    const response = await fetch("http://127.0.0.1:8080/users", {
      method: "POST",
      body: JSON.stringify(data),
      headers: {
        "Content-Type": "application/json",
      },
    });
    const responseData: ApiResponse = await response.json();
    setResponse(responseData);
    console.log(responseData);
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
