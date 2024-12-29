<template>
    <div id="header">
        <div id="header_toolbar">
            <div id="header_icon">
                <canvas id="header_icon_canvas" width="70px" height="70px">
                </canvas>
            </div>
            <div id="header_balance">
                <p id="header_balance_lable">BALANCE</p>
                <div id="header_balance_value">
                    <p id="header_balance_tickets_num">{{ tickets.balance }}</p>
                    <p id="header_balance_tickets_lable">Tickets</p>
                </div>
            </div>

            <div id="header_payment_toggle">
                <canvas id="header_payment_toggle_canvas" width="60" height="60">
                </canvas>
            </div>
        </div>
        <div id="header_payment">
        </div>
    </div>
</template>

<script setup>
import { onMounted, defineProps, defineEmits } from 'vue'
import { Logo } from './Logo.js'
import { CanvasButtonRefill } from './CanvasButton.js'
import { HSLColor } from './HSLColor.js'

const emit = defineEmits(["poChanged"])
const tickets = defineProps(['balance'])
var buttonArray = []

onMounted(() => {
    {
        let c = document.getElementById("header_payment_toggle_canvas");
        let cv = new CanvasButtonRefill(c, c.width, c.height, new HSLColor(48, 100, 53), new HSLColor(0, 0, 25)).redraw()
        buttonArray.push(cv)
    }
    {
        let c = document.getElementById("header_icon_canvas");
        // c.width=c.clientWidth
        // c.height=c.clientHeight
        new Logo(c, c.width, c.height, new HSLColor(48, 100, 53), new HSLColor(0, 0, 25)).redraw()
    }
    registerClickers()
})

function registerClickers() {
    for (let b of buttonArray) {
        b.canvas().addEventListener('click', () => {
        b.onClickApplyVisual().redraw()
        emit("poChanged", b.isEnbaled())
        });
    }
}
</script>

<style>
canvas {
    width: 100%;
    height: 100%;
}

#header {
    position: absolute;
    width: 100%;
    margin-top: 0px;
    height: 70px;
    background-color: #ffcc00ff;
    box-shadow: 0px 5px 10px #404040ff;

}

#header_toolbar {
    /**
align-items: left;
border-color: red;
*/
    height: 100%;
    display: flex;
    margin-left: 0;

}



#header_icon {
    margin: auto;
    height: 50px;
    width: 50px;
    margin-left: 10px;

}

#header_icon_canvas {
    background-color: transparent;
}

#header_payment_toggle {
    width: 60px;
    height: 60px;
    margin: auto;
    margin-right: 20px;
}

#header_payment_toggle_canvas {
    margin: auto;
    padding: auto;
}

#header_icon {
    margin-left: 5px;
    margin-right: 5px;
}
#header_balance {
    padding: 3px;
    margin: auto;
    margin-left: 0;
    /*
    margin-top: 10px;
    border-style: solid;
    display: flex;
    flex-direction: column;
    */

}

#header_balance_lable {
    margin: auto;
    padding: 0;
    font-family:  sans-serif;
    font-weight: bold;
    font-size: 22px;
    color: #404040ff;
}

#header_balance_value {
    margin: auto;
    padding: 0;
    display: flex;
    font-family:  sans-serif;
    font-weight: bold;
    font-size: 18px;
}

#header_balance_tickets_num {
    margin: 0;
    padding: 5px;
    width: 20px;
    color: white;
    text-align: center;
    background-color: #404040ff;
}

#header_balance_tickets_lable {
    margin: 0;
    padding: 5px;
    color: #404040ff;
}
</style>