"use client";

import { ChangeEvent, FormEvent, useEffect } from "react";
import Link from "next/link";
import { Fish, createFish, deleteFish, fetchFishes, uploadFishImage } from "@/app/lib/admin_api";
import { usePatchReducer } from "@/app/hooks/use_patch_reducer";

type FishesPageState = {
  fishes: Fish[];
  loading: boolean;
  errorMessage: string;
  fishName: string;
  fishCategory: string;
  fishDescription: string;
  fishImageUrl: string;
  fishLinkUrl: string;
  uploadingImage: boolean;
};

const initialState: FishesPageState = {
  fishes: [],
  loading: true,
  errorMessage: "",
  fishName: "",
  fishCategory: "",
  fishDescription: "",
  fishImageUrl: "",
  fishLinkUrl: "",
  uploadingImage: false,
};

export default function FishesPage() {
  const [state, patchState] = usePatchReducer<FishesPageState>(initialState);

  const loadFishes = async (): Promise<void> => {
    patchState({ loading: true, errorMessage: "" });
    try {
      const items = await fetchFishes();
      patchState({ fishes: items });
    } catch (error) {
      patchState({
        errorMessage: error instanceof Error ? error.message : "魚一覧の取得に失敗しました",
      });
    } finally {
      patchState({ loading: false });
    }
  };

  useEffect(() => {
    void loadFishes();
  }, []);

  const handleAddFish = async (event: FormEvent<HTMLFormElement>): Promise<void> => {
    event.preventDefault();
    patchState({ errorMessage: "" });

    try {
      const createdFish = await createFish({
        name: state.fishName,
        category: state.fishCategory,
        description: state.fishDescription,
        imageUrl: state.fishImageUrl,
        linkUrl: state.fishLinkUrl,
      });
      patchState((prev) => ({
        fishes: [createdFish, ...prev.fishes],
        fishName: "",
        fishCategory: "",
        fishDescription: "",
        fishImageUrl: "",
        fishLinkUrl: "",
      }));
    } catch (error) {
      patchState({
        errorMessage: error instanceof Error ? error.message : "魚の登録に失敗しました",
      });
    }
  };

  const handleDeleteFish = async (fishId: string): Promise<void> => {
    patchState({ errorMessage: "" });

    try {
      await deleteFish(fishId);
      patchState((prev) => ({
        fishes: prev.fishes.filter((fish) => fish.id !== fishId),
      }));
    } catch (error) {
      patchState({
        errorMessage: error instanceof Error ? error.message : "魚の削除に失敗しました",
      });
    }
  };

  const handleUploadImage = async (event: ChangeEvent<HTMLInputElement>): Promise<void> => {
    const file = event.target.files?.[0];
    if (!file) {
      return;
    }

    patchState({ errorMessage: "", uploadingImage: true });
    try {
      const imageUrl = await uploadFishImage(file);
      patchState({ fishImageUrl: imageUrl });
    } catch (error) {
      patchState({
        errorMessage: error instanceof Error ? error.message : "画像アップロードに失敗しました",
      });
    } finally {
      patchState({ uploadingImage: false });
      event.target.value = "";
    }
  };

  return (
    <div className="min-h-screen bg-slate-100 px-4 py-8 text-slate-900 md:px-10">
      <main className="mx-auto max-w-6xl space-y-8">
        <section className="rounded-2xl bg-white p-6 shadow-sm">
          <div className="flex flex-wrap items-center justify-between gap-3">
            <div>
              <h1 className="text-2xl font-bold">魚データ管理</h1>
              <p className="mt-2 text-sm text-slate-600">魚の登録と削除を行います。</p>
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
                onClick={() => void loadFishes()}
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
            <h2 className="text-lg font-bold">魚を登録</h2>
            <form className="mt-4 space-y-3" onSubmit={(event) => void handleAddFish(event)}>
              <label className="block text-sm font-medium">
                魚名
                <input
                  className="mt-1 w-full rounded-lg border border-slate-300 px-3 py-2"
                  type="text"
                  value={state.fishName}
                  onChange={(event) => patchState({ fishName: event.target.value })}
                  placeholder="例: ヒラメ"
                  required
                />
              </label>
              <label className="block text-sm font-medium">
                カテゴリ
                <input
                  className="mt-1 w-full rounded-lg border border-slate-300 px-3 py-2"
                  type="text"
                  value={state.fishCategory}
                  onChange={(event) => patchState({ fishCategory: event.target.value })}
                  placeholder="例: 白身魚"
                />
              </label>
              <label className="block text-sm font-medium">
                説明
                <textarea
                  className="mt-1 min-h-20 w-full rounded-lg border border-slate-300 px-3 py-2"
                  value={state.fishDescription}
                  onChange={(event) => patchState({ fishDescription: event.target.value })}
                  placeholder="特徴やおすすめの食べ方"
                />
              </label>
              <label className="block text-sm font-medium">
                画像アップロード（Google Photos）
                <input
                  className="mt-1 w-full rounded-lg border border-slate-300 px-3 py-2"
                  type="file"
                  accept="image/*"
                  onChange={(event) => void handleUploadImage(event)}
                  disabled={state.uploadingImage}
                />
              </label>
              <label className="block text-sm font-medium">
                画像URL
                <input
                  className="mt-1 w-full rounded-lg border border-slate-300 px-3 py-2"
                  type="url"
                  value={state.fishImageUrl}
                  onChange={(event) => patchState({ fishImageUrl: event.target.value })}
                  placeholder="Google PhotosのURLまたは画像URL"
                />
              </label>
              <label className="block text-sm font-medium">
                関連リンクURL
                <input
                  className="mt-1 w-full rounded-lg border border-slate-300 px-3 py-2"
                  type="url"
                  value={state.fishLinkUrl}
                  onChange={(event) => patchState({ fishLinkUrl: event.target.value })}
                  placeholder="紹介ページなどのURL"
                />
              </label>
              {state.uploadingImage && (
                <p className="text-sm text-slate-500">画像アップロード中...</p>
              )}
              <button
                className="rounded-lg bg-slate-900 px-4 py-2 text-sm font-semibold text-white hover:bg-slate-700"
                type="submit"
              >
                魚を追加
              </button>
            </form>
          </article>

          <article className="rounded-2xl bg-white p-6 shadow-sm">
            <h2 className="text-lg font-bold">登録済みの魚</h2>
            <ul className="mt-4 space-y-3">
              {state.fishes.length === 0 && (
                <li className="rounded-lg border border-dashed border-slate-300 p-3 text-sm text-slate-500">
                  まだ登録がありません。
                </li>
              )}
              {state.fishes.map((fish) => (
                <li key={fish.id} className="rounded-lg border border-slate-200 p-3 text-sm">
                  <div className="flex items-center justify-between gap-3">
                    <div>
                      <p className="font-semibold">{fish.name}</p>
                      {fish.imageUrl && (
                        <img
                          src={fish.imageUrl}
                          alt={fish.name}
                          className="mt-2 h-24 w-24 rounded object-cover"
                        />
                      )}
                      {fish.category && <p className="text-slate-500">カテゴリ: {fish.category}</p>}
                      {fish.description && <p className="text-slate-600">{fish.description}</p>}
                      {fish.linkUrl && (
                        <a
                          href={fish.linkUrl}
                          target="_blank"
                          rel="noreferrer"
                          className="text-sky-700 underline"
                        >
                          関連リンク
                        </a>
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
        </section>
      </main>
    </div>
  );
}
