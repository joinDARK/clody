<script lang="ts">
    import Button from "$lib/Button.svelte";
    import Header from "$lib/Header.svelte";
    import Sidebar from "$lib/Sidebar.svelte";
    import Link from "$lib/Link.svelte";
    import Table from "$lib/Table.svelte";
    import type { ITableApi } from "../../share/interfaces/api";

    console.log("Hello, SvelteKit!");

    let { data }: { data: ITableApi } = $props();
</script>

<div class="screen__grid p-3">
    <div>
        <Header title={data.data.name} spaceBetween>
            <div class="flex gap-3">
                <Button small>Табличный вид</Button>
                <Button small>Интерфейс</Button>
            </div>
            <Link label="Назад" link="/" isButton small />
        </Header>
    </div>
    <div class="flex gap-3">
        <div class="py-3 w-full">
            <Table/>
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
