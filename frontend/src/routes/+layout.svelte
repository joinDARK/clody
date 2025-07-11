<script>
    let { children } = $props();
    import "../app.css";
</script>

<main class="flex flex-col gap-2 items-center justify-center h-svh bg-neutral-900">
    {@render children()}
</main>