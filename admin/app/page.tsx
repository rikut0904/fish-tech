import Link from "next/link";

export default function AdminPage() {
  return (
    <div className="min-h-screen bg-slate-100 px-4 py-8 text-slate-900 md:px-10">
      <main className="mx-auto max-w-5xl space-y-8">
        <section className="rounded-2xl bg-white p-6 shadow-sm">
          <h1 className="text-2xl font-bold">管理画面</h1>
        </section>

        <section className="grid gap-6 md:grid-cols-2">
          <article className="rounded-2xl bg-white p-6 shadow-sm">
            <h2 className="text-lg font-bold">魚データ管理</h2>
            <p className="mt-2 text-sm text-slate-600">
              魚の登録と削除を行います。
            </p>
            <Link
              href="/fishes"
              className="mt-4 inline-block rounded-lg bg-slate-900 px-4 py-2 text-sm font-semibold text-white hover:bg-slate-700"
            >
              魚管理ページへ
            </Link>
          </article>

          <article className="rounded-2xl bg-white p-6 shadow-sm">
            <h2 className="text-lg font-bold">魚相性管理</h2>
            <p className="mt-2 text-sm text-slate-600">
              魚同士の相性登録と削除を行います。
            </p>
            <Link
              href="/pairs"
              className="mt-4 inline-block rounded-lg bg-emerald-700 px-4 py-2 text-sm font-semibold text-white hover:bg-emerald-600"
            >
              相性管理ページへ
            </Link>
          </article>
        </section>
      </main>
    </div>
  );
}
