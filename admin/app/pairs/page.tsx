"use client";

import { FormEvent, useEffect, useMemo } from "react";
import Link from "next/link";
import {
  Fish,
  FishPair,
  buildPairKey,
  createPair,
  deletePair,
  fetchFishes,
  fetchPairs,
} from "@/app/lib/admin_api";
import { usePatchReducer } from "@/app/hooks/use_patch_reducer";

type PairsPageState = {
  fishes: Fish[];
  pairs: FishPair[];
  loading: boolean;
  errorMessage: string;
  pairFishA: string;
  pairFishB: string;
  pairScore: number;
  pairMemo: string;
};

const initialState: PairsPageState = {
  fishes: [],
  pairs: [],
  loading: true,
  errorMessage: "",
  pairFishA: "",
  pairFishB: "",
  pairScore: 3,
  pairMemo: "",
};

export default function PairsPage() {
  const [state, patchState] = usePatchReducer<PairsPageState>(initialState);

  const fishMap = useMemo<Map<string, Fish>>(() => {
    return new Map<string, Fish>(state.fishes.map((fish) => [fish.id, fish]));
  }, [state.fishes]);

  const availableForPair = state.fishes.length >= 2;
  const existingPairKeys = useMemo<Set<string>>(() => {
    return new Set<string>(state.pairs.map((pair) => buildPairKey(pair.fishIdA, pair.fishIdB)));
  }, [state.pairs]);
  const isDuplicatePairSelection =
    state.pairFishA !== "" &&
    state.pairFishB !== "" &&
    existingPairKeys.has(buildPairKey(state.pairFishA, state.pairFishB));

  const loadData = async (): Promise<void> => {
    patchState({ loading: true, errorMessage: "" });
    try {
      const [fishItems, pairItems] = await Promise.all([fetchFishes(), fetchPairs()]);
      patchState({ fishes: fishItems, pairs: pairItems });
    } catch (error) {
      patchState({
        errorMessage:
        error instanceof Error
          ? error.message
          : "データ取得に失敗しました。バックエンドをご確認ください。",
      });
    } finally {
      patchState({ loading: false });
    }
  };

  useEffect(() => {
    void loadData();
  }, []);

  const handleAddPair = async (event: FormEvent<HTMLFormElement>): Promise<void> => {
    event.preventDefault();
    patchState({ errorMessage: "" });

    if (state.pairFishA === "" || state.pairFishB === "") {
      patchState({ errorMessage: "魚Aと魚Bを選択してください" });
      return;
    }
    if (state.pairFishA === state.pairFishB) {
      patchState({ errorMessage: "同じ魚同士は登録できません" });
      return;
    }
    if (isDuplicatePairSelection) {
      patchState({ errorMessage: "同じ魚ペアは既に登録されています" });
      return;
    }

    try {
      const createdPair = await createPair({
        fishIdA: state.pairFishA,
        fishIdB: state.pairFishB,
        score: state.pairScore,
        memo: state.pairMemo,
      });
      patchState((prev) => ({
        pairs: [createdPair, ...prev.pairs],
        pairFishA: "",
        pairFishB: "",
        pairScore: 3,
        pairMemo: "",
      }));
    } catch (error) {
      patchState({
        errorMessage: error instanceof Error ? error.message : "相性の登録に失敗しました",
      });
    }
  };

  const handleDeletePair = async (pairId: string): Promise<void> => {
    patchState({ errorMessage: "" });
    try {
      await deletePair(pairId);
      patchState((prev) => ({
        pairs: prev.pairs.filter((pair) => pair.id !== pairId),
      }));
    } catch (error) {
      patchState({
        errorMessage: error instanceof Error ? error.message : "相性の削除に失敗しました",
      });
    }
  };

  return (
    <div className="min-h-screen bg-slate-100 px-4 py-8 text-slate-900 md:px-10">
      <main className="mx-auto max-w-6xl space-y-8">
        <section className="rounded-2xl bg-white p-6 shadow-sm">
          <div className="flex flex-wrap items-center justify-between gap-3">
            <div>
              <h1 className="text-2xl font-bold">魚相性管理</h1>
              <p className="mt-2 text-sm text-slate-600">魚同士の相性登録と削除を行います。</p>
            </div>
            <div className="flex items-center gap-2">
              <Link
                href="/"
                className="rounded-lg border border-slate-300 px-3 py-2 text-sm font-semibold hover:bg-slate-100"
              >
                管理トップへ
              </Link>
              <button
                type="button"
                className="rounded-lg border border-slate-300 px-3 py-2 text-sm font-semibold hover:bg-slate-100"
                onClick={() => void loadData()}
                disabled={state.loading}
              >
                再読み込み
              </button>
            </div>
          </div>
          {state.loading && <p className="mt-3 text-sm text-slate-500">読み込み中...</p>}
          {state.errorMessage && (
            <p className="mt-3 rounded-lg bg-rose-50 px-3 py-2 text-sm text-rose-700">
              {state.errorMessage}
            </p>
          )}
        </section>

        <section className="grid gap-6 lg:grid-cols-2">
          <article className="rounded-2xl bg-white p-6 shadow-sm">
            <h2 className="text-lg font-bold">魚同士の相性を登録</h2>
            <form className="mt-4 space-y-3" onSubmit={(event) => void handleAddPair(event)}>
              <label className="block text-sm font-medium">
                魚 A
                <select
                  className="mt-1 w-full rounded-lg border border-slate-300 px-3 py-2"
                  value={state.pairFishA}
                  onChange={(event) => patchState({ pairFishA: event.target.value })}
                  disabled={!availableForPair}
                  required
                >
                  <option value="">選択してください</option>
                  {state.fishes.map((fish) => (
                    <option
                      key={fish.id}
                      value={fish.id}
                      disabled={
                        state.pairFishB !== "" &&
                        fish.id !== state.pairFishB &&
                        existingPairKeys.has(buildPairKey(fish.id, state.pairFishB))
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
                  value={state.pairFishB}
                  onChange={(event) => patchState({ pairFishB: event.target.value })}
                  disabled={!availableForPair}
                  required
                >
                  <option value="">選択してください</option>
                  {state.fishes.map((fish) => (
                    <option
                      key={fish.id}
                      value={fish.id}
                      disabled={
                        fish.id === state.pairFishA ||
                        (state.pairFishA !== "" &&
                          fish.id !== state.pairFishA &&
                          existingPairKeys.has(buildPairKey(state.pairFishA, fish.id)))
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
                  value={state.pairScore}
                  onChange={(event) => patchState({ pairScore: Number(event.target.value) })}
                  required
                />
              </label>
              <label className="block text-sm font-medium">
                メモ
                <textarea
                  className="mt-1 min-h-20 w-full rounded-lg border border-slate-300 px-3 py-2"
                  value={state.pairMemo}
                  onChange={(event) => patchState({ pairMemo: event.target.value })}
                  placeholder="相性の理由など"
                />
              </label>
              <button
                className="rounded-lg bg-emerald-700 px-4 py-2 text-sm font-semibold text-white hover:bg-emerald-600 disabled:bg-slate-400"
                type="submit"
                disabled={
                  !availableForPair ||
                  state.pairFishA === state.pairFishB ||
                  isDuplicatePairSelection
                }
              >
                相性を追加
              </button>
              {!availableForPair && (
                <p className="text-sm text-amber-700">相性登録には2件以上の魚登録が必要です。</p>
              )}
              {isDuplicatePairSelection && (
                <p className="text-sm text-amber-700">
                  同じ組み合わせの魚ペアは既に登録されています。
                </p>
              )}
            </form>
          </article>

          <article className="rounded-2xl bg-white p-6 shadow-sm">
            <h2 className="text-lg font-bold">登録済みの相性</h2>
            <ul className="mt-4 space-y-3">
              {state.pairs.length === 0 && (
                <li className="rounded-lg border border-dashed border-slate-300 p-3 text-sm text-slate-500">
                  まだ登録がありません。
                </li>
              )}
              {state.pairs.map((pair) => {
                const fishA = fishMap.get(pair.fishIdA)?.name ?? "削除された魚";
                const fishB = fishMap.get(pair.fishIdB)?.name ?? "削除された魚";

                return (
                  <li key={pair.id} className="rounded-lg border border-slate-200 p-3 text-sm">
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
