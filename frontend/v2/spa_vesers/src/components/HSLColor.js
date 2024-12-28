export class HSLColor {
    constructor(h, s, l) {
        this.h = h
        this.s = s
        this.l = l
    }
    str() {
        return "HSL(" + this.h + "," + this.s + "%," + this.l + "%)"
    }
    clone() {
        return new HSLColor(this.h, this.s, this.l)
    }
}