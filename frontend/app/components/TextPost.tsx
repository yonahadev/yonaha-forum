"use client";
import { ApiResponse } from "@/interfaces";
import { useRouter } from "next/navigation";
import React, { useEffect, useState } from "react";
import { FieldValues, set, useForm } from "react-hook-form";
import useAuth from "../hooks/useAuth";

interface Props {
  title: string;
  content: string;
  username: string;
  postID: number;
}

const TextPost = ({ title, content, username, postID }: Props) => {
  useEffect(() => {
    setPostText({ title: title, text_content: content });
  }, []);

  const router = useRouter();
  const [response, setResponse] = useState({ error: "", message: "" });
  const { register, handleSubmit } = useForm();

  const savePost = async (data: FieldValues) => {
    data.id = postID;
    const res = await fetch("http://127.0.0.1:8080/posts", {
      method: "PATCH",
      body: JSON.stringify(data),
      credentials: "include",
      headers: { "Content-Type": "application/json" },
    });
    const responseData: ApiResponse = await res.json();
    if (res.ok) {
      setPostText({ title: data.title, text_content: data.text_content });
      setEditing(false);
      router.refresh();
    } else {
      setResponse({ error: responseData.error!, message: "" });
    }
  };

  const deletePost = async (postID: number) => {
    const res = await fetch("http://127.0.0.1:8080/posts", {
      method: "DELETE",
      credentials: "include",
      body: JSON.stringify({ id: postID }),
      headers: {
        "Content-Type": "application/json",
      },
    });
    if (res.ok) {
      router.refresh();
    }
  };

  const { authenticated, user } = useAuth();
  const [editing, setEditing] = useState(false);
  const [postText, setPostText] = useState({ title: "", text_content: "" });

  return (
    <div className="w-1/2 bg-white p-4 h-fit mt-2">
      {editing === false ? (
        <>
          <p className="text-2xl">{title}</p>
          <p className="opacity-50">{username}</p>
          <p>{content}</p>
        </>
      ) : (
        <form onSubmit={handleSubmit(savePost)}>
          <input {...register("title")} defaultValue={postText.title}></input>
          <input
            {...register("text_content")}
            defaultValue={postText.text_content}
          ></input>
          <p>{response?.error}</p>
          <p>{response?.message}</p>
          <button type="submit">Save</button>
        </form>
      )}
      {authenticated && username === user?.username ? (
        <div className="text-xl">
          {editing === false ? (
            <>
              <button
                className="text-red-800  mr-2"
                onClick={() => deletePost(postID)}
              >
                Delete
              </button>
              <button
                className="text-blue-800 "
                onClick={() => setEditing(true)}
              >
                Update
              </button>
            </>
          ) : null}
        </div>
      ) : null}
    </div>
  );
};

export default TextPost;
