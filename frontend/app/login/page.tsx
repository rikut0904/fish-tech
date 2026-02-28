"use client";
// ログイン画面

import Link from "next/link";
import Header from "@/app/components/header";
import Footer from "@/app/components/footer";
import { useState } from "react";
import { useRouter } from "next/navigation";


export default function LoginPage() {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [error, setError] = useState("");
    const router = useRouter();

    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        const envEmail = process.env.NEXT_PUBLIC_LOGIN_EMAIL;
        const envPassword = process.env.NEXT_PUBLIC_LOGIN_PASSWORD;
        if (email === envEmail && password === envPassword) {
            setError("");
            router.push("/");
        } else {
            setError("メールアドレスまたはパスワードが正しくありません。");
        }
    };

    return (
        <div className="flex flex-col min-h-screen bg-blue-50">
            <Header />
            <main className="flex-1 flex items-center justify-center">
                <div className="bg-white p-8 rounded shadow-md w-full max-w-md">
                    <h1 className="text-2xl font-bold mb-6 text-center">ログイン</h1>
                    <form className="space-y-4" onSubmit={handleSubmit}>
                        <div>
                            <label htmlFor="email" className="block text-sm font-medium mb-1">メールアドレス</label>
                            <input
                                type="email"
                                id="email"
                                name="email"
                                className="w-full border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-400"
                                required
                                value={email}
                                onChange={e => setEmail(e.target.value)}
                            />
                        </div>
                        <div>
                            <label htmlFor="password" className="block text-sm font-medium mb-1">パスワード</label>
                            <input
                                type="password"
                                id="password"
                                name="password"
                                className="w-full border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-400"
                                required
                                value={password}
                                onChange={e => setPassword(e.target.value)}
                            />
                        </div>
                        {error && <p className="text-red-500 text-sm text-center">{error}</p>}
                        <button
                            type="submit"
                            className="w-full bg-blue-500 text-white py-2 rounded hover:bg-blue-600 font-semibold transition"
                        >
                            ログイン
                        </button>
                    </form>
                    <p className="mt-4 text-center text-sm">
                        アカウントをお持ちでない方は{' '}
                        <Link href="/register" className="text-blue-500 hover:underline">新規登録</Link>
                    </p>
                </div>
            </main>
            <Footer />
        </div>
    );
}
