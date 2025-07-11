<script lang="ts">
    import { page } from "$app/state";
</script>

<h1 class="text-white error-title">{page.status}: {page.error?.message}</h1>

<style>
    .error-title {
        font-size: 1.5rem;
    }
</style>
