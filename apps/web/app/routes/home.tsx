import type { Route } from "./+types/home";
import { useLoaderData } from "react-router";

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

type ClusterMetadata = {
  id: string;
  status: "connected" | "error";
  message: string;
  broker_count: number;
  topic_count: number;
  brokers: string[];
  total_size: number;
};
type ClustersResponse = {
  clusters: ClusterMetadata[];
  cluster_count: number;
};
const apiUrl = process.env.API_URL || "http://localhost:4000";
export async function loader() {
  console.log(`SENDING REQUEST TO ${apiUrl}`);
  const res = await fetch(`${apiUrl}/api/v1/clusters`);
  if (!res.ok) {
    throw new Error("Failed to fetch clusters from API");
  }
  try {
    const data = (await res.json()) as ClustersResponse;
    return data;
  } catch (error) {
    console.error(error);
    throw new Error("Failed to parse response from API");
  }
}
export default function Home() {
  const data = useLoaderData<typeof loader>();

  return (
    <>
      <div className="max-w-7xl pt-10 mx-auto bg-gray-100 border border-neutral-100 shadow-sm">
        <div className="mt-8 space-y-5 px-4 pb-8">
          <h1 className="text-4xl font-bold">
            Clusters ({data.cluster_count})
          </h1>
          {data.clusters.map((cluster) => {
            return (
              <div className="rounded-md bg-zinc-50 shadow-sm p-2">
                <h5 className="font-bold">{cluster.id}</h5>
                {cluster.status === "connected" ? (
                  <span className="inline-flex items-center rounded-md bg-emerald-50 px-2 py-1 text-xs font-medium text-emerald-700 ring-1 ring-emerald-600/10 ring-inset">
                    {cluster.status}
                  </span>
                ) : (
                  <span className="inline-flex items-center rounded-md bg-red-50 px-2 py-1 text-xs font-medium text-red-700 ring-1 ring-red-600/10 ring-inset">
                    {cluster.status}
                  </span>
                )}

                <div>
                  <p>Topics: {cluster.topic_count}</p>
                  <p>Size (in bytes): {cluster.total_size}</p>
                </div>

                <p className="text-neutral-500">
                  Brokers: {cluster.brokers.join(",")}
                </p>
              </div>
            );
          })}
        </div>
      </div>
    </>
  );
}
