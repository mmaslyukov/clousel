
export class Widget {
    constructor(canvas, lineWIdth, radius) {
        this.context = canvas.getContext("2d");
        this.canvas = canvas;
        this.lineWidth = lineWIdth
        this.radius = radius
        this.angleStart = 5 * Math.PI / 6//2.5
        this.angleEnd = 13 * Math.PI / 6;
        this.gaugeVirt = this.angleEnd - this.angleStart;
        this.colorsNormal = { font: "#2b6baaff", gauge: "#c5c5c5ff", actual: "#71a6dbff", status: "red" }
        this.colorsError = { font: "#c5c5c5ff", gauge: "#c5c5c5ff", actual: "#c5c5c5ff", status: "red" }
        this.colors = this.colorsNormal
        this.isOk = true;
        this.status = "";
    }

    clear() {
        this.context.reset();
        return this;
    }

    _drawInner() {
        this.context.beginPath();
        this.context.lineWidth = this.lineWidth;
        this.context.strokeStyle = this.colors.gauge; // line color
        this.context.arc(this.canvas.width / 2,
            this.canvas.height / 2,
            this.radius,
            this.angleStart,
            this.angleEnd,
            false);
        this.context.stroke();
    }

    withStatus(status) {
        this.status = status;
        this.isOk = status == "online";
        this.colors = this.isOk ? this.colorsNormal : this.colorsError;
        // if (status == "offline") {
        //     this.isOk = false;
        //     this.colors = this.colorsError;
        //     let textSize = this.radius/2;
        //     this.context.font = textSize + "px Impact";//64 for 15
        //     this.context.fillStyle = this.colors.status;
        //     this.context.textAlign = "center";
        //     this.context.fillText(status, this.canvas.width / 2, this.canvas.height * 0.5 + textSize / 2.5);
        // } else {
        //     this.colors = this.colorsNormal;
        // }
        return this;
    }

    draw(current, max) {

        this._drawInner()
        this.context.beginPath();
        this.context.lineWidth = this.lineWidth;
        this.context.strokeStyle = this.colors.actual; // line color
        this.context.arc(this.canvas.width / 2,
            this.canvas.height / 2,
            this.radius,
            this.angleStart,
            this.angleStart + this.gaugeVirt * current / max,
            false);
        this.context.stroke();
        let textSize = this.radius;
        this.context.font = textSize + "px Impact";//64 for 15
        this.context.fillStyle = this.colors.font;
        this.context.textAlign = "center";
        this.context.fillText(current.toString(), this.canvas.width / 2, this.canvas.height * 0.5 + textSize / 2.5);

        if (!this.isOk) {
            let textSize = this.radius/2;
            this.context.font = textSize + "px Impact";//64 for 15
            this.context.fillStyle = this.colors.status;
            this.context.textAlign = "center";
            this.context.fillText(this.status, this.canvas.width / 2, this.canvas.height * 0.5 + textSize / 2.5);

        }

        return this;
    }

    lable(text) {
        // var context = this.canvas.getContext("2d");
        this.context.font = this.radius / 3 + "px sans-serif";
        this.context.fillStyle = this.colors.font;
        this.context.textAlign = "center";
        this.context.fillText(text, this.canvas.width / 2, this.canvas.height / 2 + this.radius);
        return this
    }
}
