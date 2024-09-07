
class Widget {
  constructor(canvas, lineWIdth, radius) {
    this.canvas = canvas;
    this.lineWidth = lineWIdth
    this.radius = radius
    this.angleStart = 5 * Math.PI / 6//2.5
    this.angleEnd = 13 * Math.PI / 6;
    this.gaugeVirt = this.angleEnd - this.angleStart;
  }
  clear() {
    var context = this.canvas.getContext("2d");
    context.reset();
    return this;
  }
  gauge(current, max) {
    var context = this.canvas.getContext("2d");
    context.beginPath();
    context.lineWidth = this.lineWidth;
    context.strokeStyle = "#c5c5c5ff"; // line color
    context.arc(this.canvas.width / 2,
      this.canvas.height * 0.5,
      this.radius,
      this.angleStart,
      this.angleEnd,
      false);
    context.stroke();
    
    context.beginPath();
    context.lineWidth = this.lineWidth;
    context.strokeStyle = '#71a6dbff'; // line color
    context.arc(this.canvas.width / 2,
      this.canvas.height * 0.5,
      this.radius,
      this.angleStart,
      this.angleStart + this.gaugeVirt * current / max,
      //this.angleEnd,
      /// this.gaugeMax * current / max,
      false);
    context.stroke();
    let textSize = this.radius;
    context.font = textSize+"px Impact";//64 for 15
    context.fillStyle = '#2b6baaff';
    context.textAlign = "center";
    context.fillText(current.toString(), this.canvas.width / 2, this.canvas.height * 0.5 + textSize/2.5);


    return this;
  }

  lable(text) {
    var context = this.canvas.getContext("2d");
    context.font = this.radius/3+"px sans-serif";
    context.fillStyle = '#2b6baaff';
    context.textAlign = "center";
    context.fillText(text, this.canvas.width / 2, this.canvas.height/2 + this.radius);
    return this
  }
}
// export Widget;
