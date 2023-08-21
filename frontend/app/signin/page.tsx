"use client";
import { useState } from "react";
import { FieldValues, useForm } from "react-hook-form";

const page = () => {
  const [response, setResponse] = useState(0);
  const { register, handleSubmit } = useForm();
  const onSubmit = async (data: FieldValues) => {
    const response = await fetch("http://127.0.0.1:8080/users", {
      method: "POST",
      body: JSON.stringify(data),
      headers: {
        "Content-Type": "application/json",
      },
    });
    const result = await response.json();
    const status = response.status;
    setResponse(status);
  };

  return (
    <div>
      Sign in
      <form onSubmit={handleSubmit(onSubmit)}>
        <input {...register("username")} placeholder="enter username"></input>
        <input {...register("password")} placeholder="enter password"></input>
        <input type="submit" />
      </form>
      {response != 0 ? (
        <p>
          {response === 200
            ? "Successfuly signed in"
            : `Error code ${response}`}
        </p>
      ) : null}
    </div>
  );
};

export default page;
