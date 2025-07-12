<script lang="ts">
    import type { Snippet } from "svelte";

    interface Props {
        title?: string;
        children: Snippet;
    }

    let { children, title = "Шапка" }: Props = $props();
</script>

<div class="h-full flex">
    <div class="header px-4">
        <div class="text-white font-bold">{title}</div>
        {@render children()}
    </div>
</div>

<style>
    .header {
        border-radius: 1rem;
        border: 1px solid oklch(0.295 0 360);
        background-color: oklch(0.241 0 360);
        display: flex;
        align-items: center;
        gap: 0.75rem;
        width: 100%;
    }
</style>
