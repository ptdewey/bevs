import type { DrinkPage } from "$lib/types";
import type { PageLoad } from "./$types";

export const load: PageLoad = async function ({ params, fetch }) {
  const { category, drink } = params;

  const response = await fetch("/data/pages.json");
  if (!response.ok) {
    throw new Error(`Failed to load pages: ${response.status}`);
  }

  const pages: DrinkPage[] = await response.json();

  const data = pages.find(
    (p) => p.metadata.category === category && p.metadata.slug === drink,
  );

  if (!data) {
    return {
      status: 404,
      error: new Error("Page not found"),
    };
  }

  return { data };
};
