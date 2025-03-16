import type { DrinkPage } from "$lib/types";
import type { PageLoad } from "./$types";

export const load: PageLoad = async function ({ fetch }) {
  const response = await fetch("/data/pages.json");

  if (!response.ok) {
    throw new Error(`Failed to load posts: ${response.status}`);
  }

  const data: DrinkPage[] = await response.json();

  const categories: Record<string, DrinkPage[]> = {};

  for (const drink of data) {
    const category = drink.metadata.category;
    if (!categories[category]) {
      categories[category] = [];
    }
    categories[category].push(drink);
  }

  return { categories };
};
