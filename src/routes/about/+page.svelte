<script lang="ts">
  import SocialLink from "$lib/components/SocialLink.svelte";
  import { base } from "$app/paths";
</script>

<div>
  <h2>About</h2>
  <section>
    Just a place where I can keep a list of cocktails I've tried and recipes for
    cocktails I've made.
  </section>
  <section>
    <ul
      class="font-sm flex flex-col mt-1 mb-0 space-x-0 space-y-2 md:flex-row md:space-x-5 md:space-y-0 list-none pl-0"
    >
      <SocialLink title="main site" link="https://pdewey.com" />
      <SocialLink title="rss" link="{base}/rss.xml" />
    </ul>
  </section>
</div>
