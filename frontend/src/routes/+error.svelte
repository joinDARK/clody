<script lang="ts">
    import { page } from "$app/state";
</script>

<div class="w-full h-full flex items-center justify-center">
    <h1 class="text-white error-title">{page.status}: {page.error?.message}</h1>
</div>

<style>
    .error-title {
        font-size: 1.5rem;
    }
</style>
