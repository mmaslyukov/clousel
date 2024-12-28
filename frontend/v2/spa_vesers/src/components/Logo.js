export class Logo {
    constructor(canvas, w, h, colorBg, colorFg) {
        this.c = canvas
        this.cw = w
        this.ch = h
        this.w = w * 0.9
        this.h = h * 0.95
        this.cfg = colorFg
        this.cbg = colorBg
        this.ctx = canvas.getContext("2d");
        // this.ctx.translate(0.5, 0.5);
        // this.ctx.imageSmoothingEnabled = false;
    }
    _fixDPI() {
        this.c.width = this.cw * devicePixelRatio;
        this.c.height =this.ch * devicePixelRatio;
        this.ctx.scale(devicePixelRatio, devicePixelRatio);
    }
    clear() {
        this.ctx.clear()
        return this
    }
    redraw() {
        this._fixDPI()
        this.ctx.fillStyle = this.cbg.str();
        this.ctx.strokeStyle = this.cfg.str()

        this.ctx.lineWidth = 3;
        let x = this.cw - this.w
        let y = this.ch - this.h 
        this.ctx.beginPath();
        this.ctx.moveTo(x, y);
        this.ctx.lineTo(x, y = this.h);
        this.ctx.lineTo(x = this.w * 0.75 - 0.8, y);
        let r = 15
        this.ctx.arc(x, y = y - r, r, -Math.PI * 1.5, 0, true);
        this.ctx.lineTo(x + r, y = this.ch - this.h);
        this.ctx.lineTo(this.cw - this.w-1, y = this.ch - this.h);

        this.ctx.fill();
        this.ctx.stroke();
        this.ctx.closePath();


        this.ctx.beginPath();
        x = this.cw / 2
        y = this.ch / 2
        this.ctx.moveTo(x, y);
        this.ctx.lineTo(x, y = this.h * 0.87);
        this.ctx.lineTo(x = this.w * 0.7, y);
        r = 8
        this.ctx.arc(x, y = y - r, r, -Math.PI * 1.5, 0, true);
        this.ctx.lineTo(x + r, y = this.ch / 2);

        this.ctx.lineWidth = 5;
        this.ctx.stroke();

        this.ctx.font =  "12pt  sans-serif";
        this.ctx.color = this.cbg
        
        this.ctx.fillStyle = this.cfg.str();
        this.ctx.textAlign = "left";
        this.ctx.fillText("vesers", 10, 18);
    
        // this.ctx.beginPath();
        // this.ctx.fillTe;

    }
}