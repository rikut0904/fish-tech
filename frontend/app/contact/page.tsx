"use client"
import { useState } from "react";
import Header from "@/app/components/header";

export default function ContactPage() {
    const [name, setName] = useState("");
    const [email, setEmail] = useState("");
    const [message, setMessage] = useState("");
    const [submitted, setSubmitted] = useState(false);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        // 本来は API に送信するが、まだ未実装のため擬似送信とする
        setSubmitted(true);
        setName("");
        setEmail("");
        setMessage("");
    };

    return (
        <>
            <Header />
            <main className="container mx-auto p-6">
                <h1 className="text-2xl font-bold mb-4">お問い合わせ</h1>

                {submitted && (
                    <div className="bg-green-50 border border-green-200 text-green-800 p-4 mb-4 rounded">
                        送信が完了しました。返信をお待ちください。
                    </div>
                )}

                <form onSubmit={handleSubmit} className="space-y-4 max-w-lg">
                    <div>
                        <label className="block mb-1 font-medium">お名前</label>
                        <input
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            required
                            className="w-full border rounded px-3 py-2"
                        />
                    </div>

                    <div>
                        <label className="block mb-1 font-medium">メールアドレス</label>
                        <input
                            type="email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            required
                            className="w-full border rounded px-3 py-2"
                        />
                    </div>

                    <div>
                        <label className="block mb-1 font-medium">お問い合わせ内容</label>
                        <textarea
                            value={message}
                            onChange={(e) => setMessage(e.target.value)}
                            required
                            className="w-full border rounded px-3 py-2 h-32"
                        />
                    </div>

                    <div>
                        <button type="submit" className="bg-blue-700 text-white px-4 py-2 rounded hover:bg-blue-800">
                            送信
                        </button>
                    </div>
                </form>
            </main>
        </>
    );
}
