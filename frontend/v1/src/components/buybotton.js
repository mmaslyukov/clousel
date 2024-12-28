
export class BuyButton {
  constructor(canvas, w) {
    this.c = canvas
    this.ticketsMax = 5
    this.cw = w
    this.ch = w / 2
    this.ctx = canvas.getContext("2d");
    this.ctx.width = this.cw
    this.ctx.height = this.ch
    this.delta = 0;
    this.clickShift = w / 100
    this.x = this.cw * 0.01
    this.y = this.cw * 0.01
    this.w = this.cw * 0.8//0.91
    this.h = this.ch * 0.7//0.83
    this.lineWidth = w / 150
    this.lineShadow = w / 70
    this._bezel()
  }
  canvas() {
    return this.c
  }
  _bezel() {
    // this.ctx.lineWidth = this.cw/100;
    // this.ctx.fillStyle = '#ccccccff';
    // this.ctx.restore();
    // this.ctx.beginPath();
    // this.ctx.fillRect(0, 0, this.cw, this.ch)
    // this.ctx.fill()
    // this.ctx.strokeStyle = 'grey';
    // this.ctx.strokeRect(0, 0, this.cw, this.ch)
    // this.ctx.stroke()
  }
  clear() {
    this.ctx.reset();
    this._bezel();
    return this;
  }
  // click(x, y) {

  // }
  setPriceId(id) {
    this.priceId = id;
    return this
  }
  getPriceId() {
    return this.priceId;
  }
  setPrice(p) {
    this.price = p;
    return this;
  }
  getPrice() {
    return this.price
  }
  setTickets(t) {
    this.tickets = t;
    
    this.x = 10 + this.cw * (this.ticketsMax / this.tickets) / 100
    this.y = 8 + this.ch * (this.ticketsMax / this.tickets) / 100
    this.ctx.width = this.cw * (1 - this.tickets / 100)
    // this.c.width = this.cw * (1 - this.tickets / 100)
    return this
  }
  getTickets() {
    return this.tickets
  }
  redraw() {
    this._drawTickets(this.tickets)
    this._drawPriceTag(this.price)
    if (this.dim) {
      //  tx.fillStyle = 'grey';
      this.ctx.fillStyle = '#00000020';
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
    this.delta = this.clickShift
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

    let rx = x + this.cw / 25 + this.x + this.delta
    let ry = y + this.cw / 25 + this.y + this.delta
    this.ctx.lineWidth = this.lineWidth;


    this.ctx.beginPath();
    this.ctx.translate(rx, ry);
    this.ctx.rotate(Math.PI / 5);

    this.ctx.fillStyle = 'grey';
    let shadowDelta = this.lineShadow
    this.ctx.moveTo(x + shadowDelta, y + shadowDelta);

    this.ctx.lineTo(x + this.w * 0.07 + shadowDelta, y + this.h * 0.2 + shadowDelta);
    this.ctx.lineTo(x + this.w * 0.5 + shadowDelta, y + this.h * 0.2 + shadowDelta);
    this.ctx.lineTo(x + this.w * 0.5 + shadowDelta, y - this.h * 0.2 + shadowDelta);
    this.ctx.lineTo(x + this.w * 0.07 + shadowDelta, y - this.h * 0.2 + shadowDelta);

    this.ctx.lineTo(x + shadowDelta, y + shadowDelta);
    this.ctx.fill();

    this.ctx.beginPath();
    this.ctx.strokeStyle = '#c57b00ff';
    this.ctx.fillStyle = '#faf846ff';
    this.ctx.moveTo(x, y);


    this.ctx.lineTo(x + this.w * 0.07, y + this.h * 0.2);
    this.ctx.lineTo(x + this.w * 0.5, y + this.h * 0.2);
    this.ctx.lineTo(x + this.w * 0.5, y - this.h * 0.2);
    this.ctx.lineTo(x + this.w * 0.07, y - this.h * 0.2);
    this.ctx.lineTo(x, y);
    this.ctx.fill();
    this.ctx.stroke();

    // hole
    this.ctx.beginPath();
    this.ctx.fillStyle = '#ffffffff';
    this.ctx.arc(this.w * 0.06, y, this.h * 0.04, 0, Math.PI * 2, true);
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

  _ticket(x, y, num, shadow) {

    if (shadow) {
      this.ctx.fillStyle = 'grey';
      this.ctx.beginPath();
      this.ctx.fillRect(x + this.lineShadow, y + this.lineShadow, this.w, this.h)
      this.ctx.fill()
    }

    this.ctx.fillStyle = "HSL(125,55%," + (55 - 10 * num) + "%)";
    // this.ctx.fillStyle = '#43cc4fff';
    this.ctx.beginPath();
    this.ctx.fillRect(x, y, this.w, this.h)
    this.ctx.fill()

    this.ctx.lineWidth = this.lineWidth;
    this.ctx.strokeStyle = "HSL(125,55%," + (40 - 10 * num) + "%)";
    this.ctx.strokeRect(x, y, this.w, this.h)
    this.ctx.stroke()
    this.ctx.beginPath();
    this.ctx.fillStyle = "HSL(0,0%," + (100 - 10 * num) + "%)";
    // this.ctx.fillRect(x + 30, y + 1, 100, 200 - 2)
    this.ctx.fillRect(x + this.w / 15, y + 1, this.w / 5, this.h - 2)
    this.ctx.fill()
    this.ctx.restore();
  }

  _drawTickets(num) {
    let x = this.x + this.delta;
    let y = this.y + this.delta;
    let i = num > this.ticketsMax ? this.ticketsMax : num;
    let shadow = true;
    for (; i > 0; i--) {
      this._ticket(x + i * 5, y + i * 5, i - 1, shadow)
      if (shadow) {
        shadow = false
      }
    }
    this.ctx.font = this.h / 1.17 + "px Impact";
    this.ctx.fillStyle = 'white';
    this.ctx.textAlign = "right";
    this.ctx.fillText(num, x + this.w / 1.25, y + this.h / 1.13);
    let rx = x + this.w * 0.85//420
    let ry = 0
    let text = ""
    if (num == 1) {
      ry = y + this.h * 0.53
      text = "TICKET"
    } else {
      ry = y + this.h * 0.56
      text = "TICKETS"
    }
    this.ctx.font = this.h * 0.25 + "px Impact";
    this.ctx.textAlign = "center";
    // this.ctx.beginPath();
    this.ctx.translate(rx, ry);
    this.ctx.rotate(Math.PI / 2);
    this.ctx.fillText(text, 0, 0);
    this.ctx.restore();
    this.ctx.rotate(-Math.PI / 2);
    this.ctx.translate(-rx, -ry);
    return this
  }


  lable(text) {
    this.ctx.font = 40 + "px sans-serif";
    this.ctx.fillStyle = '#2b6baaff';
    this.ctx.textAlign = "center";
    this.ctx.fillText(text, 100, 200);
    return this;
  }
}
