<template>
    <div id="sign" style="display: flex; flex-direction: column; vertical-align: middle; margin: auto; padding: auto;">
        <div id="sign-logo-container">
            <img id="sign-logo" src="../assets/vesers.png"/>
            <!-- <p id="sign-logo">LOGO</p> -->
        </div>
        <div id="signin-container">
            <div id="signin-container-username">
                <p class="lable">Username</p>
                <input class="field" type="text" id="username-value" name="username">
            </div>
            <div id="signin-container-email" v-show="signData.modeSignup">
                <p class="lable">Email</p>
                <input class="field" type="email" id="email-value" name="email">
            </div>
            <div id="signin-container-password">
                <p class="lable">Password</p>
                <input class="field" type="password" id="password-value" name="password">
            </div>
            <div id="signin-container-password-repeat" v-show="signData.modeSignup">
                <p class="lable">Repeat Password</p>
                <input class="field" type="password" id="password-value-repeat" name="password">
            </div>
            <div id="signin-button-container">
                <button id="login-button" class="sign-button" v-show="!signData.modeSignup" @click="$emit('loginClicked')">LOGIN</button>
                <button id="register-lable-clickable" v-show="!signData.modeSignup" @click="regModeHander">Registeration</button>
                <button id="register-button" class="sign-button" v-show="signData.modeSignup" @click="$emit('registerClicked')">REGISTER</button>
            </div>
            <!-- <div id="signin-button-container">
            </div> -->

        </div>
        <!-- <div id="signup_container">
        <div id="signup_username_lable_container"></div>
        <div id="signup_username_value_container"></div>
        <div id="signup_email_lable_container"></div>
        <div id="signup_email_value_container"></div>
        <div id="signup_password_lable_container"></div>
        <div id="signup_password_value_container"></div>
        <div id="signup_button_container"></div>
    </div> -->
        <div></div>
        <div></div>
    </div>
</template>

<script setup>
import { onMounted, reactive, defineEmits } from 'vue'
const signData = reactive({
    modeSignup: false,
    username: "",
    password: "",
    email: "",
})
defineEmits(["loginClicked", "registerClicked"])

onMounted(() => {
})
function regModeHander() {
    console.log("click")
    signData.modeSignup = true
}
</script>

<style>
#sign {
    position: absolute;
    height: 100%;
    width: 100%;
    background-color: white;
}

#sign-logo-container {
    position: relative;
    display: flex;
    flex-direction: column;
    vertical-align: middle;
    margin: auto;
    margin-top: 100px;
    margin-bottom: 50px;
}

#sign-logo {
    display: flex;
    margin-left: auto;
    margin-right: auto;
}

#signin-container {
    position: relative;
    display: flex;
    flex-direction: column;
    height: 100%;
    align-items: center;
}
/*
#signin-container-username {
    position: relative;
    display: flex;
    flex-direction: row;
    margin: auto;
    margin-top: 25%;
    margin-bottom: 10px;
    margin-right: width/2;
}
*/
#signin-container-username,
#signin_container_password_repeat,
#signin_container_password,
#signin_container_email {
    float: right;
    position: relative;
    display: flex;
    flex-direction: column;
    margin: auto;
    margin-top: 10px;
    margin-bottom: 10px;
/*
    align-items: center;
    justify-content: center;
    right:25%;
    */
}

.field {
    height: 30px;
    border-radius: 20px;
    padding-left: 15px;
}

.lable {
    margin: 0;
    margin-left: 10px;
}

#signin-button-container {
    margin: auto;
    display: flex;
    margin-top: 10%;
    flex-direction: column;
}
#register-lable-clickable {
    margin-top: 10px;
    text-decoration: underline;
    border: none;
    background: transparent;
    color: rgb(84, 84, 243);
}

.sign-button {
    border: none;
    min-width: 120px;
    height: 50px;
    font-weight: bold;
    font-size: 20pt;
    border-radius: 10px;
    background-color: #404040ff;
    color: #ffeb99ff;
}
</style>