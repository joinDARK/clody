<script lang="ts">
    interface Props {
        columns: Column[];
    }

    interface Column {
        id: string;
        title: string;
        type?: "text" | "number" | "date" | "boolean";
        data?: any[]; // TODO: Исправить тип данных
    }

    let { columns }: Props = $props();
</script>

<div class="table">
    {#each columns as column}
        <div class="flex flex-1 text-white flex-col table__column">
            <div class="table__cell">{column.title}</div>
            {#each column.data as data}
                <div class="table__cell">{data}</div>
            {/each}
        </div>
    {/each}
    <div
        class="flex mx-3 text-white flex-col table__column items-center justify-center"
    >
        Добавить колонку
    </div>
</div>

<style>
    .table {
        display: flex;
        border-radius: 1rem;
        overflow: hidden;
        border: 1px solid oklch(0.295 0 360);
    }

    .table__column:not(:last-child) {
        border-right: 1px solid oklch(0.295 0 360);
    }

    .table__cell {
        padding: 6px 1rem;
        background-color: oklch(0.241 0 360);
    }

    .table__cell:not(:last-child) {
        border-bottom: 1px solid oklch(0.295 0 360);
    }
</style>
