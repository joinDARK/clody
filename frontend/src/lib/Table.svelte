<script lang="ts">
    import { onMount } from "svelte";

    type TestData = {
        name: string;
        count: number;
    };

    type Column<T> = {
        key: keyof T;
        name: string;
        position: number;
        accessorFn?: (item: T) => T[keyof T];
        width?: number;
    };

    const testData: TestData[] = [
        { name: "John", count: 10 },
        { name: "Jane", count: 20 },
        { name: "Bob", count: 30 },
    ];

    let columns: Column<TestData>[] = [
        {
            key: "name",
            name: "Имя",
            position: 1,
            accessorFn: (item) => item.name,
        },
        {
            key: "count",
            name: "Количество",
            position: 2,
            accessorFn: (item) => item.count,
        },
    ];

    let headerCells: HTMLDivElement[] = [];

    // 📏 метод для замера ширины
    function updateColumnWidths() {
        headerCells.forEach((cell, i) => {
            columns[i].width = cell.offsetWidth;
        });
    }

    onMount(() => {
        updateColumnWidths();
        window.addEventListener("resize", updateColumnWidths);
        return () => window.removeEventListener("resize", updateColumnWidths);
    });

    // 🎯 обработка drag-resize
    function startResize(e: MouseEvent, columnIndex: number) {
        e.preventDefault();

        const startX = e.clientX;
        const startWidth =
            columns[columnIndex].width || headerCells[columnIndex].offsetWidth;

        function onMouseMove(ev: MouseEvent) {
            const delta = ev.clientX - startX;
            columns[columnIndex].width = Math.max(50, startWidth + delta); // минимум 50px
        }

        function onMouseUp() {
            window.removeEventListener("mousemove", onMouseMove);
            window.removeEventListener("mouseup", onMouseUp);
        }

        window.addEventListener("mousemove", onMouseMove);
        window.addEventListener("mouseup", onMouseUp);
    }
</script>

<div class="table">
    <div class="flex flex-1 text-white flex-col table__container">
        <!-- HEADER -->
        <div class="table__header">
            {#each columns as column, i}
                <div
                    class="table__cell header"
                    bind:this={headerCells[i]}
                    style="width: {column.width ? column.width + 'px' : 'auto'}"
                >
                    {column.name}
                    <!-- Ресайзер справа -->
                    <div
                        class="resizer"
                        on:mousedown={(e) => startResize(e, i)}
                    ></div>
                </div>
            {/each}
        </div>

        <!-- BODY -->
        <div class="table__body">
            {#each testData as item}
                <div class="table__row">
                    {#each columns as column}
                        <div
                            class="table__cell"
                            style="width: {column.width
                                ? column.width + 'px'
                                : 'auto'}"
                        >
                            {column.accessorFn
                                ? column.accessorFn(item)
                                : item[column.key]}
                        </div>
                    {/each}
                </div>
            {/each}
        </div>
    </div>
</div>

<style>
    .table {
        display: flex;
        border-radius: 1rem;
        
        overflow: hidden;
        border: 1px solid oklch(0.295 0 360);
    }

    .table__header {
        display: flex;
        border-bottom: 1px solid oklch(0.295 0 360);
    }
    
    .table__container {
        width: 100%;
    }

    .table__row {
        display: flex;
    }

    .table__cell {
        position: relative; /* нужно для ресайзера */
        padding: 6px 1rem;
        background-color: oklch(0.241 0 360);
        overflow: hidden;
    }

    .table__row:not(:last-child) .table__cell {
        border-bottom: 1px solid oklch(0.295 0 360);
    }

    .table__cell:not(:last-child) {
        border-right: 1px solid oklch(0.295 0 360);
    }

    .table__cell:last-child {
        flex: 1;
    }

    /* 🎯 ресайзер справа */
    .resizer {
        position: absolute;
        top: 0;
        right: 0;
        width: 1.5px;
        height: 100%;
        cursor: col-resize;
        user-select: none;
    }

    .resizer:hover {
        background: oklch(0.4 0 360 / 0.5);
    }
</style>
