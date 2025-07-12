<script lang="ts">
    import type { Snippet } from "svelte";

    interface Props {
        children: Snippet;
        title?: string;
    }

    let { children, title = "" }: Props = $props();
</script>

<div class="max-w-56 h-full w-full">
    <div class="sidebar p-4">
        <div class="text-white text-xl mb-5 font-bold">{title}</div>
        <div class="sidebar__content">
            {@render children()}
        </div>
    </div>
</div>

<style>
    .sidebar {
        border-radius: 1rem;
        border: 1px solid oklch(0.295 0 360);
        background-color: oklch(0.241 0 360);
        height: 100%;
    }

    .sidebar__content {
        display: flex;
        flex-direction: column;
        gap: 0.75rem;
    }
</style>
