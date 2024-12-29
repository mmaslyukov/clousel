
import {HSLColor} from "./HSLColor.js"
export class CanvasButton {

    constructor(canvas, w, h, colorBg, colorFg) {
        this.c = canvas
        this.cw = w
        this.ch = h
        this.w = w * 0.7
        this.h = h * 0.7
        this.cfge = new HSLColor(0, 100, 28)
        this.cfgo = colorFg.clone()
        this.cbge = colorBg.clone()
        this.cbgo = colorBg.clone()
        this.cbg = this.cbgo
        this.cfg = this.cfgo
        
        this.cbge.s = 40
        this.ctx = canvas.getContext("2d");
        this.enabled = false


    }
    inBounds(mx, my) {
        return !(mx < 0 || mx > this.w || my < 0 || my > this.h);
      }
    
    canvas() {
        return this.c
    }

    clear() {
        this.ctx.reset()
        return this
    }
    isEnbaled() {
        return this.enabled
    }
    onClickApplyVisual() {
        if (this.enabled) {
            this.turnOff()
        } else {
            this.turnOn()
        }
        // this.cfg = this.cfg == this.cfgo ? this.cfge : this.cfgo;
        // this.cbg = this.cbg == this.cbgo ? this.cbge : this.cbgo;
        return this
    }
    turnOff() {
        this.enabled = false
        this.cfg = this.cfgo
        this.cbg = this.cbgo
        return this
    }
    turnOn() {
        this.cfg = this.cbgo
        this.cbg = this.cfgo
        this.enabled = true
        return this
    }
    _fixDPI() {
        this.c.width = this.cw * devicePixelRatio;
        this.c.height =this.ch * devicePixelRatio;
        this.ctx.scale(devicePixelRatio, devicePixelRatio);
    }
    redraw() {
        this.clear()
        this._fixDPI()
        this._redraw()
        return this
    }
    _redraw() {
        console.log("Redraw Generic")
    }
}


export class CanvasButtonRefill extends CanvasButton {
    constructor(canvas, w, h, colorBg, colorFg) {
        super(canvas, w, h, colorBg, colorFg)
    }
    _redraw() {
        this.ctx.fillStyle = this.cbg.str();
        this.ctx.strokeStyle = this.cfg.str()

        this.ctx.lineWidth = 3;
        this.ctx.arc(this.cw / 2, this.ch / 2, this.cw * 0.35, 0, Math.PI * 2, true);
        if (this.enabled) {
            this.ctx.stroke();
            this.ctx.fill();
        } else {
            this.ctx.fill();
            this.ctx.stroke();
        }

        this.ctx.beginPath();
        this.ctx.moveTo(this.cw / 2, (this.ch - this.h));
        this.ctx.lineTo(this.cw / 2, this.h);
        this.ctx.moveTo((this.cw - this.w), this.ch / 2);
        this.ctx.lineTo(this.w, this.ch / 2);
        this.ctx.stroke();
    }
}


export class CanvasButtonScan extends CanvasButton {
    constructor(canvas, w, h, colorBg, colorFg) {
        super(canvas, w, h, colorBg, colorFg)
        this.w = w * 0.9
        this.h = h * 0.9
    }

    _redraw() {
        this.ctx.fillStyle = this.cbg.str();
        this.ctx.strokeStyle = this.cfg.str()

        this.ctx.lineWidth = 3;
        let x = (this.cw - this.w)/2
        let y = (this.ch - this.h)/2 
        
        this.ctx.beginPath();
        this.ctx.moveTo(x,y)
        this.ctx.lineTo(x, y += this.h);
        this.ctx.lineTo(x += this.w, y );
        this.ctx.lineTo(x, y -= this.h );
        this.ctx.lineTo(x -= this.w +1, y);
        this.ctx.fill();
        this.ctx.closePath();

        function bracked(b, shift, dir) {
            b.ctx.beginPath();
            x = b.cw/10 + shift
            y = b.ch/10
            b.ctx.moveTo(x,y)
            b.ctx.lineTo(x +=10*dir, y);
            y = b.ch/10
            b.ctx.lineTo(x, y += b.h - b.ch/10);
            b.ctx.lineTo(x -=11*dir, y);
            b.ctx.stroke();
            b.ctx.closePath();
        
        }
        bracked(this, 10, -1)
        bracked(this, 46, 1)

        // function line(b, shift, len) {
        //     b.ctx.beginPath();
        //     x = shift
        //     y = b.ch*0.25 + y*(1-len)
        //     // y = y + 1 - y*len
        //     b.ctx.moveTo(x,y)
        //     b.ctx.lineTo(x, y + b.ch*0.5*len);
        //     b.ctx.lineWidth = 5;
        //     b.ctx.stroke();
        //     b.ctx.closePath();
        // }
        // line(this, this.cw/2, 1)
        // line(this, this.cw/3, 0.8)
        // line(this, this.cw*2/3, 0.8)
        this.ctx.fillStyle = this.cfg.str();
        this.ctx.textAlign = "center";
        this.ctx.font =  "12pt  Verdana";
        this.ctx.fillText("SCAN", this.cw*0.48, this.ch*0.38);
        this.ctx.font =  "bold 18pt  Verdana";
        this.ctx.fillText("QR", this.cw*0.48, this.ch*0.75);
    }
}