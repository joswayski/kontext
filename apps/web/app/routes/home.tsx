import type { Route } from "./+types/home";
import { Welcome } from "../welcome/welcome";
import { useLoaderData } from "react-router";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Kontext by Jose Valerio" },
    { name: "description", content: "Kafka visualization and business logic mapping" },
  ];
}

const apiUrl = process.env.API_URL || "http://localhost:4000";
export async function loader() {


  console.log(`SENDING REQUEST TO ${apiUrl}`)
  const res = await fetch(`${apiUrl}/`)
  if (!res.ok) {
    throw new Error("Failed to fetch from API")
  }
  try {
    const data = await res.json()
    return {
      message: data.message,
    }
  } catch (error) {
    console.error(error)
    throw new Error("Failed to parse response from API")
  }
}
export default function Home() {
  const { message } = useLoaderData<typeof loader>()
  return <div>
    <Welcome apiStatus={message}/>
  </div>;
}
