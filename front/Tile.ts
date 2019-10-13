export class Tile {
    readonly x: number;
    readonly y: number;
    readonly w: number;
    readonly h: number;
    private color: string;
    private readonly canvas: HTMLCanvasElement;
    private context: CanvasRenderingContext2D;
    changed: boolean = true;
    
    constructor(x: number, y: number, w:number, h:number, color: string) {
        this.x = x;
        this.y = y;
        this.w = w;
        this.h = h;
        this.color = color;
        this.canvas = document.createElement('canvas');
        this.canvas.width = w;
        this.canvas.height = h;
        this.context = this.canvas.getContext('2d');
    }
    
    setColor(color: string){
        if (this.color !== color) {
            this.color = color;
            this.changed = true;
        }
    }
    
    clearCanvas(){
        this.context.clearRect(0, 0, this.w, this.h);
    }
    
    draw(){
        this.clearCanvas();
        this.context.fillStyle = this.color;
        this.context.fillRect(0, 0, this.w, this.h);
        this.changed = false;
        
        return this.canvas;
    }
}