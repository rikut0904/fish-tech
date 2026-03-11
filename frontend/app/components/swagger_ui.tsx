export default function SwaggerUi() {
  return (
    <div className="min-h-screen bg-slate-100 px-4 py-8 text-slate-900">
      <div className="mx-auto max-w-6xl rounded-3xl bg-white p-4 shadow-sm md:p-6">
        <div className="mb-4 border-b border-slate-200 pb-4">
          <h1 className="text-2xl font-bold">API ドキュメント</h1>
          <p className="mt-2 text-sm text-slate-600">
            Fish-Tech の OpenAPI をブラウザ上で確認できます。
          </p>
        </div>
        <iframe
          title="Fish-Tech Swagger UI"
          src="/swagger_ui.html"
          className="h-[80vh] w-full rounded-2xl border border-slate-200 bg-white"
        />
      </div>
    </div>
  );
}
