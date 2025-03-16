<script lang="ts">
  import { base } from "$app/paths";
  import { capitalize } from "$lib/utils";

  let data = $props();
  let categories = data.data.categories;
</script>

<svelte:head>
  <title>Patrick's Bevs</title>
</svelte:head>

<div>
  <h2>Bevs</h2>
  {#each Object.keys(categories) as category}
    <div>
      <h3 class="text-lg font-bold">{capitalize(category)}</h3>
      <ul class="space-y-2">
        {#each categories[category] as d}
          <li>
            <a href={`${base}/${d.metadata.category}/${d.metadata.slug}`}>
              {d.metadata.title}
            </a>
          </li>
        {/each}
      </ul>
    </div>
  {/each}
</div>
