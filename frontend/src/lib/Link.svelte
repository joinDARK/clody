<script lang="ts">
    interface Props {
        isButton?: boolean;
        link?: string;
        lable?: string;
    }

    let { isButton = false, link = "#", lable = "" }: Props = $props();
</script>

<a href={link} class={`${isButton ? "btn btn_primary" : ""}`}>
    {lable}
</a>

<style>
    .btn {
        border-radius: 1rem;
        padding-inline: 1.25rem;
        padding-block: 0.55rem;
        cursor: pointer;
        transition: background-color 150ms ease-in, transform 150ms ease-in;
    }

    .btn_primary {
        background-color: var(--color-secondary);
        color: #fff;
    }
    .btn_primary:hover {
        background-color: var(--color-secondary-hover);
    }
    .btn_primary:active {
        transform: scale(0.97);
    }
</style>
