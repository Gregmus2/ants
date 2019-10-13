import {Tile} from "./Tile";

export class App{
    static canvas: HTMLCanvasElement;
    static ctx: CanvasRenderingContext2D;
    static elements: Array<Array<Tile>> = [];
    
    static init(){
        App.createSelect();
        App.canvas = document.createElement("canvas");
        let size = window.innerWidth < window.innerHeight ? window.innerWidth : window.innerHeight;
        App.canvas.width = size;
        App.canvas.height = size;
        App.canvas.style.border = '1px solid #000';
        App.canvas.addEventListener('contextmenu', event => event.preventDefault());
        document.body.appendChild(App.canvas);
        App.ctx = App.canvas.getContext('2d');
    }
    
    static loop(){
        App.update();
    
        for (let row of App.elements) {
            for (let el of row) {
                if (el.changed) {
                    App.ctx.clearRect(el.x, el.y, el.w, el.h);
                    App.ctx.drawImage(el.draw(), el.x, el.y);
                }
            }
        }
    }
    
    static update(){
        fetch('/get?id=0')
            .then(response => {
                switch (response.status) {
                    case 200:
                        return response.json()
                    case 404:
                        throw '404'
                }
            })
            .then(body => {
                for (let x in body) {
                    for (let y in body) {
                        App.elements[x][y].setColor(body[x][y])
                    }
                }
    
                setTimeout(() => (App.loop()), 50);
            })
            .catch(() => setTimeout(() => (App.loop()), 10000));
    }

    static createSelect(){
        fetch('/pipes')
            .then(response => response.json())
            .then(body => {
                let select = document.createElement("select");
                for (let id in body) {
                    let opt = document.createElement("option")
                    opt.value = id;
                    opt.text = id;
                    select.add(opt, null)
                }
                document.body.appendChild(select);
            })
    }
}