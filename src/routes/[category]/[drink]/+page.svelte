<script lang="ts">
  import type { DrinkPage } from "$lib/types";

  let data = $props();
  let d: DrinkPage = data.data.data;
</script>

<svelte:head>
  <title>{d.metadata.title}</title>
</svelte:head>

<article>
  <h2>{d.metadata.title}</h2>

  <div class="break-words hyphens-auto">
    {@html d.content}
  </div>

  {#if d.metadata.tags}
    <hr />
    <ul class="list-outside list-none pl-0">
      {#if d.metadata.tags.length > 1}
        <li class="pl-0">
          Tags: {d.metadata.tags.join(", ")}
        </li>
      {:else}
        <li class="pl-0">Tag: {d.metadata.tags}</li>
      {/if}
    </ul>
  {/if}
</article>
