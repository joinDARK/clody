<script lang="ts">
    import type { Snippet } from "svelte";

    interface Props {
        color?: "primary" | "neutral";
        small?: boolean;
        children: Snippet;
    }

    let { children, color = "primary", small = false }: Props = $props();

    function styleColor(color: "primary" | "neutral"): string {
        return color === "primary" ? "btn_primary" : "btn_neutral";
    }
</script>

<button class={`btn ${styleColor(color)} ${small ? 'text-sm' : ''}`}>
    {@render children()}
</button>

<style>
    .btn {
        border-radius: 1rem;
        padding-inline: 1.25rem;
        padding-block: 0.55rem;
        cursor: pointer;
        transition: background-color 150ms ease-in, transform 150ms ease-in;
    }

    .btn_primary {
        background-color: oklch(62.3% 0.214 259.815);
        color: #fff;
    }
    .btn_primary:hover {
        background-color: oklch(54.6% 0.245 262.881);
    }
    .btn_primary:active {
        transform: scale(0.97);
    }

    .btn_neutral {
        background-color: oklch(37.1% 0 0);
        color: oklch(97% 0 0);
    }
    .btn_neutral:hover {
        background-color: oklch(26.9% 0 0);
    }
    .btn_neutral:active {
        transform: scale(0.97);
    }
</style>
