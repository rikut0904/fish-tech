import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import Header from "@/app/components/header";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Fish-Tech",
  description: "金沢の魚を紹介するアプリ（fish-tech）",
  keywords: [
    "金沢",
    "魚",
    "地魚",
    "魚図鑑",
    "レシピ",
    "漁業",
    "サステナブル",
    "fish-tech",
  ],
  openGraph: {
    title: "Fish-Tech",
    description: "金沢の魚を紹介するアプリ（fish-tech）",
    siteName: "Fish-Tech",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ja">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <Header />
        {children}
      </body>
    </html>
  );
}
