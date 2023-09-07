import type { Metadata } from "next";
import { Inter } from "next/font/google";
import Navbar from "./components/Navbar";
import "./globals.css";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "yonahaforum",
  description: "Demo social media site",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className="w-full h-full bg-gray-300">
      <body className={`${inter.className} w-full h-full`}>
        <Navbar />
        {children}
      </body>
    </html>
  );
}
