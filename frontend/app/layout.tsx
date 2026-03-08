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
    url: "https://fish-tech.example",
    images: [
      {
        url: "/fishtech-og.webp",
        width: 1200,
        height: 630,
        alt: "Fish-Tech — 金沢の魚",
      },
    ],
  },
  twitter: {
    card: "summary_large_image",
    title: "Fish-Tech",
    description: "金沢の魚を紹介するアプリ（fish-tech）",
    images: ["/fishtech-og.webp"],
  },
  icons: {
    icon: "/fishtech-favicon.webp",
    shortcut: "/fishtech-favicon.webp",
    apple: "/fishtech-favicon.webp",
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
