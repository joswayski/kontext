import type { Route } from "./+types/home";
import { Welcome } from "../welcome/welcome";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Kontext" },
    {
      name: "description",
      content:
        "Automated Kafka event flow visualization and business logic mapping",
    },
  ];
}

export default function Home() {
  return <Welcome />;
}
