<script setup>
import { onMounted, watch, defineProps } from 'vue'
import { Widget } from './widget.js'
// A "ref" is a reactive data source that stores a.
// Technically, we don't need to wrap the string with ref()
// in order to display it, but we will see in the next
// example why it is needed if we ever intend to change
// the.
// const lable = ref('lable')


const props = defineProps({
    current: Number, // ref() of Number
    max: Number, // ref() of Number
    status: String,
    show: Boolean,
})

onMounted(() => {
    var c = document.getElementById("myCanvas");
    c.width = window.innerWidth;
    if (window.innerWidth > window.innerHeight) {
        c.width = 500;
    }
    c.height = c.width;

    console.log("Make widget")
    /* Get last element from array derived from URL*/
    // let id = window.location.href.split('/').at(-1);
    // console.log("url ", window.location.href);
    // console.log("id ", id);
})


watch(props, (newProps) => {
    console.log("Watch widget, tickets: ", props.current, " max: ", props.max)
    // redraw(newProps.current, newProps.max);
    let c = document.getElementById("myCanvas");
    let widget = new Widget(c, c.width / 20, c.width / 3);
    widget.clear().withStatus(props.status).draw(props.current, props.max).lable("Tickets");
})
</script>

<script>
</script>

<template>
    <div class="canvas_container" v-show="props.show">
        <canvas id="myCanvas" class="canvas"> </canvas>
    </div>

    <!-- <h1>{{ message }}</h1> -->
</template>

<style>
.canvas {
    width: 100%;
    height: 100%;
}

.canvas_container {
    /* padding-top: 40px; */
    max-width: fit-content;
    margin: auto;
    /* border-style: solid; */
}
</style>