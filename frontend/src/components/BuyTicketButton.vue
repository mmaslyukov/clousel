<script setup>
// import { defineProps, watch, onMounted } from 'vue'
import { defineProps, defineEmits, watch, onMounted, nextTick } from 'vue'
import { BuyButton } from './buybotton.js'

const emit = defineEmits(["buyClicked"])
const props = defineProps(['show', 'menu'])
var buyButtonArray = []

function redrawPriceOptions(newProps) {
    let buttonArray = []
    if (newProps == undefined || newProps.menu == undefined) {
        return
    }
    newProps.menu.forEach(function (p, i) {
        if (p == undefined) {
            console.error("P is undefined")
            return
        }
        let c = document.getElementById("Canvas" + i);
        if (c == undefined) {
            console.error("cannot find canvas: Canvas%s ", i)
            return
        }
        let bb = new BuyButton(c, c.width)
        bb.setPrice(p.price).setTickets(p.tickets).setPriceId(p.priceId).redraw()
        buttonArray.push(bb);
    })

    
    if (buttonArray.length != buyButtonArray.length) {
            buyButtonArray = buttonArray
        
        console.log("bba", buyButtonArray)
        for (let b of buttonArray) {
        // console.log("canvas>>>>", b)
                b.canvas().addEventListener('click', (event) => {
                let x = event.pageX - (b.canvas().clientLeft + b.canvas().offsetLeft);
                let y = event.pageY - (b.canvas().clientTop + b.canvas().offsetTop);
                if (b.inBounds(x, y) && !!b.onClick) {
                    b.onClick().clear().redraw()
                    window.setTimeout(() => { b.onRelease().clear().redraw(); }, 200);
                    emit('buyClicked', b)
                }
            });
        }
    }

    return
}

onMounted(() => {
    // console.log("count:", props.count);
    redrawPriceOptions(props)


})
// watch( ()=> props, (newVal)=> { 
//     redrawPriceOptions(newVal)

// })
watch(props, async (newProps) => {
    await nextTick();
    redrawPriceOptions(newProps)
})

</script>


<template>
    <div class="buy_canvas_container" v-show="props.show">
        <canvas v-for="(item, i)  in props.menu" :key=i :id="'Canvas' + i" class="buy_canvas"> </canvas>
        <div class="bottom_padding"> </div>
        <!-- <canvas v-for="x in props.count" :key=x :id="'Canvas' + x" :idx="x" class="buy_canvas"> </canvas> -->
    </div>
</template>

<style>
.buy_canvas_container {
    /* padding-top: 50px; */
    /* max-width: fit-content; */
    display: block;
    margin-left: auto;
    margin-right: auto;
    /* border-style: solid; */
}

.buy_canvas {
    /* padding-top: 5px; */
    padding-left: 0;
    padding-right: 0;
    margin-left: auto;
    margin-right: auto;
    display: block;
     /* border-style: solid;  */
}
.bottom_padding {
    padding-bottom: 5rem;
}
</style>