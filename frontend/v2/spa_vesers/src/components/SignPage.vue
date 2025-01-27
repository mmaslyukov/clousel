<template>
    <div id="sign" style="display: flex; flex-direction: column; vertical-align: middle; margin: auto; padding: auto;">
        <div id="sign-logo-container">
            <img id="sign-logo" src="../assets/vesers.png" />
            <!-- <p id="sign-logo">LOGO</p> -->
        </div>
        <div id="signin-container">
            <div id="signin-container-username" class="user-input-container">
                <p class="lable">Username</p>
                <input class="field" type="text" id="username-value" name="username" placeholder="username"
                    v-model="signData.username" @focus="onFocus">
            </div>
            <div id="signin-container-email" v-show="signData.modeSignup" class="user-input-container">
                <p class="lable">Email</p>
                <input class="field" type="email" id="email-value" name="email" placeholder="user@mail.com"
                    pattern=".+@.+\..+" v-model="signData.email" @focus="onFocus">
            </div>
            <div id="signin-container-password" class="user-input-container">
                <p class="lable">Password</p>
                <input class="field" type="password" id="password-value" name="password" placeholder="password"
                    v-model="signData.password" @focus="onFocus">
            </div>
            <div id="signin-container-password-repeat" v-show="signData.modeSignup" class="user-input-container">
                <p class="lable">Repeat Password</p>
                <input class="field" type="password" id="password-value-repeat" name="password" placeholder="password"
                    v-model="signData.passwordRepeat" @focus="onFocus">
            </div>
            <div id="signin-button-container">
                <button id="login-button" class="sign-button" v-show="!signData.modeSignup"
                    @click="onLoginClick">Login</button>
                <button id="register-lable-clickable" v-show="!signData.modeSignup"
                    @click="regModeHander">Registeration</button>
                <button id="register-button" class="sign-button" v-show="signData.modeSignup"
                    @click="onRegisterClick">Register</button>
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
import axios from 'axios'
import cookiejar from 'cookiejar'
import { onMounted, reactive, defineEmits } from 'vue'
import { config } from '../utils/config.js'
const signData = reactive({
    modeSignup: false,
    username: "",
    password: "",
    passwordRepeat: "",
    email: "",
})
const emit = defineEmits(["LoggedIn", "Registered"])

onMounted(() => {
})
function regModeHander() {
    console.log("click")
    signData.modeSignup = true
}

function onLoginClick() {
    requestLogin()
}

function onRegisterClick() {
    let pswd = document.getElementById("password-value")
    let pswdr = document.getElementById("password-value-repeat")
    let name = document.getElementById("username-value")
    let email = document.getElementById("email-value")

    if (signData.username.length == 0) {
        name.classList.add("field-error")
    } else if (signData.email.length == 0) {
        email.classList.add("field-error")
    } else  if ((signData.password.length == 0) || (signData.password != signData.passwordRepeat)) {
        pswd.classList.add("field-error")
        pswdr.classList.add("field-error")
    } else {
        pswd.classList.remove("field-error")
        pswdr.classList.remove("field-error")
        requestRegister()
    }
}

async function requestLogin() {
    let url = config.url().login()
    try {
        var details = {
            'Username': signData.username,
            'Password': signData.password,
        };
        const formBody = Object.entries(details).map(([key, value]) => encodeURIComponent(key) + '=' + encodeURIComponent(value)).join('&')
        const response = await axios.post(url, formBody);
        if (response.status == 200) {
            // console.log(response.data)
            let t = new cookiejar.Cookie()
            t.name = "Tocken"
            t.value = response.data.Tocken
            document.cookie = t.toString()
            emit('LoggedIn')
        }
    } catch (e) {
        if (e.status == 401 || e.status == 406) {
            let usernameField = document.getElementById("username-value")
            usernameField.classList.add("field-error")
            let passwordField = document.getElementById("password-value")
            passwordField.classList.add("field-error")
        } else {
            console.error(e)
        }
    }
}

async function requestRegister() {
    let url = config.url().register()
    try {
        var details = {
            'Username': signData.username,
            'Company': config.company,
            'Email': signData.email,
            'Password': signData.password,
        };
        const formBody = Object.entries(details).map(([key, value]) => encodeURIComponent(key) + '=' + encodeURIComponent(value)).join('&')


        const response = await axios.post(url, formBody);
        if (response.status == 200) {
            // console.log(response.data)
            let t = new cookiejar.Cookie()
            t.name = "Tocken"
            t.value = response.data.Tocken
            document.cookie = t.toString()
            emit('Registered')
        }

    } catch (e) {
        if (e.status == 409) {
            let usernameField = document.getElementById("username-value")
            usernameField.classList.add("field-error")
            let passwordField = document.getElementById("email-value")
            passwordField.classList.add("field-error")
        } else {
            console.error(e)
        }
    }
}
function onFocus(e) {
    console.log("onFocus")
    if (e.target.classList.contains("field-error")) {
        e.target.value = ""
    }
    e.target.classList.remove("field-error")
}

// async function requestRegister() {
// }
</script>

<style>
#sign {
    position: absolute;
    height: 100%;
    width: 100%;
    background-color: beige;
    overflow-y: auto;
}

#sign-logo-container {
    position: relative;
    display: flex;
    flex-direction: column;
    vertical-align: middle;
    margin: auto;
    margin-top: 70px;
    margin-bottom: 30px;
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
    align-items: center;
    justify-items: center;
    height: 100%;
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
#signin-container-username,
#signin_container_password_repeat,
#signin_container_password,
#signin_container_email
*/
.user-input-container {
    float: right;
    position: relative;
    display: flex;
    flex-direction: column;
    margin: auto;
    margin-top: 10px;
    margin-bottom: 0px;
    /*
    align-items: center;
    justify-content: center;
    right:25%;
    */
}

.field {
    height: 30px;
    padding: 0;
    border-radius: 20px;
    padding-top: 5px;
    padding-bottom: 3px;
    padding-left: 15px;

    background-color: cornsilk;
    font-family: sans-serif;
    color: #404040ff;
    font-size: 12pt;
}

.field:invalid {
    color: brown;
}


.field-error {
    border-color: red;
    color: red;
}

.lable {
    margin: 0;
    margin-left: 10px;
}

#signin-button-container {
    margin: auto;
    display: flex;
    margin-top: 40px;
    margin-bottom: 200px;
    display: flex;
    flex-direction: column;
    align-items: center;
}

#register-lable-clickable {
    margin-top: 20px;
    text-align: center;
    text-decoration: underline;
    border: none;
    background: transparent;
    color: rgb(84, 84, 243);
    font-size: 12pt;
}

.sign-button {
    border: none;
    height: 35px;
    width: fit-content;
    padding-left: 20px;
    padding-right: 20px;
    font-size: 12pt;
    border-radius: 10px;
    color: #404040ff;
    border-style: solid;
    border-width: 1px;
    background-color: #ffeb99ff;
}


.sign-button:active {
    background-color: #1a1a1aff;
    color: #cca400ff;
}
</style>