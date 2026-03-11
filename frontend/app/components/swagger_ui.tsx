export default function SwaggerUi() {
  const apiBaseUrl =
    process.env.NEXT_PUBLIC_SWAGGER_API_BASE_URL ??
    process.env.NEXT_PUBLIC_API_BASE_URL ??
    "http://localhost:8080/api";
  const iframeSrc = `/swagger_ui.html?apiBaseUrl=${encodeURIComponent(apiBaseUrl)}`;

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
          src={iframeSrc}
          className="h-[80vh] w-full rounded-2xl border border-slate-200 bg-white"
        />
      </div>
    </div>
  );
}
