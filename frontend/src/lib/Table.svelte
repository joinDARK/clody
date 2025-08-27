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
    let maxColumnWidths: number[] = [];
    // 📏 метод для замера ширины
    function updateColumnWidths() {
        headerCells.forEach((cell, i) => {
            columns[i].width = cell.offsetWidth;
        });
    }
    function calculateMaxWidths() {
        if (headerCells.length === 0) return;
        const canvas = document.createElement("canvas");
        const ctx = canvas.getContext("2d")!;
        const font = getComputedStyle(headerCells[0]).font;
        ctx.font = font;
        columns.forEach((column, i) => {
            const headerText = column.name;
            const cellValues = testData.map((item) =>
                String(
                    column.accessorFn
                        ? column.accessorFn(item)
                        : item[column.key],
                ),
            );
            const allTexts = [headerText, ...cellValues];
            const maxTextWidth = Math.max(
                ...allTexts.map((text) => ctx.measureText(text).width),
            );
            const paddingLeft = parseFloat(
                getComputedStyle(headerCells[i]).paddingLeft,
            );
            const paddingRight = parseFloat(
                getComputedStyle(headerCells[i]).paddingRight,
            );
            const maxContentWidth = maxTextWidth + paddingLeft + paddingRight;
            maxColumnWidths[i] = 1.5 * maxContentWidth;
        });
    }
    onMount(() => {
        updateColumnWidths();
        calculateMaxWidths();
        window.addEventListener("resize", updateColumnWidths);
        return () => window.removeEventListener("resize", updateColumnWidths);
    });
    // 🎯 обработка drag-resize
    function startResize(e: MouseEvent, columnIndex: number) {
        e.preventDefault();
        const startX = e.clientX;
        const startWidth =
            columns[columnIndex].width || headerCells[columnIndex].offsetWidth;
        const maxWidth = maxColumnWidths[columnIndex] || Infinity; // если не рассчитано, без лимита
        function onMouseMove(ev: MouseEvent) {
            const delta = ev.clientX - startX;
            const newWidth = startWidth + delta;
            columns[columnIndex].width = Math.max(
                50,
                Math.min(newWidth, maxWidth),
            );
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
                    {#if i < columns.length - 1}
                        <button
                            class="resizer"
                            aria-label="Resize column"
                            on:mousedown={(e) => startResize(e, i)}
                        ></button>
                    {/if}
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
        max-width: 100svw;
        overflow: hidden;
        border: 1px solid oklch(0.295 0 360);
    }
    .table__header {
        display: flex;
        border-bottom: 1px solid oklch(0.295 0 360);
    }
    .table__container {
        width: 100%;
        overflow-x: auto;
    }
    .table__row {
        display: flex;
    }
    .table__cell {
        position: relative; /* нужно для ресайзера */
        padding: 6px 1rem;
        background-color: oklch(0.241 0 360);
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
        z-index: 999;
        top: 3px;
        right: -2.5px;
        width: 4px;
        height: 30px;
        border-radius: 5px;
        cursor: col-resize;
        user-select: none;
    }
    .resizer:hover {
        background: var(--color-primary);
    }
</style>
