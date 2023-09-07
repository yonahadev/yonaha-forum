import Link from "next/link";
import React from "react";
import { AiOutlineSearch } from "react-icons/ai";
import LoginButton from "./LoginButton";

const Navbar = () => {
  return (
    <nav className="w-full h-16 flex justify-between items-center px-16 bg-white">
      <Link className="text-xl font-semibold" href={"/"}>
        yonahaforum
      </Link>
      <Link href={"/users"}>Users</Link>
      <Link href={"/createpost"}>Create a post</Link>
      <div className="w-1/2 h-3/4 bg-transparent bg-gray-300 px-2 flex items-center">
        <AiOutlineSearch size="20" />
        <input className="bg-transparent w-full h-full px-2 outline-none"></input>
      </div>
      <LoginButton />
    </nav>
  );
};

export default Navbar;
