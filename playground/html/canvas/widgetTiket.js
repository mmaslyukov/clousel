
class WidgetTiket {
  constructor(canvas, w) {
    this.cw = w
    this.ch = w/2
    canvas.width = this.cw
    canvas.height = this.ch
    this.ctx = canvas.getContext("2d");
    this.delta = 0;
    this.x = this.cw*0.018
    this.y = this.cw*0.018
    this.w = this.cw*0.91
    this.h = this.ch*0.83
    this._bezel()
  }
  _bezel() {
    this.ctx.lineWidth = 10;
    this.ctx.fillStyle = '#ccccccff';
    this.ctx.restore();
    this.ctx.beginPath();
    this.ctx.fillRect(0, 0, this.cw, this.ch)
    this.ctx.fill()
    this.ctx.strokeStyle = 'grey';
    this.ctx.strokeRect(0, 0, this.cw, this.ch)
    this.ctx.stroke()
  }
  clear() {
    this.ctx.reset();
    this._bezel();
    return this;
  }
  click(x, y) {

  }
  setPrice(p) {
    this.price = p;
    return this;
  }
  setTikets(t) {
    this.tikets = t;
    return this
  }
  redraw() {
    this._drawTikets(this.tikets)
    this._drawPriceTag(this.price)
    if (this.dim) {
      // this.ctx.fillStyle = 'grey';
      this.ctx.fillStyle = '#00000060';
      this.ctx.restore();
      this.ctx.beginPath();
      this.ctx.fillRect(0, 0, this.cw, this.ch)
      this.ctx.fill()
    }
    // console.log("Widget Redraw, dim:", this.dim)
    return this
  }
  onPress() {
    // console.log("Widget Pressed!")
    this.delta = 10
    this.dim = true
    return this;
  }
  onClick() {
    // console.log("Widget Clicked!")
    this.onPress();
    return this;
    // this.onRelease()
  }
  onRelease() {
    // console.log("Widget Released!")
    this.delta = 0;
    this.dim = false
    return this;
  }

  inBounds(mx, my) {
    return !(mx < this.x || mx > this.x + this.w || my < this.y || my > this.y + this.h);
  }




  _drawPriceTag(price) {
    let x = 0
    let y = 0

    let rx = x + 20 + this.x + this.delta
    let ry = y + 20 + this.y + this.delta
    this.ctx.lineWidth = 3;


    this.ctx.beginPath();
    this.ctx.translate(rx, ry);
    this.ctx.rotate(Math.PI / 5);

    this.ctx.fillStyle = 'grey';
    let shadowDelta = 7
    this.ctx.moveTo(x + shadowDelta, y + shadowDelta);
    // this.ctx.lineTo(x + 35 + shadowDelta, y + 40 + shadowDelta);
    // this.ctx.lineTo(x + 240 + shadowDelta, y + 40 + shadowDelta);
    // this.ctx.lineTo(x + 240 + shadowDelta, y - 40 + shadowDelta);
    // this.ctx.lineTo(x + 35 + shadowDelta, y - 40 + shadowDelta);


    this.ctx.lineTo(x + this.w * 0.07 + shadowDelta, y + this.h * 0.2 + shadowDelta);
    this.ctx.lineTo(x + this.w * 0.48 + shadowDelta, y + this.h * 0.2 + shadowDelta);
    this.ctx.lineTo(x + this.w * 0.48 + shadowDelta, y - this.h * 0.2 + shadowDelta);
    this.ctx.lineTo(x + this.w * 0.07 + shadowDelta, y - this.h * 0.2 + shadowDelta);

    this.ctx.lineTo(x + shadowDelta, y + shadowDelta);
    this.ctx.fill();

    this.ctx.beginPath();
    this.ctx.strokeStyle = '#c57b00ff';
    this.ctx.fillStyle = '#faf846ff';
    this.ctx.moveTo(x, y);

    // this.ctx.lineTo(x + 35, y + 40);
    // this.ctx.lineTo(x + 240, y + 40);
    // this.ctx.lineTo(x + 240, y - 40);
    // this.ctx.lineTo(x + 35, y - 40);

    this.ctx.lineTo(x + this.w * 0.07, y + this.h * 0.2);
    this.ctx.lineTo(x + this.w * 0.48, y + this.h * 0.2);
    this.ctx.lineTo(x + this.w * 0.48, y - this.h * 0.2);
    this.ctx.lineTo(x + this.w * 0.07, y - this.h * 0.2);
    this.ctx.lineTo(x, y);
    this.ctx.fill();
    this.ctx.stroke();
    
    // hole
    this.ctx.beginPath();
    this.ctx.fillStyle = '#ffffffff';
    this.ctx.arc(rx, y, this.h * 0.04, 0, Math.PI * 2, true);
    this.ctx.fill();
    this.ctx.stroke();



    let euro = "€"
    let size = this.h / 3 //70
    let tx = x + this.w / 10 //40
    let ty = y + this.h * 0.125
    this.ctx.font = size + "px Arial";
    this.ctx.fillStyle = '#c57b00ff';
    this.ctx.textAlign = "left";
    this.ctx.fillText(euro, tx, ty);
    let ew = this.ctx.measureText(euro).width
    this.ctx.font = size + "px Impact";
    this.ctx.fillStyle = '#c57b00ff';
    this.ctx.fillText(price.toFixed(2), tx + ew * 1.2, ty);
    
    this.ctx.rotate(-Math.PI / 5);
    this.ctx.translate(-rx, -ry);

    return this;
  }

