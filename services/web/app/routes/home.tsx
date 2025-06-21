import type { Route } from "./+types/home";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Kontext by Jose Valerio" },
    {
      name: "description",
      content:
        "Automated Kafka event flow visualization and business logic mapping",
    },
  ];
}

export default function Home() {
  return (
    <main className="flex items-center justify-center pt-16 pb-4">
      <div className="flex-1 flex flex-col items-center gap-16 min-h-0">
        <div className="max-w-[300px] w-full space-y-6 px-4">
          <h1>Welcome to Kontext</h1>
        </div>
      </div>
    </main>
  );
}
