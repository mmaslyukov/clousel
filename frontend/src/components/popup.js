export class PopupIconSuccess {
    constructor(canvas) {
        this._c = canvas
        this._ctx = this._c.getContext("2d")
        this._radius = 20
        this._fontHeight = 20
        this._x = 0
        this._y = 0
    }
    
    pos(x, y) {
        this._x = x
        this._y = y
        return this
    }
    
    clear() {
        this._ctx.reset()
        return this
    }
    
    redraw() {
        return this.clear().draw()
    }
    
    colorFill(c) {
        this._colorFill = c
        return this
    }
    
    colorStrike(c) {
        this._colorStrike = c
        return this
    }
    
    draw() {
        this._ctx.beginPath()
        this._ctx.fillStyle = this._colorStrike
        this._ctx.arc(this._x, this._y, this._radius, 0, Math.PI * 2, true)
        this._ctx.fill()

        this._ctx.beginPath();
        this._ctx.strokeStyle = this._colorFill
        this._ctx.moveTo(this._x-10, this._y);
        this._ctx.lineTo(this._x, this._y+10);
        this._ctx.lineTo(this._x+10, this._y-10);
        this._ctx.fill();
        this._ctx.stroke();
        return this
    }
}


export class PopupIconError {
    constructor(canvas) {
        this._c = canvas
        this._ctx = this._c.getContext("2d")
        this._radius = 20
        this._fontHeight = 20
        this._x = 0
        this._y = 0
    }
    
    pos(x, y) {
        this._x = x
        this._y = y
        return this
    }
    
    clear() {
        this._ctx.reset()
        return this
    }
    
    redraw() {
        return this.clear().draw()
    }
    
    colorFill(c) {
        this._colorFill = c
        return this
    }
    
    colorStrike(c) {
        this._colorStrike = c
        return this
    }
    
    draw() {
        this._ctx.beginPath()
        this._ctx.fillStyle = this._colorStrike
        this._ctx.arc(this._x, this._y, this._radius, 0, Math.PI * 2, true)
        this._ctx.fill()
        

        this._ctx.beginPath();
        this._ctx.strokeStyle = this._colorFill
        this._ctx.moveTo(this._x-10, this._y-10);
        this._ctx.lineTo(this._x+10, this._y+10);

        this._ctx.moveTo(this._x-10, this._y+10);
        this._ctx.lineTo(this._x+10, this._y-10);
        this._ctx.fill();
        this._ctx.stroke();
        return this
    }
}

export class PopupHeader {
    constructor(canvas) {
        this._c = canvas
        this._ctx = this._c.getContext("2d")
        this._x = 0
        this._y = 0
        this._fontSize = 18
        this._msg = "header"
    }
    
    pos(x, y) {
        this._x = x
        this._y = y
        return this
    }
    
    msg(msg) {
        this._msg = msg
        return this
    }
    
    colorFill(c) {
        this._colorFill = c
        return this
    }
    
    clear() {
        this._ctx.reset()
        return this
    }
    
    redraw() {
        return this.clear().draw()
    }
    
    draw() {
        this._ctx.beginPath()
        this._ctx.font = "bold " + this._fontSize + "pt Arial"
        this._ctx.fillStyle = this._colorFill
        this._ctx.textAlign = "left"
        this._ctx.fillText(this._msg, this._x, this._y)
        return this
    }
}

export class PopupBody {
    constructor(canvas) {
        this._c = canvas
        this._ctx = this._c.getContext("2d")
        this._x = 0
        this._y = 0
        this._fontSize = 14
        this._msg = "lable"
        this._colorFill = 'black'
    }
    
    pos(x, y) {
        this._x = x
        this._y = y
        return this
    }
    
    colorFill(c) {
        this._colorFill = c
        return this
    }
    
    msg(msg) {
        this._msg = msg
        return this
    }
    
    clear() {
        this._ctx.reset()
        return this
    }
    
    redraw() {
        return this.clear().draw()
    }
    
    draw() {
        this._ctx.beginPath()
        this._ctx.font = this._fontSize + "pt Arial"
        this._ctx.fillStyle = this._colorFill
        this._ctx.textAlign = "left"
        this._ctx.fillText(this._msg, this._x, this._y)
        return this
    }
}


export class PopupSuccess {
    constructor(canvas, w, h) {
        this._c = canvas
        this._ctx = this._c.getContext("2d")
        this._w = w
        this._h = h
        this._lw = 4
        this._colorFill = "#e2f9e2ff"
        this._colorStrike = "#2e8e37ff"
        this._icon = new PopupIconSuccess(canvas).pos(50, 60).colorFill(this._colorFill).colorStrike(this._colorStrike)
        this._header = new PopupHeader(canvas).msg("SUCCESS").pos(90, 55).colorFill(this._colorStrike)
        this._body = new PopupBody(canvas).pos(90, 80).colorFill(this._colorStrike)
    }
    
    msg(msg) {
        this._body.msg(msg)
        return this
    }

    colorFill(c) {
        this._colorFill = c
        return this
    }
    
    colorStrike(c) {
        this._colorStrike = c
        return this
    }

    clear() {
        this._ctx.reset()
        return this
    }
    
    redraw() {
        return this.clear().draw()
    }

    draw() {
        this._ctx.beginPath();
        this._ctx.fillStyle = this._colorFill
        this._ctx.strokeStyle = this._colorStrike
        this._ctx.lineWidth = this._lw
        this._ctx.roundRect(10, 10, this._c.width - 20, 100, 10)
        this._ctx.fill()
        this._ctx.stroke();
        this._icon.draw()
        this._header.draw()
        this._body.draw()
        return this
    }
}


export class PopupError {
    constructor(canvas, w, h) {
        this._c = canvas
        this._ctx = this._c.getContext("2d")
        this._w = w
        this._h = h
        this._lw = 4
        this._colorFill = "#f9e2e2ff"
        this._colorStrike = "#8e2e2eff"
        this._icon = new PopupIconError(canvas).pos(50, 60).colorFill(this._colorFill).colorStrike(this._colorStrike)
        this._header = new PopupHeader(canvas).msg("FAILURE").pos(90, 55).colorFill(this._colorStrike)
        this._body = new PopupBody(canvas).pos(90, 80).colorFill(this._colorStrike)
    }
    
    msg(msg) {
        this._body.msg(msg)
        return this
    }

    colorFill(c) {
        this._colorFill = c
        return this
    }
    
    colorStrike(c) {
        this._colorStrike = c
        return this
    }

    clear() {
        this._ctx.reset()
        return this
    }
    
    redraw() {
        return this.clear().draw()
    }

    draw() {
        this._ctx.beginPath();
        this._ctx.fillStyle = this._colorFill
        this._ctx.strokeStyle = this._colorStrike
        this._ctx.lineWidth = this._lw
        this._ctx.roundRect(10, 10, this._c.width - 20, 100, 10)
        this._ctx.fill()
        this._ctx.stroke();
        this._icon.draw()
        this._header.draw()
        this._body.draw()
        return this
    }
}
