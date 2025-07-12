<script lang="ts">
    import type { PageProps } from "./$types";
    import Button from "$lib/Button.svelte";
    import Header from "$lib/Header.svelte";
    import Sidebar from "$lib/Sidebar.svelte";

    console.log("Hello, SvelteKit!");

    let { data }: PageProps = $props();
</script>

<div class="screen__grid p-3">
    <div>
        <Header>
            <Button small>Табличный вид</Button>
            <Button small>Интерфейс</Button>
        </Header>
    </div>
    <div class="flex gap-3">
        <Sidebar title="Тестовая БД">
            <Button small>Таблица 1</Button>
            <Button small color="neutral">Таблица 2</Button>
        </Sidebar>
        <div class="py-3 pr-3 w-full">
            <h1 class="text-6xl font-bold text-blue-500 font-sans">
                Welcome to Dater
            </h1>
            <div class="text-white">{@html data?.message}</div>
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
