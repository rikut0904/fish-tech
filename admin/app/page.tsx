"use client";

import { FormEvent, useEffect, useMemo, useState } from "react";

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

function buildPairKey(fishIdA: string, fishIdB: string): string {
  return [fishIdA, fishIdB].sort().join(":");
}

async function parseErrorMessage(response: Response): Promise<string> {
  try {
    const data = (await response.json()) as { error?: string };
    return data.error ?? "APIエラーが発生しました";
  } catch {
    return "APIエラーが発生しました";
  }
}

export default function AdminPage() {
  const [fishes, setFishes] = useState<Fish[]>([]);
  const [pairs, setPairs] = useState<FishPair[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [errorMessage, setErrorMessage] = useState<string>("");

  const [fishName, setFishName] = useState<string>("");
  const [fishCategory, setFishCategory] = useState<string>("");
  const [fishDescription, setFishDescription] = useState<string>("");

  const [pairFishA, setPairFishA] = useState<string>("");
  const [pairFishB, setPairFishB] = useState<string>("");
  const [pairScore, setPairScore] = useState<number>(3);
  const [pairMemo, setPairMemo] = useState<string>("");

  const fishMap = useMemo<Map<string, Fish>>(() => {
    return new Map<string, Fish>(fishes.map((fish) => [fish.id, fish]));
  }, [fishes]);

  const availableForPair = fishes.length >= 2;
  const existingPairKeys = useMemo<Set<string>>(() => {
    return new Set<string>(pairs.map((pair) => buildPairKey(pair.fishIdA, pair.fishIdB)));
  }, [pairs]);
  const isDuplicatePairSelection =
    pairFishA !== "" &&
    pairFishB !== "" &&
    existingPairKeys.has(buildPairKey(pairFishA, pairFishB));

  const loadData = async (): Promise<void> => {
    setLoading(true);
    setErrorMessage("");

    try {
      const [fishesResponse, pairsResponse] = await Promise.all([
        fetch(`${API_BASE_URL}/fishes`, { cache: "no-store" }),
        fetch(`${API_BASE_URL}/pairs`, { cache: "no-store" }),
      ]);

      if (!fishesResponse.ok) {
        throw new Error(await parseErrorMessage(fishesResponse));
      }
      if (!pairsResponse.ok) {
        throw new Error(await parseErrorMessage(pairsResponse));
      }

      const fishesData = (await fishesResponse.json()) as ListResponse<Fish>;
      const pairsData = (await pairsResponse.json()) as ListResponse<FishPair>;

      setFishes(fishesData.items);
      setPairs(pairsData.items);
    } catch (error) {
      setErrorMessage(
        error instanceof Error
          ? error.message
          : "データ取得に失敗しました。バックエンドをご確認ください。",
      );
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void loadData();
  }, []);

  const handleAddFish = async (event: FormEvent<HTMLFormElement>): Promise<void> => {
    event.preventDefault();
    setErrorMessage("");

    try {
      const response = await fetch(`${API_BASE_URL}/admin/fishes`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          name: fishName,
          category: fishCategory,
          description: fishDescription,
        }),
      });

      if (!response.ok) {
        throw new Error(await parseErrorMessage(response));
      }

      const createdFish = (await response.json()) as Fish;
      setFishes((prev) => [createdFish, ...prev]);
      setFishName("");
      setFishCategory("");
      setFishDescription("");
    } catch (error) {
      setErrorMessage(
        error instanceof Error ? error.message : "魚の登録に失敗しました",
      );
    }
  };

  const handleDeleteFish = async (fishId: string): Promise<void> => {
    setErrorMessage("");

    try {
      const response = await fetch(`${API_BASE_URL}/admin/fishes/${fishId}`, {
        method: "DELETE",
      });

      if (!response.ok) {
        throw new Error(await parseErrorMessage(response));
      }

      setFishes((prev) => prev.filter((fish) => fish.id !== fishId));
      setPairs((prev) =>
        prev.filter((pair) => pair.fishIdA !== fishId && pair.fishIdB !== fishId),
      );
      if (pairFishA === fishId) {
        setPairFishA("");
      }
      if (pairFishB === fishId) {
        setPairFishB("");
      }
    } catch (error) {
      setErrorMessage(
        error instanceof Error ? error.message : "魚の削除に失敗しました",
      );
    }
  };

  const handleAddPair = async (event: FormEvent<HTMLFormElement>): Promise<void> => {
    event.preventDefault();
    setErrorMessage("");

    if (pairFishA === "" || pairFishB === "") {
      setErrorMessage("魚Aと魚Bを選択してください");
      return;
    }
    if (pairFishA === pairFishB) {
      setErrorMessage("同じ魚同士は登録できません");
      return;
    }
    if (isDuplicatePairSelection) {
      setErrorMessage("同じ魚ペアは既に登録されています");
      return;
    }

    try {
      const response = await fetch(`${API_BASE_URL}/admin/pairs`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          fishIdA: pairFishA,
          fishIdB: pairFishB,
          score: pairScore,
          memo: pairMemo,
        }),
      });

      if (!response.ok) {
        throw new Error(await parseErrorMessage(response));
      }

      const createdPair = (await response.json()) as FishPair;
      setPairs((prev) => [createdPair, ...prev]);
      setPairFishA("");
      setPairFishB("");
      setPairScore(3);
      setPairMemo("");
    } catch (error) {
      setErrorMessage(
        error instanceof Error ? error.message : "相性の登録に失敗しました",
      );
    }
  };

  const handleDeletePair = async (pairId: string): Promise<void> => {
    setErrorMessage("");

    try {
      const response = await fetch(`${API_BASE_URL}/admin/pairs/${pairId}`, {
        method: "DELETE",
      });

      if (!response.ok) {
        throw new Error(await parseErrorMessage(response));
      }

      setPairs((prev) => prev.filter((pair) => pair.id !== pairId));
    } catch (error) {
      setErrorMessage(
        error instanceof Error ? error.message : "相性の削除に失敗しました",
      );
    }
  };

  return (
    <div className="min-h-screen bg-slate-100 px-4 py-8 text-slate-900 md:px-10">
      <main className="mx-auto max-w-6xl space-y-8">
        <section className="rounded-2xl bg-white p-6 shadow-sm">
          <h1 className="text-2xl font-bold">管理画面</h1>
          <p className="mt-2 text-sm text-slate-600">
            バックエンドAPI経由で魚と魚同士の相性データを管理できます。
          </p>
          <div className="mt-3 flex items-center gap-3">
            <button
              type="button"
              className="rounded-lg border border-slate-300 px-3 py-2 text-sm font-semibold hover:bg-slate-100"
              onClick={() => void loadData()}
              disabled={loading}
            >
              再読み込み
            </button>
            {loading && <p className="text-sm text-slate-500">読み込み中...</p>}
          </div>
          {errorMessage && (
            <p className="mt-3 rounded-lg bg-rose-50 px-3 py-2 text-sm text-rose-700">
              {errorMessage}
            </p>
          )}
        </section>

        <section className="grid gap-6 lg:grid-cols-2">
          <article className="rounded-2xl bg-white p-6 shadow-sm">
            <h2 className="text-lg font-bold">魚を登録</h2>
            <form className="mt-4 space-y-3" onSubmit={(event) => void handleAddFish(event)}>
              <label className="block text-sm font-medium">
                魚名
                <input
                  className="mt-1 w-full rounded-lg border border-slate-300 px-3 py-2"
                  type="text"
                  value={fishName}
                  onChange={(event) => setFishName(event.target.value)}
                  placeholder="例: ヒラメ"
                  required
                />
              </label>
              <label className="block text-sm font-medium">
                カテゴリ
                <input
                  className="mt-1 w-full rounded-lg border border-slate-300 px-3 py-2"
                  type="text"
                  value={fishCategory}
                  onChange={(event) => setFishCategory(event.target.value)}
                  placeholder="例: 白身魚"
                />
              </label>
              <label className="block text-sm font-medium">
                説明
                <textarea
                  className="mt-1 min-h-20 w-full rounded-lg border border-slate-300 px-3 py-2"
                  value={fishDescription}
                  onChange={(event) => setFishDescription(event.target.value)}
                  placeholder="特徴やおすすめの食べ方"
                />
              </label>
              <button
                className="rounded-lg bg-slate-900 px-4 py-2 text-sm font-semibold text-white hover:bg-slate-700"
                type="submit"
              >
                魚を追加
              </button>
            </form>
          </article>

          <article className="rounded-2xl bg-white p-6 shadow-sm">
            <h2 className="text-lg font-bold">魚同士の相性を登録</h2>
            <form className="mt-4 space-y-3" onSubmit={(event) => void handleAddPair(event)}>
              <label className="block text-sm font-medium">
                魚 A
                <select
                  className="mt-1 w-full rounded-lg border border-slate-300 px-3 py-2"
                  value={pairFishA}
                  onChange={(event) => setPairFishA(event.target.value)}
                  disabled={!availableForPair}
                  required
                >
                  <option value="">選択してください</option>
                  {fishes.map((fish) => (
                    <option
                      key={fish.id}
                      value={fish.id}
                      disabled={
                        pairFishB !== "" &&
                        fish.id !== pairFishB &&
                        existingPairKeys.has(buildPairKey(fish.id, pairFishB))
                      }
                    >
                      {fish.name}
                    </option>
                  ))}
                </select>
              </label>
              <label className="block text-sm font-medium">
                魚 B
                <select
                  className="mt-1 w-full rounded-lg border border-slate-300 px-3 py-2"
                  value={pairFishB}
                  onChange={(event) => setPairFishB(event.target.value)}
                  disabled={!availableForPair}
                  required
                >
                  <option value="">選択してください</option>
                  {fishes.map((fish) => (
                    <option
                      key={fish.id}
                      value={fish.id}
                      disabled={
                        fish.id === pairFishA ||
                        (pairFishA !== "" &&
                          fish.id !== pairFishA &&
                          existingPairKeys.has(buildPairKey(pairFishA, fish.id)))
                      }
                    >
                      {fish.name}
                    </option>
                  ))}
                </select>
              </label>
              <label className="block text-sm font-medium">
                相性スコア (1〜5)
                <input
                  className="mt-1 w-full rounded-lg border border-slate-300 px-3 py-2"
                  type="number"
                  min={1}
                  max={5}
                  value={pairScore}
                  onChange={(event) => setPairScore(Number(event.target.value))}
                  required
                />
              </label>
              <label className="block text-sm font-medium">
                メモ
                <textarea
                  className="mt-1 min-h-20 w-full rounded-lg border border-slate-300 px-3 py-2"
                  value={pairMemo}
                  onChange={(event) => setPairMemo(event.target.value)}
                  placeholder="相性の理由など"
                />
              </label>
              <button
                className="rounded-lg bg-emerald-700 px-4 py-2 text-sm font-semibold text-white hover:bg-emerald-600 disabled:bg-slate-400"
                type="submit"
                disabled={!availableForPair || pairFishA === pairFishB || isDuplicatePairSelection}
              >
                相性を追加
              </button>
              {!availableForPair && (
                <p className="text-sm text-amber-700">
                  相性登録には2件以上の魚登録が必要です。
                </p>
              )}
              {isDuplicatePairSelection && (
                <p className="text-sm text-amber-700">
                  同じ組み合わせの魚ペアは既に登録されています。
                </p>
              )}
            </form>
          </article>
        </section>

        <section className="grid gap-6 lg:grid-cols-2">
          <article className="rounded-2xl bg-white p-6 shadow-sm">
            <h2 className="text-lg font-bold">登録済みの魚</h2>
            <ul className="mt-4 space-y-3">
              {fishes.length === 0 && (
                <li className="rounded-lg border border-dashed border-slate-300 p-3 text-sm text-slate-500">
                  まだ登録がありません。
                </li>
              )}
              {fishes.map((fish) => (
                <li
                  key={fish.id}
                  className="rounded-lg border border-slate-200 p-3 text-sm"
                >
                  <div className="flex items-center justify-between gap-3">
                    <div>
                      <p className="font-semibold">{fish.name}</p>
                      {fish.category && (
                        <p className="text-slate-500">カテゴリ: {fish.category}</p>
                      )}
                      {fish.description && (
                        <p className="text-slate-600">{fish.description}</p>
                      )}
                    </div>
                    <button
                      className="rounded-md border border-rose-400 px-3 py-1 text-xs font-semibold text-rose-700 hover:bg-rose-50"
                      type="button"
                      onClick={() => void handleDeleteFish(fish.id)}
                    >
                      削除
                    </button>
                  </div>
                </li>
              ))}
            </ul>
          </article>

          <article className="rounded-2xl bg-white p-6 shadow-sm">
            <h2 className="text-lg font-bold">登録済みの相性</h2>
            <ul className="mt-4 space-y-3">
              {pairs.length === 0 && (
                <li className="rounded-lg border border-dashed border-slate-300 p-3 text-sm text-slate-500">
                  まだ登録がありません。
                </li>
              )}
              {pairs.map((pair) => {
                const fishA = fishMap.get(pair.fishIdA)?.name ?? "削除された魚";
                const fishB = fishMap.get(pair.fishIdB)?.name ?? "削除された魚";

                return (
                  <li
                    key={pair.id}
                    className="rounded-lg border border-slate-200 p-3 text-sm"
                  >
                    <div className="flex items-start justify-between gap-3">
                      <div>
                        <p className="font-semibold">
                          {fishA} × {fishB}
                        </p>
                        <p className="text-slate-600">スコア: {pair.score}</p>
                        {pair.memo && <p className="text-slate-600">{pair.memo}</p>}
                      </div>
                      <button
                        className="rounded-md border border-rose-400 px-3 py-1 text-xs font-semibold text-rose-700 hover:bg-rose-50"
                        type="button"
                        onClick={() => void handleDeletePair(pair.id)}
                      >
                        削除
                      </button>
                    </div>
                  </li>
                );
              })}
            </ul>
          </article>
        </section>
      </main>
    </div>
  );
}
