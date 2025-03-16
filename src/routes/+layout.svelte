<script lang="ts">
	import Header from '$lib/components/Header.svelte';
	import '../app.css';
	let { children } = $props();
</script>

<div class="font-display bg-background text-primary flex min-h-screen justify-center p-6">
	<div class="prose prose-a:no-underline mt-6 w-full max-w-prose text-[16px] leading-relaxed">
		<Header />
		<main class="link-color prose prose-sm">
			{@render children()}
		</main>
	</div>
</div>
