"use client";
import { useRouter } from "next/navigation";
import { ArrowLeft, Plus, Edit, Trash2, User as UserIcon } from "lucide-react";
import { useEffect } from "react";
import Header from "@/app/components/header";
import Footer from "@/app/components/footer";
// 仮データ
const diaryData = [
    {
        id: 1,
        fishName: "アジ",
        date: "2026-02-28",
        location: "金沢港",
        comment: "脂がのっていて美味しかった！",
        imageUrl: "https://images.unsplash.com/photo-1504674900247-0877df9cc836?auto=format&fit=crop&w=200&q=80"
    },
    {
        id: 2,
        fishName: "カレイ",
        date: "2026-02-15",
        location: "近江町市場",
        comment: "煮付けが最高。",
        imageUrl: "https://images.unsplash.com/photo-1465101046530-73398c7f28ca?auto=format&fit=crop&w=200&q=80"
    }
];

export default function MyPage() {
    const router = useRouter();

    useEffect(() => {
        if (typeof window !== "undefined") {
            const isLoggedIn = localStorage.getItem("isLoggedIn");
            if (!isLoggedIn) {
                router.push("/login");
            }
        }
    }, [router]);

    return (
        <div className="flex flex-col min-h-screen bg-gray-50 pb-20">
            <Header />
            <main className="flex-1 px-4 py-6">
                {/* ユーザー情報 */}
                <div className="bg-white rounded-lg shadow-sm p-6 mb-6">
                    <div className="flex items-center gap-4 mb-4">
                        <div className="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center">
                            <UserIcon className="w-8 h-8 text-blue-600" />
                        </div>
                        <div>
                            <h2 className="text-lg font-medium">ユーザー名 さん</h2>
                            <p className="text-sm text-gray-600">食べた魚: {diaryData.length}種類</p>
                        </div>
                    </div>
                </div>
                {/* 日記リスト */}
                <div className="space-y-3 mb-6">
                    {diaryData.map((diary) => (
                        <div key={diary.id} className="bg-white rounded-lg shadow-sm p-4">
                            <div className="flex gap-3">
                                {diary.imageUrl && (
                                    <div className="w-20 h-20 flex-shrink-0 rounded-lg overflow-hidden">
                                        <img
                                            src={diary.imageUrl}
                                            alt={diary.fishName}
                                            className="w-full h-full object-cover"
                                        />
                                    </div>
                                )}
                                <div className="flex-1 min-w-0">
                                    <div className="flex items-start justify-between mb-2">
                                        <h3 className="font-medium text-lg">{diary.fishName}</h3>
                                        <div className="flex gap-2">
                                            <button className="p-1 hover:bg-gray-100 rounded">
                                                <Edit className="w-4 h-4 text-gray-600" />
                                            </button>
                                            <button className="p-1 hover:bg-gray-100 rounded">
                                                <Trash2 className="w-4 h-4 text-gray-600" />
                                            </button>
                                        </div>
                                    </div>
                                    <p className="text-sm text-gray-600 mb-1">
                                        {diary.date} {diary.location}
                                    </p>
                                    {diary.comment && (
                                        <p className="text-sm text-gray-700">「{diary.comment}」</p>
                                    )}
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
                {/* 追加ボタン */}
                <button
                    onClick={() => router.push("/mypage/diaries/new")}
                    className="w-full flex items-center justify-center bg-blue-500 text-white py-3 rounded-lg font-semibold text-lg hover:bg-blue-600 transition"
                >
                    <Plus className="w-5 h-5 mr-2" />
                    新しい記録を追加
                </button>
            </main>
            <Footer />
        </div>
    );
}
