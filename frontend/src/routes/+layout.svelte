<script lang="ts">
    let { children } = $props();
    import "../app.css";
</script>

<main class="screen">
    {@render children()}
</main>

<style>
    .screen {
        width: 100%;
        height: 100svh;
        background-color: var(--color-neutral-900);
    }
</style>
