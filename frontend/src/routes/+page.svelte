<script lang="ts">
    // import Button from "$lib/Button.svelte";
    import Header from "$lib/Header.svelte";
    import Sidebar from "$lib/Sidebar.svelte";
    import Link from "$lib/Link.svelte";
    import type { IBaseApi } from "@share/interfaces/api";

    console.log("Hello, SvelteKit!");

    let { data }: { data: IBaseApi } = $props();
    const { tables, name } = data.data;
</script>

<div class="screen__grid p-3">
    <div>
        <Header>
            ds
        </Header>
    </div>
    <div class="flex gap-3">
        <Sidebar title={name}>
            <!-- <Button small>Таблица 1</Button>
            <Button small color="neutral">Таблица 2</Button> -->
            {#if tables != undefined && data != undefined}
                {#each tables as table}
                    <Link
                        isButton
                        label={table.name}
                        link={`/table/${table.id}`}
                    />
                {/each}
            {/if}
        </Sidebar>
        <div class="py-3 pr-3 w-full">
            <h1 class="text-6xl font-bold text-primary font-sans">
                Welcome to Dater
            </h1>
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
