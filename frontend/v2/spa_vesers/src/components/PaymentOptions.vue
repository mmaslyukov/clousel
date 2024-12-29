<template>
    <div id="po-area">
        <div id="po-options">
            <div style="margin-left: 10px;"></div>
            <div v-for="(item, i)  in props.details" :key=i style="display:flex">
                <input type="radio" :id="'radioApple'+i" name="radioFruit" :value="item.prodId">
                <label :for="'radioApple'+i" class="po-prod-p" @click="updateInfoData(item.tickets, item.price, item.prodId)">
                    <span style="display: block; ">
                        <div style="display: flex; flex-direction: row; height: 30px;">
                            <p id="po-prod-tn" style="margin:0; margin-top: 10px; margin-left: 6px; width:20px; text-align: right; "
                                class="po-font">{{item.tickets}}</p>
                            <p style="margin-top: 10px; margin-left: 5px;" class="po-font">Tickets</p>
                        </div>
                        <div
                            style="display: flex; flex-direction: row; margin:0; height: 30px; justify-content: center; align-items: center;">
                            <p class="po-font" style="font-size: 20pt; font-weight: bold;">{{item.price}}</p>
                            <p class="po-font" style="font-size: 20pt; font-weight: bold;">€</p>
                        </div>
                    </span>
                </label>
            </div>
            
            <div style="margin-right: 10px;"></div>

            
        </div>
        <div style="display: flex; flex-direction: row;">
            <div id="po-info">
                <p id="po-info-ppt-lable" class="po-font">Price per 1 ticket:</p>
                <div style="display: flex; flex-direction: row; ">
                    <p id="po-info-ppt-value" class="po-font ppt-value">{{infoPricePerTicket}}</p>
                    <p class="po-font ppt-value">€</p>
                </div>
                <div style="display: flex; flex-direction: row; ">
                    <p id="po-info-ppt-save" class="po-font ppt-save">{{infoSaved}}</p>
                    <p class="po-font ppt-save">€ Saved</p>
                </div>
            </div>
            <div style="display: flex; flex-direction: row;justify-content: center; align-items: center">
                <button id="po-buy" class="po-font" type="button" @click="emit('buyClicked', selectedPriceOption)">BUY</button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { onMounted, defineProps, defineEmits, ref } from 'vue'

const props = defineProps(["details"])
const emit = defineEmits(["buyClicked"])
const infoPricePerTicket = ref(0)
const infoSaved = ref(0)
var selectedPriceOption={
    price:Number,
    tickets: Number,
    priodId:String,
}
onMounted(() => {

})
function calculatePricePerTicket(tickets, price) {
    if (tickets == 0) {
        return 0
    }
    return (price / tickets).toFixed(2)
}
function updateInfoData(tickets, price, prodId) {
    let maxPricePerTicket = 0    
    for (let o of props.details) {
        let cp = calculatePricePerTicket(o.tickets, o.price)
        maxPricePerTicket = cp > maxPricePerTicket ? cp : maxPricePerTicket
    }
    infoPricePerTicket.value = calculatePricePerTicket(tickets, price)
    infoSaved.value = (maxPricePerTicket * tickets - infoPricePerTicket.value * tickets).toFixed(2)
    selectedPriceOption.price = price
    selectedPriceOption.tickets = tickets
    selectedPriceOption.prodId = prodId
}

// function onBuyClicked() {
//     console.log("onBuyClicked")
// }
</script>

<style>
.po-font {
    font-family: sans-serif;
    color: #404040ff;
}

#po-area {
    /*
    border-top: solid;
    */
    border-color: #404040ff;
    box-shadow: 0px 5px 10px #404040ff;
    position: absolute;
    top: 70px;
    width: 100%;
    height: 220px;
    background-color: #ffcc00ff;

}

#po-options {
    border: none;
    padding: 0;
    margin: auto;
    background-color: #404040ff;
    height: 100px;
    display: flex;
    flex-direction: row;
    overflow-y: hidden;
    overflow-x: scroll;
}

#po-options input[type="radio"] {
    opacity: 0;
    position: fixed;
    width: 0;
}

/*
    */
#po-options input[type="radio"]:checked+label {
    /*
    background-color: #ffcc00ff;
    */
    background-color: #ffeb99ff;
    border-style: solid;
    border-color: #404040ff;
}

.po-prod-p {
    min-width: 100px;
    height: 70px;
    margin: 15px;
    margin-left: 10px;
    margin-right: 10px;
    background-color: white;
    box-shadow: 0px 0px 10px black;

}

/*
background-color: #ffeb99ff;
.po-prod {
    width: 100px;
}

.po-prod-li {
    list-style-type: none;
    margin: 15px;
    margin-left: 20px;
    margin-right: 0px;
    width: 100px;
    background-color: #ffeb99ff;
    box-shadow: 0px 0px 10px black;
}

.po-prod-li-last {
    list-style-type: none;
    margin-right: 20px;
}
*/
#po-info {
    /*
    */
    height: 75px;
    width: 130px;
    margin: auto;
    margin-top: 15px;
    margin-left: 20px;
    margin-bottom: 10px;
    border-style: solid;
    border-color: #404040ff;
    border-radius: 15px;
}

#po-info-ppt-lable {
    margin: 0px;
    margin-top: 8px;
    margin-left: 15px;
    font-size: 10pt;
}

.ppt-value {
    margin: 0;
    margin-top: 0px;
    font-size: 16pt;
    font-weight: bold;
}

#po-info-ppt-value {
    margin-left: 15px;
}

.ppt-save {
    margin-top: 1px;
    font-size: 10pt;
    font-weight: bold;
    color: #c14b1dff;
}

#po-buy:active {
    background-color: #1a1a1aff;
    color: #cca400ff;
}

#po-buy {
    border: none;
    width: 120px;
    height: 50px;
    margin-right: 25px;
    font-weight: bold;
    font-size: 20pt;
    background-color: #404040ff;
    color: #ffeb99ff;
    border-radius: 10px;
}

#po-info-ppt-save {
    margin-left: 15px;
}

/*
font-weight: bold;
*/
</style>