<script lang="ts">
    import type { PageProps } from "./$types";
    
    console.log("Hello, SvelteKit!");
    
    let { data }: PageProps = $props();
</script>

<h1 class="text-6xl font-bold text-blue-500 font-sans">Welcome to Dater</h1>
<div class="text-white">{@html data?.message}</div>
