"use client";
import { FieldValues, useForm } from "react-hook-form";

const page = () => {
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
    console.log(`${result} account added`);
  };

  return (
    <div>
      Create an account!
      <form onSubmit={handleSubmit(onSubmit)}>
        <input {...register("username")} placeholder="enter username"></input>
        <input type="submit" />
      </form>
    </div>
  );
};

export default page;
