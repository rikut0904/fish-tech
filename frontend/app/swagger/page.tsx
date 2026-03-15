import type { Metadata } from "next";
import SwaggerUi from "@/app/components/swagger_ui";

export const metadata: Metadata = {
  title: "Swagger UI | Fish-Tech",
  description: "Fish-Tech API の Swagger UI",
};

export default function SwaggerPage() {
  return <SwaggerUi />;
}