  _tiket(x, y, num, shadow) {

    if (shadow) {
      this.ctx.fillStyle = 'grey';
      this.ctx.beginPath();
      this.ctx.fillRect(x + 5, y + 5, this.w, this.h)
      this.ctx.fill()
    }

    this.ctx.fillStyle = "HSL(125,55%," + (55 - 10 * num) + "%)";
    // this.ctx.fillStyle = '#43cc4fff';
    this.ctx.beginPath();
    this.ctx.fillRect(x, y, this.w, this.h)
    this.ctx.fill()
    this.ctx.beginPath();
    this.ctx.fillStyle = "HSL(0,0%," + (100 - 10 * num) + "%)";
    // this.ctx.fillRect(x + 30, y + 1, 100, 200 - 2)
    this.ctx.fillRect(x + this.w / 15, y + 1, this.w / 5, this.h - 2)
    this.ctx.fill()
    this.ctx.restore();
  }

  _drawTikets(num) {
    let x = this.x + this.delta;
    let y = this.y + this.delta;
    let i = num > 5 ? 5 : num;
    let shadow = true;
    for (; i > 0; i--) {
      this._tiket(x + i * 5, y + i * 5, i, shadow)
      if (shadow) {
        shadow = false
      }
    }
    // this._tiket(x + i * 5, y + i * 5, false)

    this.ctx.font = this.h / 1.17 + "px Impact";
    this.ctx.fillStyle = 'white';
    this.ctx.textAlign = "right";
    this.ctx.fillText(num, x + this.w / 1.25, y + this.h / 1.17);
    let rx = x + this.w * 0.85//420
    let ry = y + this.h * 0.5//100
    this.ctx.font = this.h * 0.3 + "px Impact";
    this.ctx.textAlign = "center";
    // this.ctx.beginPath();
    this.ctx.translate(rx, ry);
    this.ctx.rotate(Math.PI / 2);
    if (num == 1) {
      this.ctx.fillText("TIKET", 0, 0);
    } else {
      this.ctx.fillText("TIKETS", 0, 0);
    }
    this.ctx.restore();
    this.ctx.rotate(-Math.PI / 2);
    this.ctx.translate(-rx, -ry);
    return this
  }


  lable(text) {
    // var context = this.ctx.getContext("2d");
    this.ctx.font = 40 + "px sans-serif";
    this.ctx.fillStyle = '#2b6baaff';
    this.ctx.textAlign = "center";
    this.ctx.fillText(text, 100, 200);
    return this;
  }
}
// export Widget;
