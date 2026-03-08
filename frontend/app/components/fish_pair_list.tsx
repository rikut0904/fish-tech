"use client";

import { useEffect, useMemo, useState } from "react";

type Fish = {
  id: string;
  name: string;
  category: string;
  description: string;
};

type FishPair = {
  id: string;
  fishIdA: string;
  fishIdB: string;
  score: number;
  memo: string;
};

type ListResponse<T> = {
  items: T[];
};

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080/api";

async function fetchList<T>(path: string): Promise<T[]> {
  const response = await fetch(`${API_BASE_URL}${path}`, { cache: "no-store" });
  if (!response.ok) {
    throw new Error("データ取得に失敗しました");
  }
  const data = (await response.json()) as ListResponse<T>;
  return data.items;
}

export default function FishPairList() {
  const [fishes, setFishes] = useState<Fish[]>([]);
  const [pairs, setPairs] = useState<FishPair[]>([]);
  const [error, setError] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const load = async (): Promise<void> => {
      setLoading(true);
      setError("");

      try {
        const [fishItems, pairItems] = await Promise.all([
          fetchList<Fish>("/fishes"),
          fetchList<FishPair>("/pairs"),
        ]);
        setFishes(fishItems);
        setPairs(pairItems);
      } catch {
        setError("相性一覧の取得に失敗しました");
      } finally {
        setLoading(false);
      }
    };

    void load();
  }, []);

  const fishNameById = useMemo<Map<string, string>>(() => {
    return new Map<string, string>(fishes.map((fish) => [fish.id, fish.name]));
  }, [fishes]);

  if (loading) {
    return <p className="text-sm text-blue-700">読み込み中...</p>;
  }

  if (error) {
    return <p className="text-sm text-rose-700">{error}</p>;
  }

  if (pairs.length === 0) {
    return <p className="text-sm text-slate-600">相性データはまだありません。</p>;
  }

  return (
    <ul className="space-y-3">
      {pairs.map((pair) => {
        const fishA = fishNameById.get(pair.fishIdA) ?? "不明な魚";
        const fishB = fishNameById.get(pair.fishIdB) ?? "不明な魚";

        return (
          <li key={pair.id} className="rounded-xl border border-blue-100 bg-white p-4 shadow-sm">
            <p className="font-semibold text-slate-800">
              {fishA} × {fishB}
            </p>
            <p className="text-sm text-slate-600">スコア: {pair.score}</p>
            {pair.memo && <p className="mt-1 text-sm text-slate-600">{pair.memo}</p>}
          </li>
        );
      })}
    </ul>
  );
}
