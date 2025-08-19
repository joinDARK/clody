<script lang="ts">
    import Button from "$lib/Button.svelte";
    import Header from "$lib/Header.svelte";
    import Sidebar from "$lib/Sidebar.svelte";
    import Link from "$lib/Link.svelte";
    import type { ITableApi } from "../../share/interfaces/api";

    console.log("Hello, SvelteKit!");

    let { data }: { data: ITableApi } = $props();
</script>

<div class="screen__grid p-3">
    <div>
        <Header title={data.data.name}>
            <Button small>Табличный вид</Button>
            <Button small>Интерфейс</Button>
        </Header>
    </div>
    <div class="flex gap-3">
        <div class="py-3 pr-3 w-full">
            
        </div>
    </div>
</div>

<style>
    .screen__grid {
        display: grid;
        grid-template-rows: 50px 1fr;
        height: inherit;
        gap: 0.75rem;
    }
</style>
