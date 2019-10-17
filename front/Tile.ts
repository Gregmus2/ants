export class Tile {
    readonly x: number;
    readonly y: number;
    readonly w: number;
    readonly h: number;
    private color: string;
    private ctx: CanvasRenderingContext2D;
    changed: boolean = true;
    
    constructor(x: number, y: number, w:number, h:number, color: string, ctx: CanvasRenderingContext2D) {
        this.x = x;
        this.y = y;
        this.w = w;
        this.h = h;
        this.color = color;
        this.ctx = ctx;
    }
    
    setColor(color: string){
        if (this.color !== color) {
            this.color = color;
            this.changed = true;
        }
    }

    draw(){
        this.ctx.fillStyle = this.color;
        this.ctx.fillRect(this.x, this.y, this.w, this.h);
    }
}