<script lang="ts">
    import type { Snippet } from "svelte";

    interface Props {
        title?: string;
        children: Snippet;
        spaceBetween?: boolean;
    }

    let { children, title = "Шапка", spaceBetween = false }: Props = $props();
</script>

<div class="h-full flex">
    <div class="header px-4">
        <div class="text-white font-bold w-fit">{title}</div>
        <div class={`${spaceBetween ? "flex justify-between flex-1" : ""}`}>
            {@render children()}
        </div>
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
        height: 50px;
    }
</style>
